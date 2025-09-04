package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/juanrojas09/gocourse_meta/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)
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

	Response struct {
		Status int            `json:"status"`
		Err    string         `json:"err,omitempty"`
		Data   interface{}    `json:"data,omitempty"`
		Meta   *meta.Metadata `json:"meta,omitempty"`
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
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		fmt.Println(id)

		if id == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Error al eliminar user"})
			return
		}
		id, err := s.Delete(id)
		if err != nil {
			log.Println(err)
		}

		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   id,
		})
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest
		fmt.Println(&req)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		usr, err := s.Create(req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
		}
		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   usr,
		})
	}
}
func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()

		filters := Filters{
			FirstName: v.Get("first_name"),
			LastName:  v.Get("last_name"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		count, err := s.Count(filters)
		log.Println(count)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorRes{err.Error()})
		}
		defaultLimit := os.Getenv("PAGINATION_PER_PAGE_DEFAULT")
		meta, err := meta.New(count, page, limit, defaultLimit)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorRes{err.Error()})
		}
		usr, err := s.GetAll(filters, meta.Offset(), meta.Limit())

		if err != nil {
			json.NewEncoder(w).Encode(ErrorRes{err.Error()})
		}
		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   usr,
			Meta:   meta,
		})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		usr, err := s.Update(&req)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
		}
		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   usr,
		})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get user")

		vars := mux.Vars(r)
		id := vars["id"]

		fmt.Println(id)
		if id == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Error al obtener user"})
			return
		}
		usr, err := s.GetById(id)
		if err != nil {
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
		}

		json.NewEncoder(w).Encode(&Response{
			Status: 200,
			Data:   usr,
		})

	}
}
