package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/saltbo/gopkg/ginutil"
	"github.com/storyicon/grbac"

	"github.com/saltbo/moreu/client"
	"github.com/saltbo/moreu/model"
	"github.com/saltbo/moreu/service"
)

func LoginAuth() gin.HandlerFunc {
	return LoginAuthWithRoles(nil)
}

func LoginAuthWithRoles(roles grbac.Rules) gin.HandlerFunc {
	if roles != nil {
		defaultRules = append(defaultRules, roles...)
	}

	ctrl, err := grbac.New(grbac.WithRules(defaultRules))
	if err != nil {
		log.Fatalln(err)
	}

	return func(c *gin.Context) {
		token, err := tokenCookieGet(c)
		if errors.Is(err, http.ErrNoCookie) {
			token, _ = service.TokenCreate("guest", 30, model.RoleGuest) // 未登录状态颁发一个匿名Token
		}

		rc, err := service.TokenVerify(token)
		if err != nil {
			tokenError(c, err)
			return
		}

		state, err := ctrl.IsRequestGranted(c.Request, rc.Roles)
		if err != nil {
			grantedError(c, err)
			return
		}

		if !state.IsGranted() {
			notGrantedError(c)
			return
		}

		uxSet(c, rc.Subject)
	}
}

func tokenError(c *gin.Context, err error) {
	accept := c.Request.Header.Get("Accept")
	if strings.Contains(accept, gin.MIMEJSON) {
		ginutil.JSONUnauthorized(c, err)
	} else {
		ginutil.FoundRedirect(c, service.Link2SignIn(c.Request.URL.RequestURI()))
	}
}

func grantedError(c *gin.Context, err error) {
	accept := c.Request.Header.Get("Accept")
	if strings.Contains(accept, gin.MIMEJSON) {
		ginutil.JSONServerError(c, err)
	} else {
		ginutil.FoundRedirect(c, service.Link2ServerError(err))
	}
}

func notGrantedError(c *gin.Context) {
	accept := c.Request.Header.Get("Accept")
	if strings.Contains(accept, gin.MIMEJSON) {
		ginutil.JSONForbidden(c, fmt.Errorf("access deny"))
	} else {
		ginutil.FoundRedirect(c, service.Link2SignIn(c.Request.URL.RequestURI()))
	}
}

// auth k-v
const (
	ctxUxKey = "ctx-ux"

	cookieTokenKey = "moreu-token"
	cookieRoleKey  = "moreu-role"
)

func uxSet(c *gin.Context, ux string) {
	c.Set(ctxUxKey, ux)
	client.InjectUx(c.Request, ux)
}

func uxGet(c *gin.Context) string {
	return c.GetString(ctxUxKey)
}

func tokenCookieSet(c *gin.Context, token string, expireSec int) {
	c.SetCookie(cookieTokenKey, token, expireSec, "/", "", false, true)
}

func tokenCookieGet(c *gin.Context) (string, error) {
	return c.Cookie(cookieTokenKey)
}
