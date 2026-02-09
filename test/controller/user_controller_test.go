package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpcontroller "product-app/internal/adapters/http/controller"
	"product-app/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupUserController() *httpcontroller.UserController {
	fakeRepo := NewFakeUserRepository(nil)
	userService := usecase.NewUserService(fakeRepo)
	return httpcontroller.NewUserController(userService)
}

func Test_ShouldRegisterUser(t *testing.T) {
	e := echo.New()
	userJSON := `{
		"username": "johndoe",
		"email": "john@example.com",
		"password": "secret123",
		"first_name": "John",
		"last_name": "Doe"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userController := setupUserController()

	err := userController.Register(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func Test_ShouldLoginUser(t *testing.T) {
	e := echo.New()
	userController := setupUserController()

	registerJSON := `{
		"username": "johndoe",
		"email": "john@example.com",
		"password": "secret123",
		"first_name": "John",
		"last_name": "Doe"
	}`
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(registerJSON))
	registerReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()
	registerCtx := e.NewContext(registerReq, registerRec)

	err := userController.Register(registerCtx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, registerRec.Code)

	loginJSON := `{
		"username_or_email": "johndoe",
		"password": "secret123"
	}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(loginJSON))
	loginReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	loginRec := httptest.NewRecorder()
	loginCtx := e.NewContext(loginReq, loginRec)

	err = userController.Login(loginCtx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, loginRec.Code)

	var response map[string]interface{}
	json.Unmarshal(loginRec.Body.Bytes(), &response)
	assert.Equal(t, "Login successful", response["message"])
	assert.NotEmpty(t, response["token"])
}

func Test_ShouldGetUserById(t *testing.T) {
	e := echo.New()
	userController := setupUserController()

	registerJSON := `{
		"username": "johndoe",
		"email": "john@example.com",
		"password": "secret123",
		"first_name": "John",
		"last_name": "Doe"
	}`
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(registerJSON))
	registerReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()
	registerCtx := e.NewContext(registerReq, registerRec)

	_ = userController.Register(registerCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := userController.GetUserById(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "johndoe", response["username"])
	assert.Equal(t, "john@example.com", response["email"])
}

func Test_ShouldUpdateUser(t *testing.T) {
	e := echo.New()
	userController := setupUserController()

	registerJSON := `{
		"username": "johndoe",
		"email": "john@example.com",
		"password": "secret123",
		"first_name": "John",
		"last_name": "Doe"
	}`
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(registerJSON))
	registerReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()
	registerCtx := e.NewContext(registerReq, registerRec)

	_ = userController.Register(registerCtx)

	updateJSON := `{
		"username": "johnupdated",
		"email": "john.updated@example.com",
		"first_name": "John",
		"last_name": "Updated"
	}`
	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", strings.NewReader(updateJSON))
	updateReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	updateRec := httptest.NewRecorder()
	updateCtx := e.NewContext(updateReq, updateRec)
	updateCtx.SetParamNames("id")
	updateCtx.SetParamValues("1")

	err := userController.UpdateUser(updateCtx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, updateRec.Code)
}

func Test_ShouldDeleteUser(t *testing.T) {
	e := echo.New()
	userController := setupUserController()

	registerJSON := `{
		"username": "johndoe",
		"email": "john@example.com",
		"password": "secret123",
		"first_name": "John",
		"last_name": "Doe"
	}`
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(registerJSON))
	registerReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	registerRec := httptest.NewRecorder()
	registerCtx := e.NewContext(registerReq, registerRec)

	_ = userController.Register(registerCtx)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := userController.DeleteUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
