package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/constants"
	"boilerplate-api/dtos"
	"boilerplate-api/models"
	"boilerplate-api/paginations"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUsers(pagination paginations.UserPagination) ([]dtos.GetUserResponse, int64, error) {
	args := m.Called(pagination)
	return args.Get(0).([]dtos.GetUserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) GetOneUser(Id string) (dtos.GetUserResponse, error) {
	args := m.Called(Id)
	return args.Get(0).(dtos.GetUserResponse), args.Error(1)
}

func (m *MockUserService) CreateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserService) WithTrx(trxHandle *gorm.DB) services.UserService {
	args := m.Called(trxHandle)
	return args.Get(0).(services.UserService)
}
func (m *MockUserService) GetOneUserWithEmail(Email string) (models.User, error) {
	args := m.Called(Email)
	return args.Get(0).(models.User), args.Error(1)
}
func (m *MockUserService) GetOneUserWithPhone(Phone string) (models.User, error) {
	args := m.Called(Phone)
	return args.Get(0).(models.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockUserService := &MockUserService{}
	userController := UserController{
		userService: mockUserService,
	}
	user := models.User{
		Email:    "example@email.com",
		FullName: "John Doe",
		Phone:    "123-456-7890",
		Gender:   "Male",
		Password: "secretPassword",
	}

	reqData := dtos.CreateUserRequestData{
		User:            user,
		ConfirmPassword: "password",
	}
	mockDBTransaction := &gorm.DB{}
	ctx := context.WithValue(context.Background(), "db_trx", mockDBTransaction)

	reqBody, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req.WithContext(ctx)

	expectedError := error(nil)

	mockUserService.On("GetOneUserWithEmail", reqData.Email).
		Return(models.User{}, errors.New("User not found"))
	mockUserService.On("GetOneUserWithPhone", reqData.Phone).
		Return(models.User{}, errors.New("User not found"))
	mockUserService.On("WithTrx", mock.AnythingOfType("*gorm.DB")).Return(mockUserService)
	mockUserService.On("CreateUser", reqData.User).
		Return(expectedError)

	userController.CreateUser(c)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
	expectedResponseBody := `{"message": "User Created Successfully"}`
	assert.JSONEq(t, expectedResponseBody, w.Body.String(), "Expected response body")
	mockUserService.AssertCalled(t, "GetOneUserWithEmail", reqData.Email)
	mockUserService.AssertCalled(t, "GetOneUserWithPhone", reqData.Phone)
	mockUserService.AssertCalled(t, "WithTrx", mock.AnythingOfType("*gorm.DB"))
	mockUserService.AssertCalled(t, "CreateUser", reqData.User)
}

func TestGetAllUsers(t *testing.T) {
	mockUserService := new(MockUserService)
	userController := UserController{
		userService: mockUserService,
	}

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	user := models.User{
		Email:    "example@email.com",
		FullName: "John Doe",
		Phone:    "123-456-7890",
		Gender:   "Male",
		Password: "secretPassword",
	}
	expectedUser := dtos.GetUserResponse{
		User:     user,
		Password: "secretPassword",
	}
	expectedUsers := []dtos.GetUserResponse{expectedUser}
	expectedCount := int64(len(expectedUsers))
	expectedError := error(nil)

	mockUserService.On("GetAllUsers", mock.AnythingOfType("paginations.UserPagination")).
		Return(expectedUsers, expectedCount, expectedError)

	userController.GetAllUsers(c)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
	expectedResponseBody := `{"data":[{"id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null,"email":"example@email.com","full_name":"John Doe","phone":"123-456-7890","gender":"Male","password":"secretPassword","Password":"secretPassword"}],"count":1}`
	assert.Equal(t, expectedResponseBody, w.Body.String(), "Expected response body")
	mockUserService.AssertCalled(t, "GetAllUsers", mock.AnythingOfType("paginations.UserPagination"))
}

func TestGetUserProfile(t *testing.T) {
	mockUserService := new(MockUserService)
	userController := UserController{
		userService: mockUserService,
	}

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	user := models.User{
		Email:    "example@email.com",
		FullName: "John Doe",
		Phone:    "123-456-7890",
		Gender:   "Male",
		Password: "secretPassword",
	}
	expectedUser := dtos.GetUserResponse{
		User:     user,
		Password: "password",
	}
	expectedError := error(nil)
	mockUserService.On("GetOneUser", mock.AnythingOfType("string")).
		Return(expectedUser, expectedError)

	c.Set(constants.UserID, "123")
	userController.GetUserProfile(c)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
	expectedResponseBody := `{
		"data": {
			"Password": "password",
			"created_at": "0001-01-01T00:00:00Z",
			"deleted_at": null,
			"email": "example@email.com",
			"full_name": "John Doe",
			"gender": "Male",
			"id": 0,
			"password": "secretPassword",
			"phone": "123-456-7890",
			"updated_at": "0001-01-01T00:00:00Z"
		}
	}`
	assert.JSONEq(t, expectedResponseBody, w.Body.String(), "Expected response body")
	mockUserService.AssertCalled(t, "GetOneUser", "123")
}
