/*
 * Moreu API
 *
 * This is a moreu server.
 *
 * API version: 1.0
 * Contact: saltbo@foxmail.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package client

type BindBodyToken struct {
	Captcha string `json:"captcha,omitempty"`
	Email string `json:"email"`
	Password string `json:"password"`
}
