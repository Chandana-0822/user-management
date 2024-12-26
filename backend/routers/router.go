package routers

import (
	"backend/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "backend/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func InitRouter() *echo.Echo {
	e := echo.New()

	// Enable CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Serve Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Serve the swagger.yaml file
	e.GET("/swagger/doc.yaml", func(c echo.Context) error {
		return c.File("docs/swagger.yaml") // Update the path if swagger.yaml is in another location
	})
	e.GET("/swagger/doc.json", func(c echo.Context) error {
		return c.File("docs/swagger.json") // Path to doc.json
	})

	v1 := e.Group("/v1")
	v1.GET("/users", handlers.GetUsersHandler)
	v1.POST("/users", handlers.CreateUserHandler)
	v1.PUT("/users/:id", handlers.UpdateUserHandler)
	v1.DELETE("/users/:id", handlers.DeleteUserHandler)

	// Search username route
	v1.POST("/users/search", handlers.SearchUsernameHandler)

	return e
}
