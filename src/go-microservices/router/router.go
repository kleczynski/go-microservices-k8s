package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kleczynski/go-microservices-k8s/handlers"
)

func configureRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", handlers.HealthHanler)
	r.HandleFunc("/details", handlers.DetailsHandler)
	r.HandleFunc("/", handlers.RootHandler)
	r.HandleFunc("/api", handlers.APIListHandler)
	configureAPIRoutes(r)
	return r
}

func configureAPIRoutes(r *mux.Router) {
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/create", handlers.Create).Methods("POST")
	api.HandleFunc("/read", handlers.Read).Methods("GET")
	api.HandleFunc("/update", handlers.Update).Methods("PUT")
	api.HandleFunc("/delete/{name}", handlers.Delete).Methods("DELETE")

}

func StartServer(port string) error {
	r := configureRoutes()
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s\n", err)
	}
	return nil
}
