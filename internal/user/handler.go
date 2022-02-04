package user

import (
	"github.com/gin-gonic/gin"
	"github.com/rafimuhammad01/portofolio-api/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var (
		errors      []string
		requestBody LoginAPIRequest
	)

	// Input Validation
	err := c.Bind(&requestBody)
	if err != nil {
		errors = append(errors, err.Error())
	}

	if requestBody.Username == "" {
		errors = append(errors, "username is required")
	}

	if requestBody.Password == "" {
		errors = append(errors, "password is required")
	}

	if len(errors) != 0 {
		c.JSON(http.StatusBadRequest, &CreateUserAPIResponse{
			Status:  http.StatusBadRequest,
			Message: "bad request",
			Errors:  errors,
		})
		return
	}

	res, err := h.service.Login(requestBody)
	if err != nil {
		logrus.Error("[error] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	c.JSON(res.Status, res)
}

func (h *Handler) GetAllUser(c *gin.Context) {
	res, err := h.service.List()
	if err != nil {
		logrus.Error("[error] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}
	c.JSON(res.Status, res)
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var (
		errors      []string
		requestBody CreateUserAPIRequest
	)

	// Input Validation
	err := c.Bind(&requestBody)
	if err != nil {
		errors = append(errors, err.Error())
	}

	if requestBody.FullName == "" {
		errors = append(errors, "full_name is required")
	}

	if requestBody.Username == "" {
		errors = append(errors, "username is required")
	}

	if requestBody.Password == "" {
		errors = append(errors, "password is required")
	}

	if requestBody.FullName != "" && len(requestBody.FullName) < 3 {
		errors = append(errors, "full name should be greater than 3 characters")
	}

	if requestBody.Username != "" && len(requestBody.Username) < 3 {
		errors = append(errors, "username should be greater than 3 characters")
	}

	if requestBody.Password != "" && len(requestBody.Password) < 3 {
		errors = append(errors, "password should be greater than 8 characters")
	}

	if len(errors) != 0 {
		c.JSON(http.StatusBadRequest, &CreateUserAPIResponse{
			Status:  http.StatusBadRequest,
			Message: "bad request",
			Errors:  errors,
		})
		return
	}

	res, err := h.service.Create(requestBody)
	if err != nil {
		logrus.Error("[error] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	c.JSON(res.Status, res)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	payload, err := utils.GetPayloadFromContext(c)
	if err != nil {
		logrus.Error("[error] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	var (
		errors []string
	)

	userID := payload.UserID

	if len(errors) != 0 {
		c.JSON(http.StatusBadRequest, &CreateUserAPIResponse{
			Status:  http.StatusBadRequest,
			Message: "bad request",
			Errors:  errors,
		})
		return
	}

	res, err := h.service.Get(userID)
	if err != nil {
		logrus.Error("[error] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	c.JSON(res.Status, res)
}

func (h Handler) RefreshToken(c *gin.Context) {
	var accessToken RefreshTokenAPIRequest

	if err := c.BindJSON(&accessToken); err != nil {
		logrus.Error("[error] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	res, err := h.service.RefreshToken(accessToken.AccessToken)
	if err != nil {
		logrus.Error("[error] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	c.JSON(res.Status, res)
}
