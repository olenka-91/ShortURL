package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"encoding/json"

	"github.com/olenka-91/shorturl/internal/models"
)

func FromContext(ctx context.Context) (int, bool) {
	uid, ok := ctx.Value(models.UserKey).(int)
	return uid, ok && uid != 0
}

func (h *Handler) UserURLs(res http.ResponseWriter, req *http.Request) {
	uid, ok := FromContext(req.Context())
	if !ok {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	var URLs []models.URLsForUser

	URLs, err := h.services.ListURLsByUser(req.Context(), uid)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(URLs) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(URLs)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	//res.Header().Set("Content-Length", "30")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Handler) PostShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		ctx := req.Context()
		uid, ok := FromContext(ctx)
		if !ok {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

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
		resBody, err := h.services.ShortURL(ctx, reqBody, uid)
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
		uid, ok := FromContext(req.Context())
		if !ok {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
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
		ctx := req.Context()
		resBody, err := h.services.ShortURL(ctx, []byte(input.URL), uid)
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
		ctx := req.Context()
		batchSize := 2
		batch := make([]models.BatchInput, 0, batchSize)
		batchOutput := make([]models.BatchOutput, 0, batchSize)
		var output []models.BatchOutput
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()

		uid, ok := FromContext(req.Context())
		if !ok {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

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
				if batchOutput, err = h.services.PostURLBatch(ctx, batch, uid); err != nil {
					res.WriteHeader(http.StatusInternalServerError)
					return
				}
				batch = batch[:0]
				output = append(output, batchOutput...)
			}
		}

		if batchOutput, err = h.services.PostURLBatch(ctx, batch, uid); err != nil {
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

		uid, ok := FromContext(req.Context())
		if !ok {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := req.Context()
		resText, err := h.services.LongURL(ctx, reqId, uid)
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
