package v1

import (
	"errors"
	"github.com/DaniiLBez/avito_internship_task/internal/entities"
	"github.com/DaniiLBez/avito_internship_task/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type userRoutes struct {
	userService service.User
}

func newUserRoutes(g *echo.Group, userService service.User) {
	r := userRoutes{userService: userService}

	g.GET("/:id", r.getById)
	g.GET("", r.getByName)
}

func (r *userRoutes) getById(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request params"+err.Error())
		return err
	}

	var user entities.User
	user, err = r.userService.GetUserById(c.Request().Context(), userId)

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(c, http.StatusNotFound, "user with id: "+strconv.Itoa(userId)+" not found")
			return err
		}

		if errors.Is(err, service.ErrCannotGetUser) {
			newErrorResponse(c, http.StatusBadRequest, "cannot get user with id: "+strconv.Itoa(userId))
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error: "+err.Error())
		return err
	}

	type response struct {
		User entities.User `json:"user"`
	}

	return c.JSON(http.StatusOK, response{User: user})
}

func (r *userRoutes) getByName(c echo.Context) error {
	userName := c.QueryParam("username")
	if userName == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid query params")
		return c.String(http.StatusBadRequest, "Username parameter is required")
	}

	user, err := r.userService.GetUserByUsername(c.Request().Context(), userName)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(c, http.StatusNotFound, "user with name: "+userName+" not found")
			return err
		}

		if errors.Is(err, service.ErrCannotGetUser) {
			newErrorResponse(c, http.StatusBadRequest, "cannot get user with name: "+userName)
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error: "+err.Error())
		return err
	}

	type response struct {
		User entities.User `json:"user"`
	}

	return c.JSON(http.StatusOK, response{User: user})
}
