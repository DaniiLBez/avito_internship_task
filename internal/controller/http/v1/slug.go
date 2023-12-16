package v1

import (
	"errors"
	"github.com/DaniiLBez/avito_internship_task/internal/entities"
	"github.com/DaniiLBez/avito_internship_task/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type slugRoutes struct {
	slugService service.Slug
}

func newSlugRoutes(g *echo.Group, slugService service.Slug) {
	r := &slugRoutes{slugService: slugService}

	g.POST("/create", r.createSlug)
	g.POST("/add", r.addUserToSlug)
	g.POST("/delete", r.deleteSlug)
	g.GET("/active/:id", r.getActiveSlugs)
}

type createSlugRequest struct {
	Name string `json:"name" validator:"required,min=4,max=32"`
}

func (r *slugRoutes) createSlug(c echo.Context) error {
	var slugName createSlugRequest

	if err := c.Bind(&slugName); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(slugName); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.slugService.CreateSlug(c.Request().Context(), slugName.Name)

	if err != nil {
		if errors.Is(err, service.ErrSlugAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}

type addSlugsRequest struct {
	SlugsToAdd    []string `json:"slugs-to-add" validate:"required"`
	SlugsToDelete []string `json:"slugs-to-delete" validate:"required"`
	UserId        int      `json:"user-id" validate:"required"`
}

func (r *slugRoutes) addUserToSlug(c echo.Context) error {
	var slugsInput addSlugsRequest

	if err := c.Bind(&slugsInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(slugsInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.slugService.AddUserToSlug(
		c.Request().Context(),
		slugsInput.SlugsToAdd,
		slugsInput.SlugsToDelete,
		slugsInput.UserId,
	)

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrSlugNotFound) {
			newErrorResponse(c, http.StatusBadRequest, "invalid request body"+err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error"+err.Error())
		return err
	}

	type response struct {
		Success bool `json:"success"`
	}

	return c.JSON(http.StatusOK, response{Success: true})
}

type deleteRequest struct {
	Name string `json:"name" validator:"required,min=4,max=32"`
}

func (r *slugRoutes) deleteSlug(c echo.Context) error {
	var deleteInput deleteRequest

	if err := c.Bind(&deleteInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body"+err.Error())
		return err
	}

	if err := c.Validate(deleteInput); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.slugService.DeleteSlug(c.Request().Context(), deleteInput.Name)

	if err != nil {
		if errors.Is(err, service.ErrCannotDeleteSlug) {
			newErrorResponse(c, http.StatusBadRequest, "cannot delete the slug with name: "+deleteInput.Name)
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error"+err.Error())
		return err
	}

	type response struct {
		Success bool `json:"success"`
	}

	return c.JSON(http.StatusOK, response{Success: true})
}

func (r *slugRoutes) getActiveSlugs(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request params"+err.Error())
		return err
	}

	var slugs []entities.Slug
	slugs, err = r.slugService.GetActiveSlugs(c.Request().Context(), userId)
	if err != nil {
		if errors.Is(err, service.ErrCannotGetSlug) {
			newErrorResponse(c, http.StatusBadRequest, "cannot get slugs"+err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error"+err.Error())
	}

	type response struct {
		Slugs []entities.Slug `json:"slugs"`
	}

	return c.JSON(http.StatusOK, response{Slugs: slugs})
}
