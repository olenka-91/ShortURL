package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func (h *Handler) PostShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		fmt.Println("запрос Post")
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(reqBody) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		resBody, err := h.services.ShortURL(reqBody)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Content-Length", "30")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(resBody))

	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) GetUnShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Println("запрос Get")
		//reqId := req.URL.Query().Get("id")
		reqId := req.URL.Path[1:]

		resText, err := h.services.LongURL(reqId)
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
