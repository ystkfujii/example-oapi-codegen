package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	mw "github.com/oapi-codegen/echo-middleware"

	api "github.com/ystkfujii/example-oapi-codegen/openapi"
)

func setupServer() *echo.Echo {
	e := echo.New()

	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}
	
	// Clear servers to avoid Host validation issues in tests
	swagger.Servers = nil

	validatorOptions := mw.Options{
		ErrorHandler: customErrorHandler,
	}
	e.Use(mw.OapiRequestValidatorWithOptions(swagger, &validatorOptions))

	server := NewServer()
	api.RegisterHandlers(e, server)

	return e
}

func TestGetUsersEmpty(t *testing.T) {
	e := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, rec.Code, rec.Body.String())
		return
	}

	var users []api.User
	if err := json.Unmarshal(rec.Body.Bytes(), &users); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("Expected empty users list, got %d users", len(users))
	}
}

func TestPostUser(t *testing.T) {
	tests := []struct {
		name           string
		user           api.NewUser
		expectedStatus int
		expectedError  string
		checkUser      bool
	}{
		{
			name: "valid user",
			user: api.NewUser{
				Name: api.Name{
					First: "Alice",
					Last:  "Smith",
				},
				Age: 10,
			},
			expectedStatus: http.StatusCreated,
			checkUser:      true,
		},
		{
			name: "invalid age - too high",
			user: api.NewUser{
				Name: api.Name{
					First: "Alice",
					Last:  "Smith",
				},
				Age: 16, // Maximum is 15
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "\"number must be at most 15\"\n",
		},
		{
			name: "invalid first name - contains numbers",
			user: api.NewUser{
				Name: api.Name{
					First: "Alice123", // Should only contain alphabets
					Last:  "Smith",
				},
				Age: 10,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "\"string doesn't match the regular expression \\\"^[a-zA-Z]+$\\\"\"\n",
		},
		{
			name: "invalid first name - empty",
			user: api.NewUser{
				Name: api.Name{
					First: "", // Empty name violates minLength: 1
					Last:  "Smith",
				},
				Age: 10,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "\"minimum string length is 1\"\n",
		},
		{
			name: "invalid first name - too long",
			user: api.NewUser{
				Name: api.Name{
					First: strings.Repeat("A", 101), // Longer than 100 characters
					Last:  "Smith",
				},
				Age: 10,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "\"maximum string length is 100\"\n",
		},
		{
			name: "missing last name",
			user: api.NewUser{
				Name: api.Name{
					First: "Alice",
					// Last is required but missing
				},
				Age: 10,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "\"minimum string length is 1\"\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupServer()
			body, _ := json.Marshal(tt.user)
			
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedError != "" {
				responseBody := rec.Body.String()
				if responseBody != tt.expectedError {
					t.Errorf("Expected response '%s', got: %s", tt.expectedError, responseBody)
				}
			}

			if tt.checkUser {
				var user api.User
				if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if user.Name.First != tt.user.Name.First || user.Name.Last != tt.user.Name.Last || user.Age != tt.user.Age || user.Id != 1 {
					t.Errorf("Expected user {Id: 1, Name: %+v, Age: %d}, got %+v", tt.user.Name, tt.user.Age, user)
				}
			}
		})
	}
}

func TestGetUserById(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		createUser     bool
		expectedStatus int
		expectedUser   *api.User
	}{
		{
			name:           "get existing user",
			userID:         "1",
			createUser:     true,
			expectedStatus: http.StatusOK,
			expectedUser:   &api.User{Id: 1, Name: api.Name{First: "Bob", Last: "Johnson"}, Age: 5},
		},
		{
			name:           "get non-existing user",
			userID:         "999",
			createUser:     false,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupServer()

			if tt.createUser {
				// First create a user
				newUser := api.NewUser{
					Name: api.Name{
						First: "Bob",
						Last:  "Johnson",
					},
					Age: 5,
				}
				body, _ := json.Marshal(newUser)
				
				req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
			}

			// Then get the user by ID
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedUser != nil {
				var user api.User
				if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if user.Name.First != tt.expectedUser.Name.First || user.Name.Last != tt.expectedUser.Name.Last || user.Age != tt.expectedUser.Age || user.Id != tt.expectedUser.Id {
					t.Errorf("Expected user %+v, got %+v", tt.expectedUser, user)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		createUser     bool
		expectedStatus int
		verifyDeleted  bool
	}{
		{
			name:           "delete existing user",
			userID:         "1",
			createUser:     true,
			expectedStatus: http.StatusNoContent,
			verifyDeleted:  true,
		},
		{
			name:           "delete non-existing user",
			userID:         "999",
			createUser:     false,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupServer()

			if tt.createUser {
				// First create a user
				newUser := api.NewUser{
					Name: api.Name{
						First: "Charlie",
						Last:  "Brown",
					},
					Age: 8,
				}
				body, _ := json.Marshal(newUser)
				
				req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
			}

			// Then delete the user
			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.verifyDeleted {
				// Verify user is deleted
				req = httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
				rec = httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				if rec.Code != http.StatusNotFound {
					t.Errorf("Expected status %d after deletion, got %d", http.StatusNotFound, rec.Code)
				}
			}
		})
	}
}