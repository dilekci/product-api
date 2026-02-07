package controller

import (
	"net/http"
	"product-app/controller/response"
	"product-app/domain"
	"product-app/middleware"
	"product-app/service"
	"time"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService service.IUserService
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type UpdateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserResponse struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (userController *UserController) RegisterRoutes(e *echo.Echo) {
	// Public routes (no authentication required)
	e.POST("/api/v1/auth/register", userController.Register)
	e.POST("/api/v1/auth/login", userController.Login)

	// Protected routes (authentication required)
	protected := e.Group("/api/v1/users", middleware.JWTMiddleware())
	protected.GET("/:id", userController.GetUserById)
	protected.PUT("/:id", userController.UpdateUser)
	protected.DELETE("/:id", userController.DeleteUser)
}

func (userController *UserController) Register(c echo.Context) error {
	req, err := bindRegisterRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if err := userController.userService.Register(req.Username, req.Email, req.Password, req.FirstName, req.LastName); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, MessageResponse{
		Message: "User registered successfully",
	})
}

func (userController *UserController) Login(c echo.Context) error {
	req, err := bindLoginRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	user, err := userController.userService.Login(req.UsernameOrEmail, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.Id, user.Username, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, buildLoginResponse(token, user))
}

func (userController *UserController) GetUserById(c echo.Context) error {
	userId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid user ID",
		})
	}

	user, err := userController.userService.GetById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, buildUserResponse(user))
}

func (userController *UserController) UpdateUser(c echo.Context) error {
	userId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid user ID",
		})
	}

	updateReq, err := bindUpdateUserRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Get existing user
	user, err := userController.userService.GetById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	applyUserUpdate(&user, updateReq)

	if err := userController.userService.UpdateUser(user); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, MessageResponse{
		Message: "User updated successfully",
	})
}

func (userController *UserController) DeleteUser(c echo.Context) error {
	userId, err := parsePositiveIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid user ID",
		})
	}

	if err := userController.userService.DeleteById(userId); err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, MessageResponse{
		Message: "User deleted successfully",
	})
}

func bindRegisterRequest(c echo.Context) (RegisterRequest, error) {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return RegisterRequest{}, err
	}
	return req, nil
}

func bindLoginRequest(c echo.Context) (LoginRequest, error) {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return LoginRequest{}, err
	}
	return req, nil
}

func bindUpdateUserRequest(c echo.Context) (UpdateUserRequest, error) {
	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return UpdateUserRequest{}, err
	}
	return req, nil
}

func buildUserResponse(user domain.User) UserResponse {
	return UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func buildLoginResponse(token string, user domain.User) LoginResponse {
	return LoginResponse{
		Message: "Login successful",
		Token:   token,
		User:    buildUserResponse(user),
	}
}

func applyUserUpdate(user *domain.User, updateReq UpdateUserRequest) {
	user.Username = updateReq.Username
	user.Email = updateReq.Email
	user.FirstName = updateReq.FirstName
	user.LastName = updateReq.LastName
}
