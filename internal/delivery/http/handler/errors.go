package handler

import (
	"errors"
	"net/http"

	httpinternal "money_tracker/internal/delivery/http/internal"
	"money_tracker/internal/domain"
)

func httpStatusByError(err error) (int, string) {
	var entityAlreadyExistsErr *domain.EntityAlreadyExistsError
	if errors.As(err, &entityAlreadyExistsErr) {
		return http.StatusConflict, entityAlreadyExistsErr.Error()
	}

	var invalidArgumentErr *domain.InvalidArgumentError
	if errors.As(err, &invalidArgumentErr) {
		return http.StatusBadRequest, invalidArgumentErr.Error()
	}

	var entityNotFoundErr *domain.EntityNotFoundError
	if errors.As(err, &entityNotFoundErr) {
		return http.StatusNotFound, entityNotFoundErr.Error()
	}

	var authorizationErr *domain.AuthorizationError
	if errors.As(err, &authorizationErr) {
		return http.StatusForbidden, authorizationErr.Error()
	}

	return http.StatusInternalServerError, "Internal server error"
}

func processError(w http.ResponseWriter, err error) {
	status, msg := httpStatusByError(err)
	httpinternal.RespondWithError(w, status, msg)
}
