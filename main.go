package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

var es *elasticsearch.Client

type City struct {
	Name       string `json:"name"`
	Population int    `json:"population"`
}

func main() {
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
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/city", cityHandler)
	http.HandleFunc("/population", populationHandler)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func cityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var city City
	if err := json.NewDecoder(r.Body).Decode(&city); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	docID := strings.ToLower(city.Name)
	data, _ := json.Marshal(city)

	res, err := es.Index(
		"cities",
		bytes.NewReader(data),
		es.Index.WithDocumentID(docID),
		es.Index.WithRefresh("true"),
	)
	if err != nil {
		http.Error(w, "Error indexing document", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		http.Error(w, fmt.Sprintf("Indexing error: %s", res.String()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("City stored/updated successfully"))
}

func populationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	cityName := r.URL.Query().Get("name")
	if cityName == "" {
		http.Error(w, "City name required", http.StatusBadRequest)
		return
	}

	docID := strings.ToLower(cityName)

	res, err := es.Get("cities", docID)
	if err != nil {
		http.Error(w, "Error retrieving document", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		http.Error(w, "City not found", http.StatusNotFound)
		return
	}

	var doc struct {
		Source City `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&doc); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(doc.Source)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
