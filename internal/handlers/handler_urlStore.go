package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"encoding/json"

	"github.com/olenka-91/shorturl/internal/models"
)

func (h *Handler) PostShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		ctx := context.Background()
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(reqBody) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var MyErr *models.DBError
		resBody, err := h.services.ShortURL(ctx, reqBody)
		if err != nil {
			if errors.As(err, &MyErr) {
				res.WriteHeader(http.StatusConflict)
				return
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

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

		var MyErr *models.DBError
		ctx := context.Background()
		resBody, err := h.services.ShortURL(ctx, []byte(input.URL))
		if err != nil {
			if errors.As(err, &MyErr) {
				res.WriteHeader(http.StatusConflict)
				return
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
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

func (h *Handler) PostShortURLJSONBatch(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var err error
		ctx := context.Background()
		batchSize := 2
		batch := make([]models.BatchInput, 0, batchSize)
		batchOutput := make([]models.BatchOutput, 0, batchSize)
		var output []models.BatchOutput
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()

		if _, err := decoder.Token(); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		for decoder.More() {
			var input models.BatchInput
			if err := decoder.Decode(&input); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			// if len(input) == 0 {
			// 	res.WriteHeader(http.StatusBadRequest)
			// 	return
			// }

			batch = append(batch, input)
			if len(batch) == batchSize {
				if batchOutput, err = h.services.PostURLBatch(ctx, batch); err != nil {
					res.WriteHeader(http.StatusInternalServerError)
					return
				}
				batch = batch[:0]
				output = append(output, batchOutput...)
			}
		}

		if batchOutput, err = h.services.PostURLBatch(ctx, batch); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		output = append(output, batchOutput...)

		if _, err := decoder.Token(); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

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
		//reqId := req.URL.Query().Get("id")
		reqId := req.URL.Path[1:]

		ctx := context.Background()
		resText, err := h.services.LongURL(ctx, reqId)
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
