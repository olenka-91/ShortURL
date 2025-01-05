package main

import (
	"fmt"
	"io"
	"net/http"
)

func shortURL(longURL []byte) ([]byte, error) {

	return longURL, nil
}

func longURL(shortURL string) (string, error) {

	return shortURL, nil
}

func postShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		fmt.Println("запрос Post")
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		resBody, err := shortURL(reqBody)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Content-Length", "30")
		res.Write([]byte(resBody))
		res.WriteHeader(http.StatusCreated)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func getUnShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Println("запрос Get")
		reqId := req.URL.Query().Get("id")

		resText, err := longURL(reqId)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Location", resText)

		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postShortURL)
	mux.HandleFunc(`/{id}`, getUnShortURL)

	fmt.Println("Сервер запускается")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
	fmt.Println("Сервер запущен")
}
