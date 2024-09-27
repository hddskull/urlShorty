package root

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//test post with empty body
//test post with normal url
//test get with porper url
//test get with wrong url

const (
	localhost = "http://localhost:8080/"
)

func TestRootHandler(t *testing.T) {
	type want struct {
		contentType string
		code        int
	}

	type testCase struct {
		name   string
		method string
		url    string
		body   string
		want   want
	}

	tests := []testCase{
		{
			name:   "Post: Bad Request",
			method: http.MethodPost,
			url:    localhost,
			body:   "",
			want: want{
				contentType: "text/plain",
				code:        http.StatusBadRequest,
			},
		},
		{
			name:   "Get: Bad Request",
			method: http.MethodGet,
			url:    fmt.Sprint(localhost, "AbCdE"),
			body:   "",
			want: want{
				contentType: "text/plain",
				code:        http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "plain/text")

			w := httptest.NewRecorder()
			var h func(http.ResponseWriter, *http.Request)

			if tt.method == http.MethodGet {
				h = RootGetHandler
			} else {
				h = RootPostHandler
			}

			h(w, req)

			result := w.Result()
			defer result.Body.Close()

			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)
			assert.Equal(t, tt.want.code, result.StatusCode)
		})
	}
}

func TestFullRootHandler(t *testing.T) {
	type want struct {
		contentType string
		code        int
	}

	type testCase struct {
		name   string
		method string
		url    string
		body   string
		want   want
	}

	var shortURL = ""

	pc := testCase{
		name:   "Post: Created",
		method: http.MethodPost,
		url:    localhost,
		body:   "https://yandex.ru/",
		want: want{
			contentType: "text/plain",
			code:        http.StatusCreated,
		},
	}

	gc := testCase{
		name:   "Get: temp redirect",
		method: http.MethodGet,
		url:    "",
		body:   "",
		want: want{
			contentType: "text/plain",
			code:        http.StatusTemporaryRedirect,
		},
	}

	t.Run("Post + Get yandex url", func(t *testing.T) {
		//post request
		h := RootPostHandler

		req := httptest.NewRequest(pc.method, pc.url, strings.NewReader(pc.body))
		req.Header.Set("Content-Type", "plain/text")

		w := httptest.NewRecorder()
		h(w, req)

		result := w.Result()

		assert.Contains(t, result.Header.Get("Content-Type"), pc.want.contentType)
		assert.Equal(t, pc.want.code, result.StatusCode)

		bytes, err := io.ReadAll(result.Body)
		result.Body.Close()

		require.NoError(t, err)
		shortURL = string(bytes)

		// get request
		h = RootGetHandler

		req = httptest.NewRequest(gc.method, shortURL, nil)
		req.Header.Set("Content-Type", "plain/text")

		w = httptest.NewRecorder()
		h(w, req)

		result = w.Result()
		defer result.Body.Close()

		assert.Contains(t, result.Header.Get("Content-Type"), gc.want.contentType)
		assert.Equal(t, gc.want.code, result.StatusCode)
		assert.Equal(t, pc.body, result.Header.Get("Location"))

	})

}
