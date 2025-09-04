package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/juanrojas09/gocourse_user/internal/users"
	"github.com/juanrojas09/gocourse_user/pkg/bootstrap"
	"github.com/juanrojas09/gocourse_user/pkg/handler"
)

func main() {
	//cargar envs
	err := godotenv.Load()
	errorHandler(err)

	//inyectar servicios del bootstrap
	db, err := bootstrap.InitDb()
	errorHandler(err)
	log := bootstrap.InitLogger()

	ctx := context.Background()
	userRepository := users.NewRepository(db, log)
	userService := users.NewService(userRepository, log)
	userEndpoints := users.MakeEndpoints(userService)
	h := handler.NewUserHttpServer(ctx, userEndpoints)
	address := os.Getenv("API_URL") + ":" + os.Getenv("API_PORT")
	srv := http.Server{
		Handler: accessControl(h),
		Addr:    address,
	}

	errCh := make(chan error)
	go func() {
		log.Println("listening in", address)
		errCh <- srv.ListenAndServe() // si hay un error, lo asigna
	}()
	err = <-errCh // el canal se va a quedar esperando hasta que reciba algo, cuando lo recibe continua la ejecucion
	if err != nil {
		errorHandler(err)
	}

}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", "GET<POST,PATCH,OPTIONS,HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
