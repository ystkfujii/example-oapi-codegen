package main

import (
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	mw "github.com/oapi-codegen/echo-middleware"

	api "github.com/ystkfujii/example-oapi-codegen/openapi"
)

func rootError(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}

func customErrorHandler(ctx echo.Context, err *echo.HTTPError) error {
	e := rootError(err)
	var sErr *openapi3.SchemaError
	if errors.As(e, &sErr) {
		return ReturnFormatError(ctx, sErr.Reason)
	}
	return ReturnFormatError(ctx, "not according to the api specification")
}

func main() {
	e := echo.New()

	swagger, err := api.GetSwagger()
	if err != nil {
		e.Logger.Fatal(err)
	}

	validatorOptions := mw.Options{
		ErrorHandler: customErrorHandler,
	}
	e.Use(mw.OapiRequestValidatorWithOptions(swagger, &validatorOptions))

	server := NewServer()
	api.RegisterHandlers(e, server)

	e.Logger.Fatal(e.Start(":8080"))
}
