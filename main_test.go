package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Setup Elasticsearch client for testing
	cfg := elasticsearch.Config{
		Addresses: []string{
			getEnv("ES_URL", "http://localhost:9200"),
		},
		Username: getEnv("ES_USERNAME", ""),
		Password: getEnv("ES_PASSWORD", ""),
	}

	var err error
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		panic("Failed to create ES client in tests: " + err.Error())
	}

	code := m.Run()
	os.Exit(code)
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "OK", string(body))
}

func TestCityHandlerAndPopulationHandler(t *testing.T) {
	// POST to /city
	city := City{Name: "Testopolis", Population: 123456}
	body, _ := json.Marshal(city)

	postReq := httptest.NewRequest(http.MethodPost, "/city", bytes.NewReader(body))
	postW := httptest.NewRecorder()

	cityHandler(postW, postReq)

	postResp := postW.Result()
	assert.Equal(t, http.StatusOK, postResp.StatusCode)

	var postResponse Response
	json.NewDecoder(postResp.Body).Decode(&postResponse)
	assert.Equal(t, "City stored/updated successfully", postResponse.Message)

	// GET from /population
	getReq := httptest.NewRequest(http.MethodGet, "/population?name=Testopolis", nil)
	getW := httptest.NewRecorder()

	populationHandler(getW, getReq)

	getResp := getW.Result()
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var returned City
	json.NewDecoder(getResp.Body).Decode(&returned)
	assert.Equal(t, "Testopolis", returned.Name)
	assert.Equal(t, 123456, returned.Population)
}

func TestCityHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/city", bytes.NewReader([]byte("{bad json")))
	w := httptest.NewRecorder()

	cityHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var errResp Response
	json.NewDecoder(w.Body).Decode(&errResp)
	assert.Equal(t, "Invalid JSON", errResp.Message)
}

func TestPopulationHandler_NoName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/population", nil)
	w := httptest.NewRecorder()

	populationHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	var errResp Response
	json.NewDecoder(w.Body).Decode(&errResp)
	assert.Equal(t, "City name required", errResp.Message)
}

func TestPopulationHandler_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/population?name=nonexistentcity", nil)
	w := httptest.NewRecorder()

	populationHandler(w, req)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	var errResp Response
	json.NewDecoder(w.Body).Decode(&errResp)
	assert.Equal(t, "City not found", errResp.Message)
}
