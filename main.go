package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}

type Stringer interface {
	String() string
}

type HttpMethod string

func (h HttpMethod) String() string {
	return string(h)
}

const (
	GET     HttpMethod = http.MethodGet
	DELETE  HttpMethod = http.MethodDelete
	POST    HttpMethod = http.MethodPost
	PUT     HttpMethod = http.MethodPut
	OPTIONS HttpMethod = http.MethodOptions
)

type HogeFuga string

func (h HogeFuga) String() string {
	return string(h)
}

const (
	Hoge HogeFuga = "hoge"
	Fuga HogeFuga = "fuga"
)

func toResponse[T Stringer](value T, w http.ResponseWriter) {
	marshal, _ := json.Marshal(Response{Message: value.String()})
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshal)
}

func main() {
	port := ":" + os.Getenv("PORT")
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		log.Println(r.URL.Query())

		switch r.Method {
		case http.MethodGet:
			res := GET.String() + ":" + time.Now().String()
			if r.URL.Query().Get("q") != "" {
				res = res + ":" + r.URL.Query().Get("q")
			}
			toResponse(HogeFuga(res), w)
		case http.MethodDelete:
			toResponse(DELETE, w)
		case http.MethodPost:
			req := &Request{}
			json.NewDecoder(r.Body).Decode(req)
			log.Println(req.Message)
			switch req.Message {
			case "hoge":
				toResponse(Hoge, w)
			case "fuga":
				toResponse(Fuga, w)
			default:
				// return 500 error
				w.WriteHeader(http.StatusInternalServerError)
				//toResponse(POST, w)
			}
		case http.MethodPut:
			toResponse(PUT, w)
		case http.MethodOptions:
			toResponse(OPTIONS, w)
		default:
			fmt.Printf("Unknown method: %s", r.Method)
		}
	})
	http.HandleFunc("/notfound", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusNotFound)
			toResponse(HogeFuga(GET.String()+" Error"), w)
		}
	})
	http.HandleFunc("/internal", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusInternalServerError)
			toResponse(HogeFuga(GET.String()+" Error"), w)
		}
	})

	http.ListenAndServe(port, http.DefaultServeMux)
}
