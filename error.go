package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	api "github.com/ystkfujii/example-oapi-codegen/openapi"
)

// Error response functions
func ReturnUserNotFoundError(ctx echo.Context, userID int) error {
	details := fmt.Sprintf("User with ID %d does not exist", userID)
	errorResp := api.Error{
		Code:    http.StatusNotFound,
		Message: "User not found",
		Details: &details,
	}
	return ctx.JSON(errorResp.Code, errorResp)
}

func ReturnInvalidRequestBodyError(ctx echo.Context, err error) error {
	details := err.Error()
	errorResp := api.Error{
		Code:    http.StatusBadRequest,
		Message: "Invalid request body",
		Details: &details,
	}
	return ctx.JSON(errorResp.Code, errorResp)
}

func ReturnInvalidUserDataError(ctx echo.Context, err error) error {
	details := err.Error()
	errorResp := api.Error{
		Code:    http.StatusBadRequest,
		Message: "Invalid user data",
		Details: &details,
	}
	return ctx.JSON(errorResp.Code, errorResp)
}

func ReturnInvalidUserIDError(ctx echo.Context, err error) error {
	details := fmt.Sprintf("Invalid format for parameter id: %s", err)
	errorResp := api.Error{
		Code:    http.StatusBadRequest,
		Message: "Invalid user ID",
		Details: &details,
	}
	return ctx.JSON(errorResp.Code, errorResp)
}

func ReturnFormatError(ctx echo.Context, reson string) error {
	errorResp := api.Error{
		Code:    http.StatusBadRequest,
		Message: "Invalid format",
		Details: &reson,
	}
	return ctx.JSON(errorResp.Code, errorResp)
}
