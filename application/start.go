package application

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const bpmnFolderLocation = "./bpmn/"

func Start() {
	environment := getEnvironment()
	h := Handler{newHandler(environment.zeebeConfig)}
	router := mux.NewRouter()

	//define routes
	router.HandleFunc("/new/resource", h.DeployResource).Methods(http.MethodPost)
	router.HandleFunc("/new/instance", h.CreateInstance).Methods(http.MethodPost)

	//starting server
	log.Fatal(http.ListenAndServe(environment.serviceConfig.address, router))
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
