package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestES(t *testing.T) {
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
		t.Fatalf("Error setting up Elasticsearch: %s", err)
	}

	// Delete the index before each test run
	es.Indices.Delete([]string{"cities"}, es.Indices.Delete.WithIgnoreUnavailable(true))
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)
	res := w.Result()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "OK", string(body))
}

func TestInsertAndRetrieveCity(t *testing.T) {
	setupTestES(t)

	// Start HTTP server handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/city", cityHandler)
	mux.HandleFunc("/population", populationHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Insert city
	city := City{Name: "Berlin", Population: 3769000}
	body, _ := json.Marshal(city)
	resp, err := http.Post(server.URL+"/city", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Retrieve city
	resp, err = http.Get(server.URL + "/population?name=Berlin")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var got City
	json.NewDecoder(resp.Body).Decode(&got)
	assert.Equal(t, "Berlin", got.Name)
	assert.Equal(t, 3769000, got.Population)
}

func TestGetNonExistentCity(t *testing.T) {
	setupTestES(t)

	req := httptest.NewRequest(http.MethodGet, "/population?name=Nowhereville", nil)
	w := httptest.NewRecorder()
	populationHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
