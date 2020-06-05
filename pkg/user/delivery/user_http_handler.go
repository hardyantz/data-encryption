package delivery

import (
	"net/http"

	"github.com/hardyantz/data-encryption/helpers"
	"github.com/hardyantz/data-encryption/pkg/user/domain"
	"github.com/hardyantz/data-encryption/pkg/user/usecase"
	"github.com/labstack/echo"
)

type HTTPHandler struct {
	userUseCase *usecase.UserUseCase
}

func (h *HTTPHandler) Mount(group *echo.Group) {
	group.POST("/", h.Create)
	group.POST("/email", h.GetEmail)
	group.GET("/:id", h.GetOne)
	group.GET("/", h.GetAll)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func NewHTTPHandler(mu *usecase.UserUseCase) *HTTPHandler {
	return &HTTPHandler{mu}
}

var meta helpers.Meta

func (h *HTTPHandler) Create(c echo.Context) error {
	newUser := new(domain.User)

	if err := c.Bind(&newUser); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, "invalid params", nil, meta)
	}

	if err := h.userUseCase.Create(newUser); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusCreated, "Data successfully created", newUser, meta)
}

func (h *HTTPHandler) GetAll(c echo.Context) error {
	var params domain.Parameters

	params.Limit = c.QueryParam("limit")
	params.OrderBy = c.QueryParam("orderBy")
	params.Sort = c.QueryParam("sort")

	users, err := h.userUseCase.GetAll(params)
	if err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusOK, "Data successfully fetch", users, meta)
}

func (h *HTTPHandler) GetOne(c echo.Context) error {
	id := c.Param("id")

	user, err := h.userUseCase.GetOne(id)
	if err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusOK, "Data successfully fetch", user, meta)
}

func (h *HTTPHandler) GetEmail(c echo.Context) error {
	type emailBody struct {
		Email string `json:"email"`
	}
	bindEmail := new(emailBody)
	err := c.Bind(&bindEmail)
	if err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	user, err := h.userUseCase.GetOneEmail(bindEmail.Email)
	if err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusOK, "Data successfully fetch", user, meta)
}

func (h *HTTPHandler) Update(c echo.Context) error {
	user := new(domain.User)
	id := c.Param("id")

	if err := c.Bind(&user); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, "invalid params", nil, meta)
	}

	if err := h.userUseCase.Update(id, user); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusCreated, "Data successfully update", user, meta)
}

func (h *HTTPHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	if err := h.userUseCase.Delete(id); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusCreated, "Data successfully created", nil, meta)
}
