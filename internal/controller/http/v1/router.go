package v1

import (
	"github.com/DaniiLBez/avito_internship_task/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
	"os"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))

	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error {
		return c.NoContent(200)
	})

	auth := handler.Group("/auth")
	{
		newAuthRoutes(auth, services.Auth)
	}

	authMiddleware := &AuthMiddleware{services.Auth}
	v1 := handler.Group("/api/v1", authMiddleware.UserIdentity)
	{
		newUserRoutes(v1.Group("/user"), services.User)
		newSlugRoutes(v1.Group("/slug"), services.Slug)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		slog.Error("cannot open file", err)
	}
	return file
}
