package v1

import (
	"github.com/DaniiLBez/avito_internship_task/internal/service"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"net/http"
	"strings"
)

const (
	userIdCtx = "userId"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (h *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			slog.Error("AuthMiddleware.UserIdentity: bearerToken: ", ErrInvalidAuthHeader)
			newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			return nil
		}

		userId, err := h.authService.ParseToken(token)
		if err != nil {
			slog.Error("AuthMiddleware.UserIdentity: h.authService.ParseToken: ", err)
			newErrorResponse(c, http.StatusUnauthorized, ErrCannotParseToken.Error())
			return err
		}

		c.Set(userIdCtx, userId)

		return next(c)
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
