package routes

import (
	client "github.com/braunkc/todo-app/api-service-demo/internal/grpc"
	"github.com/braunkc/todo-app/api-service-demo/internal/http/handlers"
	"github.com/braunkc/todo-app/api-service-demo/internal/http/middlewares"
	"github.com/braunkc/todo-app/api-service-demo/internal/token"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, jwtService token.JWTService, dbService client.DatabaseService) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("register", handlers.Register(jwtService, dbService))
			v1.POST("login", handlers.Login(jwtService, dbService))
			v1.POST("logout", handlers.Logout())

			user := v1.Group("/user")
			user.Use(middlewares.AuthMiddleware(jwtService))
			{
				user.DELETE("/", handlers.DeleteUser(jwtService, dbService))
			}

			task := v1.Group("/task")
			task.Use(middlewares.AuthMiddleware(jwtService))
			{
				task.POST("/", handlers.CreateTask(dbService))
				task.PATCH("/", handlers.UpdateTask(dbService))
				task.DELETE("/", handlers.DeleteTask(dbService))
			}

			// return tasks in json
			v1.POST("/tasks", middlewares.AuthMiddleware(jwtService), handlers.GetTasks(dbService))
		}
	}

	// landing
	r.GET("/", handlers.RenderLanding())
	r.GET("auth", handlers.RenderAuth())
	r.GET("/tasks", middlewares.AuthMiddleware(jwtService), handlers.RenderTasks())
}
