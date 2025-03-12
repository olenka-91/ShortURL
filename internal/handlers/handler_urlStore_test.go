package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/olenka-91/shorturl/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func TestPostShortURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		method  string
		request string
		body    io.Reader
		want    want
	}{
		{
			name:    "positive test #1",
			method:  http.MethodPost,
			request: "/",
			body:    strings.NewReader("yandex.ru"),
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name:    "no body test #2",
			method:  http.MethodPost,
			body:    io.NopCloser(errReader{}),
			request: "/",
			want: want{
				code:        http.StatusInternalServerError,
				contentType: "",
			},
		},
		{
			name:    "Get #3",
			method:  http.MethodGet,
			request: "/",
			body:    strings.NewReader(""),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
	}

	serv := service.NewService("")
	handl := NewHandler(serv)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, test.body)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			handl.PostShortURL(w, request)

			res := w.Result()
			// проверяем код ответа и тип
			assert.Equal(t, res.StatusCode, test.want.code)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)

			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestPostShortURLJSON(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	input, _ := json.Marshal(ShortURLInput{URL: "https://practicum.yandex.ru"})

	tests := []struct {
		name    string
		method  string
		request string
		body    io.Reader
		want    want
	}{
		{
			name:    "positive test #1",
			method:  http.MethodPost,
			request: "/api/shorten",
			body:    strings.NewReader(string(input)),
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name:    "Other type test #2",
			method:  http.MethodPost,
			body:    strings.NewReader("abc:123"),
			request: "/api/shorten",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name:    "Get #3",
			method:  http.MethodGet,
			request: "/api/shorten",
			body:    strings.NewReader(""),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
		},
	}

	serv := service.NewService("")
	handl := NewHandler(serv)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, test.body)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			handl.PostShortURLJSON(w, request)

			res := w.Result()
			// проверяем код ответа и тип
			assert.Equal(t, res.StatusCode, test.want.code)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)

			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestGetUnShortURL(t *testing.T) {
	serv := service.NewService("")
	handl := NewHandler(serv)
	location := "yandex.ru"
	urlID, _ := serv.ShortURL([]byte(location))

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		method  string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			method:  http.MethodGet,
			request: "/" + urlID,
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain",
				response:    location,
			},
		},
		{
			name:    "No such adress #2",
			method:  http.MethodGet,
			request: "/asdsdkiopkjhhggkjhklkkljhfgytrkulkj",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
				response:    "",
			},
		},
		{
			name:    "Post #3",
			method:  http.MethodPost,
			request: "/",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
				response:    "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, strings.NewReader(""))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			handl.GetUnShortURL(w, request)

			res := w.Result()
			// проверяем код ответа и тип
			assert.Equal(t, res.StatusCode, test.want.code)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
			assert.Equal(t, res.Header.Get("Location"), test.want.response)

		})
	}
}
