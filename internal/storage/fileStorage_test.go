package storage

import (
	"encoding/json"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var fileName = "test.json"
var testStorage = FileStorage{}

func setupTest() error {
	config.StorageFileName = fileName
	//create file
	testFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	//create file data
	m := []model.StorageModel{
		{
			UUID:        "test",
			ShortURL:    "t.ru",
			OriginalURL: "test.ru",
		},
	}
	//marshal file data
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	//write data to file
	_, err = testFile.Write(bytes)
	if err != nil {
		return err
	}
	//close file
	err = testFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func cleanupTest() error {
	err := os.Remove(fileName)
	return err
}

func TestReadAllFromFile(t *testing.T) {
	err := setupTest()
	require.NoError(t, err)

	t.Run("read from file", func(t *testing.T) {

		slice, err := testStorage.readAllFromFile()
		require.NoError(t, err)
		assert.NotEmpty(t, slice)
	})

	err = cleanupTest()
	require.NoError(t, err)
}

func TestCheckExistence(t *testing.T) {
	err := setupTest()
	require.NoError(t, err)

	tests := []struct {
		name         string
		originalURL  string
		modelIsEmpty bool
	}{
		{
			name:         "model exists",
			originalURL:  "test.ru",
			modelIsEmpty: false,
		},
		{
			name:         "model doesn't exist",
			originalURL:  "nothing.com",
			modelIsEmpty: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			model, err := testStorage.checkExistence(tc.originalURL)
			require.NoError(t, err)
			if tc.modelIsEmpty {
				assert.Empty(t, model)
			} else {
				assert.NotEmpty(t, model)
			}
		})
	}
	err = cleanupTest()
	require.NoError(t, err)
}

func TestSaveToFile(t *testing.T) {
	err := setupTest()
	require.NoError(t, err)

	t.Run("save to file", func(t *testing.T) {
		model, err := model.NewFileStorageModel("someurl.com")
		require.NoError(t, err)

		err = testStorage.saveToFile(model)
		require.NoError(t, err)
	})

	err = cleanupTest()
	require.NoError(t, err)
}

func TestFileStorage_Save(t *testing.T) {
	err := setupTest()
	require.NoError(t, err)

	tests := []struct {
		name        string
		originalURL string
		returnError bool
	}{
		{
			name:        "normal URL",
			originalURL: "https://practicum.yandex.ru/",
			returnError: false,
		},
		{
			name:        "empty URL",
			originalURL: "",
			returnError: true,
		},
		{
			name:        "another URL",
			originalURL: "https://yandex.ru/",
			returnError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			shortened, err := testStorage.Save(tc.originalURL)
			if tc.returnError {
				assert.Equal(t, "", shortened)
				require.Error(t, err)
			} else {
				assert.NotEmpty(t, shortened)
				require.NoError(t, err)
			}
		})
	}

	err = cleanupTest()
	require.NoError(t, err)
}

func TestFileStorage_Get(t *testing.T) {
	err := setupTest()
	require.NoError(t, err)

	shortPracticum, err := testStorage.Save("https://practicum.yandex.ru/")
	if err != nil {
		t.Fatal(err)
	}
	_, err = testStorage.Save("https://yandex.ru/")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		originalURL string
		shortURL    string
		returnError bool
	}{
		{
			name:        "normal URL",
			originalURL: "https://practicum.yandex.ru/",
			shortURL:    shortPracticum,
			returnError: false,
		},
		{
			name:        "empty URL",
			originalURL: "",
			shortURL:    "",
			returnError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			original, err := testStorage.Get(tc.shortURL)

			if tc.returnError {
				assert.Equal(t, "", original)
				require.Error(t, err)
			} else {
				assert.Equal(t, tc.originalURL, original)
				require.NoError(t, err)
			}
		})
	}

	err = cleanupTest()
	require.NoError(t, err)
}
