package errors

import (
	"fmt"
	"net/http"
)

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct {
	Detail string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Detail)
}

func (e NotFoundError) StatusCode() int {
	return http.StatusNotFound // 404
}

// BadRequestError represents a 400 Bad Request error.
type BadRequestError struct {
	Detail string
}

func (e BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %s", e.Detail)
}

func (e BadRequestError) StatusCode() int {
	return http.StatusBadRequest // 400
}

// UnauthorizedError represents a 401 Unauthorized error.
type UnauthorizedError struct {
	Detail string
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Detail)
}

func (e UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized // 401
}

// InternalServerError represents a 500 Internal Server Error.
type InternalServerError struct {
	Detail string
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %s", e.Detail)
}

func (e InternalServerError) StatusCode() int {
	return http.StatusInternalServerError // 500
}
