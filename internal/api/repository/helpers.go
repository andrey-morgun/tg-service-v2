package repository

import (
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/rabbit"
	"io/ioutil"
	"net/http"
	"strconv"
)

func HandleHttpError(resp *http.Response) error {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errs.InternalError{Cause: "parsing body error"}
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errs.BadRequestError{Cause: string(data)}
	case http.StatusNotFound:
		return errs.NotFoundError{What: string(data)}
	case http.StatusUnauthorized:
		return errs.UnauthorizedError{Cause: string(data)}
	case http.StatusForbidden:
		return errs.ForbiddenError{Cause: string(data)}
	case http.StatusBadGateway:
		return errs.InternalError{Cause: strconv.Itoa(resp.StatusCode)}
	//	TODO: add StatusConflictError to lib
	//case http.StatusConflict:
	//	return
	default:
		return errs.InternalError{Cause: string(data)}
	}
}

func HandleBrokerError(resp rabbit.ResponseModel) error {
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errs.BadRequestError{Cause: string(resp.Payload)}
	case http.StatusNotFound:
		return errs.NotFoundError{What: string(resp.Payload)}
	case http.StatusUnauthorized:
		return errs.UnauthorizedError{Cause: string(resp.Payload)}
	case http.StatusForbidden:
		return errs.ForbiddenError{Cause: string(resp.Payload)}
	//	TODO: add StatusConflictError to lib
	//case http.StatusConflict:
	//	return
	default:
		return errs.InternalError{Cause: string(resp.Payload)}
	}
}
