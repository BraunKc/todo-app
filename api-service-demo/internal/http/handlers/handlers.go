package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/braunkc/todo-app/api-service-demo/internal/dto"
	client "github.com/braunkc/todo-app/api-service-demo/internal/grpc"
	"github.com/braunkc/todo-app/api-service-demo/internal/token"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func Register(jwtService token.JWTService, dbService client.DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		confirmPassword := c.PostForm("confirm_password")

		if username == "" || password == "" {
			c.Abort()
			return
		}

		if password != confirmPassword {
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := dbService.CreateUser(ctx, &dto.CreateUserRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		token, err := jwtService.Generate(resp.User.ID)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.SetCookie("Authorization", token, 3600*24*7, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/tasks")
	}
}

func Login(jwtService token.JWTService, dbService client.DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "" || password == "" {
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		user, err := dbService.Authenticate(ctx, username, password)
		if err != nil {
			c.Abort()
			return
		}

		token, err := jwtService.Generate(user.ID)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.SetCookie("Authorization", token, 3600*24*7, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/tasks")
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("Authorization", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, nil)
	}
}

func DeleteUser(jwtService token.JWTService, dbService client.DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := dbService.DeleteUserByID(ctx, &dto.DeleteUserByIDRequest{
			ID: userID.(string),
		})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.SetCookie("Authorization", "", 0, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/landing")
	}
}

func CreateTask(dbService client.DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTaskRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		md := metadata.New(map[string]string{
			"userID": userID.(string),
		})

		ctx := c.Request.Context()
		ctx = metadata.NewOutgoingContext(ctx, md)
		task, err := dbService.CreateTask(ctx, &req)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}

func UpdateTask(dbService client.DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.UpdateTaskRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resp, err := dbService.GetTask(context.Background(), &dto.GetTaskRequest{
			ID: req.ID,
		})
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if resp.Task.UserID != userID {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		md := metadata.New(map[string]string{
			"userID": userID.(string),
		})

		ctx := c.Request.Context()
		ctx = metadata.NewOutgoingContext(ctx, md)
		task, err := dbService.UpdateTask(ctx, &req)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

func DeleteTask(dbService client.DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.DeleteTasksByIDRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		for _, taskID := range req.IDs {
			resp, err := dbService.GetTask(context.Background(), &dto.GetTaskRequest{
				ID: taskID,
			})
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if resp.Task.UserID != userID.(string) {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		_, err := dbService.DeleteTasksByID(context.Background(), &req)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

func GetTasks(dbService client.DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.GetTasksRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		md := metadata.New(map[string]string{
			"userID": userID.(string),
		})

		ctx := c.Request.Context()
		ctx = metadata.NewOutgoingContext(ctx, md)
		resp, err := dbService.GetTasks(ctx, &req)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func RenderLanding() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "landing.html", nil)
	}
}

func RenderAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.html", nil)
	}
}

func RenderTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "tasks.html", nil)
	}
}
