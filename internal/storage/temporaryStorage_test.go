package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	save = "save"
	get  = "get"
)

func TestTemporaryStorage(t *testing.T) {
	storage := newTemporaryStorage()

	//test for errors
	tests := []struct {
		name  string
		fType string
		url   string
		id    string
	}{
		{
			name:  "save() empty url",
			fType: save,
			url:   "",
			id:    "",
		}, {
			name:  "get() empty id",
			fType: save,
			url:   "",
			id:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var str string
			var err error

			if tc.fType == save {
				str, err = storage.Save(tc.url)
			} else {
				str, err = storage.Get(tc.id)
			}

			assert.Error(t, err)
			assert.Equal(t, str, "")
		})
	}

	//test for success
	var url = "https://practicum.yandex.ru/"

	t.Run("save() + get() success", func(t *testing.T) {

		id, err := storage.Save(url)
		require.NoError(t, err)

		savedURL, err := storage.Get(id)
		require.NoError(t, err)
		assert.Equal(t, url, savedURL)
	})
}
