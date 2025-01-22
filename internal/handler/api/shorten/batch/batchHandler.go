package batch

import (
	"compress/gzip"
	"context"
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
			custom.JSONError(w, err, http.StatusInternalServerError)
		}
		reader = gz
	}
	defer reader.Close()

	//decode batch
	reqBatch := requestBatch{}

	if err := json.NewDecoder(reader).Decode(&reqBatch); err != nil {
		utils.SugaredLogger.Debugln("BatchHandler decoding error:", err)
		custom.JSONError(w, err, http.StatusBadRequest)
		return
	}

	//validate and convert to storage model
	storageModels, err := validateAndConvertBatch(r.Context(), reqBatch)
	if err != nil {
		utils.SugaredLogger.Debugln("BatchHandler validation or conversion error:", err)
		custom.JSONError(w, err, http.StatusBadRequest)
		return
	}

	//save batch
	err = storage.Current.SaveBatch(r.Context(), storageModels)
	if err != nil {
		utils.SugaredLogger.Debugln("BatchHandler saving to storage error:", err)
		custom.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	//convert to response model
	savedBatch, err := convertToResponseModel(storageModels)
	if err != nil {
		utils.SugaredLogger.Debugln("BatchHandler validation or conversion error:", err)
		custom.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(savedBatch)
	if err != nil {
		utils.SugaredLogger.Debugln("BatchHandler encoding error:", err)
		custom.JSONError(w, err, http.StatusInternalServerError)
		return
	}
}

func validateAndConvertBatch(ctx context.Context, batch requestBatch) ([]model.StorageModel, error) {
	//check if batch contains elements
	if len(batch) == 0 {
		return nil, custom.ErrEmptyBatch
	}

	sessionID, ok := model.SessionIDFromContext(ctx)

	if !ok {
		return nil, custom.ErrNoSessionID
	}

	models := make([]model.StorageModel, len(batch))

	//validate batch item's field and convert to storage model
	for i, item := range batch {
		if item.OriginalURL == "" || item.CorrelationID == "" {
			return nil, custom.ErrInvalidBatch
		}

		m, err := model.NewFileStorageModel(item.OriginalURL, item.CorrelationID, sessionID)
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
