package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	// Хранение сокращённых URL-адресов
	urlStore = make(map[string]string)
	mu       sync.Mutex // Защита от конкуренции
)

func shortURL(longURL []byte) (string, error) {
	bytes := make([]byte, 6) // 6 байт = 8 символов в base64
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	// Кодируем байты в строку base64

	shortURL64 := base64.URLEncoding.EncodeToString(bytes)

	mu.Lock()
	defer mu.Unlock()
	urlStore[shortURL64] = string(longURL)

	fmt.Println("longURL=", string(longURL))
	fmt.Println("shortURL=", shortURL64)
	return shortURL64, nil

}

func longURL(shortURL string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	longURL, ok := urlStore[shortURL]
	if !ok {
		return "", fmt.Errorf("Not found")
	}
	fmt.Println("shortURL=", shortURL)
	fmt.Println("longURL=", longURL)
	return longURL, nil
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
		//res.WriteHeader(http.StatusCreated)
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
