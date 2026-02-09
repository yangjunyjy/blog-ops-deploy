package main

import (
	"my-blog-backend/internal/server"
)

// @title           My Blog API
// @version         1.0
// @description     This is a blog backend API with Swagger documentation.

// @contact.name   API Support
// @contact.url    http://www.myblog.com/support
// @contact.email  support@myblog.com

// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host           localhost:8080
// @BasePath       /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	server.Run()
}
