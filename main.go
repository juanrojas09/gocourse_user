package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/juanrojas09/gocourse_user/bootstrap"
	"github.com/juanrojas09/gocourse_user/users"
)

func main() {
	//cargar envs
	err := godotenv.Load()
	errorHandler(err)
	//instanciar router con mux
	router := mux.NewRouter()

	//inyectar servicios del bootstrap
	db, err := bootstrap.InitDb()
	errorHandler(err)
	log := bootstrap.InitLogger()

	userRepository := users.NewRepository(db, log)
	userService := users.NewService(userRepository, log)
	userEndpoints := users.MakeEndpoints(userService)

	router.HandleFunc("/users/{id}", userEndpoints.Get).Methods(http.MethodGet)
	router.HandleFunc("/users", userEndpoints.GetAll).Methods(http.MethodGet)
	router.HandleFunc("/users", userEndpoints.Create).Methods(http.MethodPost)
	router.HandleFunc("/users", userEndpoints.Update).Methods(http.MethodPatch)
	router.HandleFunc("/users/{id}", userEndpoints.Delete).Methods(http.MethodDelete)

	srv := http.Server{
		Handler: router,
		Addr:    os.Getenv("API_URL") + ":" + os.Getenv("API_PORT"),
	}

	if err := http.ListenAndServe(srv.Addr, srv.Handler); err != nil {
		errorHandler(err)
	}

}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
