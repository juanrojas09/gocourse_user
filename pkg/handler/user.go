package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/juanrojas09/gocourse_user/internal/users"
	"github.com/ncostamagna/go-http-utils/response"
)

func NewUserHttpServer(ctx context.Context, endpoints users.Endpoints) http.Handler {
	r := mux.NewRouter()
	//metodo predeterminado cuando falla algun endpoint o un error no manejado.
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/users", httptransport.NewServer(
		(endpoint.Endpoint)(endpoints.Create),
		decodeCreateUser, encodeCreateUser, opts...,
	)).Methods(http.MethodPost)

	r.Handle("/users/{id}", httptransport.NewServer(
		(endpoint.Endpoint)(endpoints.Get),
		decodeGetUser, encodeGetUser, opts...,
	)).Methods(http.MethodGet)

	r.Handle("/users", httptransport.NewServer(
		(endpoint.Endpoint)(endpoints.GetAll),
		decodeGetAllUser, encodeGetUser, opts...,
	)).Methods(http.MethodGet)

	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateUser, encodeUpdateUser))

	r.Handle("/users/{id}", httptransport.NewServer(endpoint.Endpoint(endpoints.Delete),
		decodeDeleteUser, encodeDeleteUser))

	return r
}

// metodo para enviar directamente el objeto serializado al controlador
func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {
	var req users.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	return req, nil

}

func encodeCreateUser(_ context.Context, w http.ResponseWriter, resp interface{}) (err error) {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(resp)

}

func decodeGetUser(_ context.Context, r *http.Request) (interface{}, error) {

	pathVar := mux.Vars(r)
	id := pathVar["id"]

	req := users.GetRequest{
		ID: id,
	}

	return req, nil

}

func encodeGetUser(_ context.Context, w http.ResponseWriter, resp interface{}) (err error) {
	r := resp.(response.Response)
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(resp)

}

func decodeGetAllUser(_ context.Context, r *http.Request) (interface{}, error) {

	queryParams := r.URL.Query()
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	page, _ := strconv.Atoi(queryParams.Get("page"))

	req := users.GetAllReq{
		FirstName: queryParams.Get("first_name"),
		LastName:  queryParams.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}

	return req, nil

}

func decodeDeleteUser(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	id := params["id"]

	req := users.DeleteRequest{
		ID: id,
	}
	return req, nil

}

func encodeDeleteUser(_ context.Context, w http.ResponseWriter, resp interface{}) (error error) {
	r := resp.(response.Response)
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(w)
}

func decodeUpdateUser(_ context.Context, r *http.Request) (interface{}, error) {

	var req users.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest("Error decoding update request, review body")
	}

	return req, nil

}

func encodeUpdateUser(_ context.Context, w http.ResponseWriter, resp interface{}) error {

	r := resp.(response.Response)

	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(resp)

}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
