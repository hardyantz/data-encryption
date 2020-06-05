package helpers

import "github.com/labstack/echo"

type JsonSchema struct {
	Code    int         `json:"statusCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

type Meta map[string]interface{}

func HTTPResponse(c echo.Context, statusCode int, message string, data interface{}, meta Meta) error {
	res := JsonSchema{
		Code:    statusCode,
		Message: message,
		Data:    data,
		Meta:    meta,
	}

	return c.JSON(statusCode, res)
}
