package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"receptor/backend/internal/service"
)

// RespondError maps service/repository errors to HTTP (handler stays free of business rules).
func RespondError(c *gin.Context, err error) {
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		JSONError(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrEmailAlreadyExists):
		JSONError(c, http.StatusConflict, err.Error())
	case errors.Is(err, pgx.ErrNoRows):
		JSONError(c, http.StatusNotFound, "not found")
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		JSONError(c, http.StatusConflict, "conflict")
	default:
		JSONError(c, http.StatusInternalServerError, "internal server error")
	}
}
