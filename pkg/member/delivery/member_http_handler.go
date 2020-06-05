package delivery

import (
	"net/http"

	"github.com/hardyantz/data-encryption/helpers"
	"github.com/hardyantz/data-encryption/pkg/member/domain"
	"github.com/hardyantz/data-encryption/pkg/member/usecase"
	"github.com/labstack/echo"
)

type HTTPHandler struct {
	memberUseCase *usecase.MemberUseCase
}

func (h *HTTPHandler) Mount(group *echo.Group) {
	group.POST("/", h.Create)
	group.POST("/email", h.GetEmail) // add save redis
	group.GET("/:id", h.GetOne) // add save redis
	group.GET("/", h.GetAll) // add save redis
	group.PUT("/:id", h.Update) // update redis
	group.DELETE("/:id", h.Delete) // delete redis
}

func NewHTTPHandler(mu *usecase.MemberUseCase) *HTTPHandler {
	return &HTTPHandler{mu}
}

var meta helpers.Meta

func (h *HTTPHandler) Create(c echo.Context) error {
	newMember := new(domain.Member)

	if err := c.Bind(&newMember); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, "invalid params", nil, meta)
	}

	if err := h.memberUseCase.Create(newMember); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusCreated, "Data successfully created", newMember, meta)
}

func (h *HTTPHandler) GetAll(c echo.Context) error {
	var params domain.Parameters

	params.Limit = c.QueryParam("limit")
	params.OrderBy = c.QueryParam("orderBy")
	params.Sort = c.QueryParam("sort")

	members, err := h.memberUseCase.GetAll(params)
	if err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusOK, "Data successfully fetch", members, meta)
}

func (h *HTTPHandler) GetOne(c echo.Context) error {
	id := c.Param("id")

	member, err := h.memberUseCase.GetOne(id)
	if err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusOK, "Data successfully fetch", member, meta)
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

	member, err := h.memberUseCase.GetOneEmail(bindEmail.Email)
	if err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusOK, "Data successfully fetch", member, meta)
}

func (h *HTTPHandler) Update(c echo.Context) error {
	member := new(domain.Member)
	id := c.Param("id")

	if err := c.Bind(&member); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, "invalid params", nil, meta)
	}

	if err := h.memberUseCase.Update(id, member); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusCreated, "Data successfully update", member, meta)
}

func (h *HTTPHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	if err := h.memberUseCase.Delete(id); err != nil {
		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
	}

	return helpers.HTTPResponse(c, http.StatusCreated, "Data successfully created", nil, meta)
}
//
//func (h *HTTPHandler) GetEmailLogin(c echo.Context) error {
//	type emailBody struct {
//		Email string `json:"email"`
//	}
//	bindEmail := new(emailBody)
//	err := c.Bind(&bindEmail)
//	if err != nil {
//		return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
//	}
//
//	return helpers.HTTPResponse(c, http.StatusOK, "success", helpers.HashAndSalt([]byte(bindEmail.Email)), nil)
//
//	//member, err := h.memberUseCase.GetEmailLogin(bindEmail.Email)
//	//if err != nil {
//	//	return helpers.HTTPResponse(c, http.StatusBadRequest, err.Error(), nil, meta)
//	//}
//	//
//	//return helpers.HTTPResponse(c, http.StatusCreated, "Email found", member, meta)
//}
