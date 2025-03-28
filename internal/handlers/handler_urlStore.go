package handlers

import (
	"fmt"
	"io"
	"net/http"

	"encoding/json"
)

func (h *Handler) PostShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
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
		//res.Header().Set("Content-Length", "30")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(resBody))

	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

type ShortURLInput struct {
	URL string `json:"url"`
}

type ShortURLOutput struct {
	Result string `json:"result"`
}

func (h *Handler) PostShortURLJSON(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		input := ShortURLInput{}
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(input.URL) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		resBody, err := h.services.ShortURL([]byte(input.URL))
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		output := ShortURLOutput{Result: resBody}
		resp, err := json.Marshal(output)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		//res.Header().Set("Content-Length", "30")
		res.WriteHeader(http.StatusCreated)
		_, err = res.Write(resp)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

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

func (h *Handler) GetDBPing(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		err := h.services.PingDB()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		} else {
			res.WriteHeader(http.StatusOK)
		}

	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}
