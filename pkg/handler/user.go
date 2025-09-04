package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/juanrojas09/gocourse_user/internal/users"
)

func NewUserHttpServer(ctx context.Context, endpoints users.Endpoints) http.Handler {
	r := mux.NewRouter()

	r.Handle("/users", httptransport.NewServer(
		(endpoint.Endpoint)(endpoints.Create),
		decodeCreateUser, encodeCreateUser,
	)).Methods(http.MethodPost)

	return r
}

// metodo para enviar directamente el objeto serializado al controlador
func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {

	var req users.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil

}

func encodeCreateUser(_ context.Context, w http.ResponseWriter, resp interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(resp)

}
