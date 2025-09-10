package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ncostamagna/go-http-utils/meta"
	"github.com/ncostamagna/go-http-utils/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	ErrorRes struct {
		Error string `json:"error"`
	}

	UpdateReq struct {
		ID        *string `json:"id"`
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	DeleteRequest struct {
		ID string `json:"id"`
	}

	GetRequest struct {
		ID string `json:"id"`
	}
	GetAllReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Limit     int
		Page      int
	}

	Response struct {
		Status int         `json:"status"`
		Err    string      `json:"err,omitempty"`
		Data   interface{} `json:"data,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndpoints(s Service) Endpoints {

	return Endpoints{
		Create: makeCreateEndpoint(s),
		Update: makeUpdateEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Get:    makeGetEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}

}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(DeleteRequest)

		if req.ID == "" {

			return nil, response.BadRequest("Id cannot be null")
		}
		id, err := s.Delete(ctx, req.ID)
		if err != nil {
			log.Println(err)
			if errors.As(err, &ErrUserNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return nil, response.OK("User deleted successfully", id, nil)

	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest) //conversion
		fmt.Println(&req)
		if req.FirstName == "" {

			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}
		if req.LastName == "" {

			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		usr, err := s.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("Success", usr, nil), nil

	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)
		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		count, err := s.Count(ctx, filters)
		log.Println(count)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		defaultLimit := os.Getenv("PAGINATION_PER_PAGE_DEFAULT")
		meta, err := meta.New(count, req.Page, req.Limit, defaultLimit)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		usr, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())

		if err != nil {
			if errors.As(err, &ErrUserNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("successfull user fetch", usr, meta), nil

	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(UpdateReq)

		usr, err := s.Update(ctx, &req)
		if err != nil {
			if errors.As(err, &ErrUserNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("User updated successfully", usr, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("get user")
		req := request.(GetRequest)

		if req.ID == "" {
			return response.BadRequest("Invalid request data"), nil
		}
		usr, err := s.GetById(ctx, req.ID)
		if err != nil {
			return response.InternalServerError(err.Error()), nil
		}

		return response.OK("successfull user fetch", usr, nil), nil

	}
}
