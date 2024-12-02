package batch

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/hddskull/urlShorty/internal/storage"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
	"net/http"
)

type (
	batchRequestModel struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
	batchResponseModel struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	requestBatch  []batchRequestModel
	responseBatch []batchResponseModel
)

func BatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reader := r.Body

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			utils.SugaredLogger.Debugln("BatchHandler gzip decompression error:", err)
			formattedError := custom.ErrorResponseModel{Message: err.Error()}
			custom.JSONError(w, formattedError, http.StatusInternalServerError)
		}
		reader = gz
	}
	defer reader.Close()

	//decode batch
	reqBatch := requestBatch{}

	if err := json.NewDecoder(reader).Decode(&reqBatch); err != nil {
		utils.SugaredLogger.Debugln("BatchHandler decoding error:", err)
		formattedError := custom.ErrorResponseModel{Message: err.Error()}
		custom.JSONError(w, formattedError, http.StatusBadRequest)
		return
	}

	//validate and convert to storage model
	storageModels, err := validateAndConvertBatch(reqBatch)
	if err != nil {
		utils.SugaredLogger.Debugln("BatchHandler validation or conversion error:", err)
		formattedError := custom.ErrorResponseModel{Message: err.Error()}
		custom.JSONError(w, formattedError, http.StatusBadRequest)
		return
	}

	//save batch
	savedModels, err := storage.Current.SaveBatch(storageModels)
	if err != nil {
		utils.SugaredLogger.Debugln("BatchHandler saving to storage error:", err)
		formattedError := custom.ErrorResponseModel{Message: err.Error()}
		custom.JSONError(w, formattedError, http.StatusInternalServerError)
		return
	}

	//convert to response model
	savedBatch, err := convertToResponseModel(savedModels)
	if err != nil {
		utils.SugaredLogger.Debugln("BatchHandler validation or conversion error:", err)
		formattedError := custom.ErrorResponseModel{Message: err.Error()}
		custom.JSONError(w, formattedError, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(savedBatch)
	if err != nil {
		utils.SugaredLogger.Debugln("BatchHandler encoding error:", err)
		formattedError := custom.ErrorResponseModel{Message: err.Error()}
		custom.JSONError(w, formattedError, http.StatusInternalServerError)
		return
	}
}

func validateAndConvertBatch(batch requestBatch) ([]model.StorageModel, error) {
	//check if batch contains elements
	if len(batch) == 0 {
		return nil, custom.ErrEmptyBatch
	}

	models := make([]model.StorageModel, len(batch))

	//validate batch item's field and convert to storage model
	for i, item := range batch {
		if item.OriginalURL == "" || item.CorrelationID == "" {
			return nil, custom.ErrInvalidBatch
		}

		m, err := model.NewFileStorageModel(item.OriginalURL, item.CorrelationID)
		if err != nil {
			return nil, err
		}
		models[i] = *m
	}

	return models, nil
}

func convertToResponseModel(arr []model.StorageModel) (responseBatch, error) {
	respBatch := make(responseBatch, len(arr))
	for i, item := range arr {
		brm := batchResponseModel{
			CorrelationID: item.UUID,
			ShortURL:      fmt.Sprint("http://", config.Address.BaseURL, "/", item.ShortURL),
		}
		respBatch[i] = brm
	}
	return respBatch, nil
}
