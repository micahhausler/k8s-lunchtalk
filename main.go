package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"io/ioutil"

	"github.com/micahhausler/k8s-lunchtalk/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type APIResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Time    string `json:"time"`
}

func RootListener(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := APIResponse{
		Message: "Hello Chadevs!",
		Time:    time.Now().String(),
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(response)
}

func Proxy(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := APIResponse{
		Time:  time.Now().String(),
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	req.ParseForm()
	q := req.Form.Get("q")

	resp, err := http.Get(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Error = err.Error()
		encoder.Encode(response)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	response.Message = string(body)

	encoder.Encode(response)
}


func Error(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := APIResponse{
		Error: "Error!",
		Time:  time.Now().String(),
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(response)
}

func main() {
	listenHostPort := "0.0.0.0:3000"

	r := http.NewServeMux()
	r.HandleFunc("/", RootListener)
	r.HandleFunc("/proxy/", Proxy)
	r.HandleFunc("/error/", Error)

	handler := middlewares.Apply(
		r,
		middlewares.InstrumentRoute(),
		middlewares.Logging(),
	)

	http.Handle("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %s", listenHostPort)
	log.Fatal(http.ListenAndServe(listenHostPort, nil))
}
