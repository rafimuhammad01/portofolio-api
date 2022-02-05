package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rafimuhammad01/portofolio-api/internal/jwt"
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
		errorList   []string
		requestBody LoginAPIRequest
	)

	// Input Validation
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		errorList = append(errorList, err.Error())
	}

	if requestBody.Username == "" {
		errorList = append(errorList, "username is required")
	}

	if requestBody.Password == "" {
		errorList = append(errorList, "password is required")
	}

	if len(errorList) != 0 {
		c.JSON(http.StatusBadRequest, &LoginAPIResponse{
			Status:  http.StatusBadRequest,
			Message: "bad request",
			Errors:  errorList,
		})
		return
	}

	accessToken, refreshToken, expAt, err := h.service.Login(requestBody.Username, requestBody.Password, c)
	switch errors.Cause(err) {
	case ErrInvalidUsernameOrPassword:
		fallthrough
	case ErrUserNotFound:
		c.JSON(http.StatusUnauthorized, &LoginAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Errors:  []string{ErrInvalidUsernameOrPassword.Error()},
		})
		return
	case ErrInternalServer:
		logrus.Error("[error while using login service]", err.Error())
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	data := jwt.JWTAPIResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    expAt,
	}

	c.JSON(http.StatusCreated, LoginAPIResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &data,
	})
}

func (h *Handler) GetAllUser(c *gin.Context) {
	res, err := h.service.List()
	switch errors.Cause(err) {
	case ErrUserNotFound:
		res = &ListUser{
			Users: []User{},
			Count: 0,
		}
		return
	case ErrInternalServer:
		logrus.Error("[error while using list user service]", err.Error())
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	c.JSON(http.StatusOK, &ListUserAPIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    res,
	})
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var (
		errorList   []string
		requestBody CreateUserAPIRequest
	)

	// Input Validation
	err := c.Bind(&requestBody)
	if err != nil {
		errorList = append(errorList, err.Error())
	}

	if requestBody.FullName == "" {
		errorList = append(errorList, "full_name is required")
	}

	if requestBody.Username == "" {
		errorList = append(errorList, "username is required")
	}

	if requestBody.Password == "" {
		errorList = append(errorList, "password is required")
	}

	if requestBody.FullName != "" && len(requestBody.FullName) < 3 {
		errorList = append(errorList, "full name should be greater than 3 characters")
	}

	if requestBody.Username != "" && len(requestBody.Username) < 3 {
		errorList = append(errorList, "username should be greater than 3 characters")
	}

	if requestBody.Password != "" && len(requestBody.Password) < 3 {
		errorList = append(errorList, "password should be greater than 8 characters")
	}

	if len(errorList) != 0 {
		c.JSON(http.StatusBadRequest, &CreateUserAPIResponse{
			Status:  http.StatusBadRequest,
			Message: "bad request",
			Errors:  errorList,
		})
		return
	}

	res, err := h.service.Create(requestBody.Username, requestBody.FullName, requestBody.Password)
	switch errors.Cause(err) {
	case ErrUsernameAlreadyExist:
		c.JSON(http.StatusBadRequest, &CreateUserAPIResponse{
			Status:  http.StatusBadRequest,
			Message: "bad request",
			Errors:  []string{err.Error()},
		})
		return
	case ErrInternalServer:
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	c.JSON(http.StatusCreated, &CreateUserAPIResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    res,
	})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	payload, err := utils.GetPayloadFromContext(c)
	if err != nil {
		logrus.Error("[error while extracting context] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	userID := payload.UserID

	res, err := h.service.Get(userID)
	if err != nil {
		if errors.Cause(err) == ErrUserNotFound {
			c.JSON(http.StatusNotFound, GetUserByIDAPIResponse{
				Status:  http.StatusNotFound,
				Message: "not found",
				Data:    &User{},
			})
			return
		}
		logrus.Error("[error] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	c.JSON(http.StatusOK, GetUserByIDAPIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    res,
	})
}

func (h Handler) RefreshToken(c *gin.Context) {
	var refreshTokenBody RefreshTokenAPIRequest

	if err := c.BindJSON(&refreshTokenBody); err != nil {
		logrus.Error("[error while binding body to struct] ", err)
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	accessToken, refreshToken, duration, err := h.service.RefreshToken(refreshTokenBody.RefreshToken, c)
	switch errors.Cause(err) {
	case jwt.ErrInvalidToken:
		fallthrough
	case jwt.ErrExpiredToken:
		c.JSON(http.StatusUnauthorized, LoginAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Errors:  []string{err.Error()},
		})
	case jwt.ErrIntervalServer:
		logrus.Error("[error while using refresh token service] ", err.Error())
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorHandler())
		return
	}

	c.JSON(http.StatusOK, LoginAPIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &jwt.JWTAPIResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiredAt:    duration,
		},
	})
}
