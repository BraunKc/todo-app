package server

import (
	client "github.com/braunkc/todo-app/api-service-demo/internal/grpc"
	"github.com/braunkc/todo-app/api-service-demo/internal/http/routes"
	"github.com/braunkc/todo-app/api-service-demo/internal/token"
	"github.com/gin-gonic/gin"
)

func New(jwtService token.JWTService, dbService client.DatabaseService) *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("./web/templates/*")
	r.Static("/css", "./web/static/css")
	r.Static("/js", "./web/static/js")

	routes.Setup(r, jwtService, dbService)

	return r
}
