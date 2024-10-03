package shorten

import (
	"encoding/json"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	localhost = "http://localhost:8080/"
	fileName  = "test.json"
)

func setupTest() {
	config.Setup()
	config.StorageFileName = fileName
	storage.SetupStorage()
}

func cleanupTest() error {
	err := os.Remove(fileName)
	return err
}

// Post success
// Post fail - empty model
// Post fail - wrong model
// Post fail - plain/text
func TestPostHandler(t *testing.T) {
	setupTest()

	type want struct {
		contentType string
		code        int
	}

	type testCase struct {
		name         string
		method       string
		url          string
		contentType  string
		requestModel interface{} //requestPostModel
		want         want
	}

	type wrongModel struct {
		Text string `json:"text"`
	}

	tests := []testCase{
		{
			name:        "Post success",
			method:      http.MethodPost,
			url:         localhost,
			contentType: "application/json",
			requestModel: requestPostModel{
				URL: "https://yandex.ru/",
			},
			want: want{
				contentType: "application/json",
				code:        http.StatusCreated,
			},
		},
		{
			name:         "Post fail empty model",
			method:       http.MethodPost,
			url:          localhost,
			contentType:  "application/json",
			requestModel: requestPostModel{},
			want: want{
				contentType: "application/json",
				code:        http.StatusBadRequest,
			},
		},
		{
			name:        "Post fail wrong model",
			method:      http.MethodPost,
			url:         localhost,
			contentType: "application/json",
			requestModel: wrongModel{
				Text: "wrong model, wrong text",
			},
			want: want{
				contentType: "application/json",
				code:        http.StatusBadRequest,
			},
		},
		{
			name:         "Post fail plain text",
			method:       http.MethodPost,
			url:          localhost,
			contentType:  "plain/text",
			requestModel: "https://yandex.ru/",
			want: want{
				contentType: "application/json",
				code:        http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			jsonBytes, _ := json.Marshal(tc.requestModel)
			jsonString := string(jsonBytes)

			req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(jsonString))
			req.Header.Set("Content-Type", tc.contentType)

			w := httptest.NewRecorder()
			PostHandler(w, req)

			result := w.Result()
			defer result.Body.Close()

			assert.Contains(t, result.Header.Get("Content-Type"), tc.want.contentType)
			assert.Equal(t, tc.want.code, result.StatusCode)
		})
	}

	err := cleanupTest()
	assert.NoError(t, err)
}
