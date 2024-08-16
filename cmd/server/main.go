package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// MemStorage представляет хранилище метрик
type MemStorage struct {
	mu       sync.Mutex
	gauges   map[string]float64
	counters map[string]int64
}

// NewMemStorage создает новый экземпляр MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func main() {
	memStorage := NewMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		handlePost(w, r, memStorage)
	})

	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request, storage *MemStorage) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	url := r.URL.Path
	split := strings.Split(strings.Trim(url, "/"), "/")

	if len(split) != 3 {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return
	}

	metricType := split[0]
	metricName := split[1]
	metricValue := split[2]

	if metricName == "" {
		http.Error(w, "Metric name is required", http.StatusNotFound)
		return
	}

	var err error
	switch metricType {
	case "gauge":
		var value float64
		value, err = strconv.ParseFloat(metricValue, 64)
		if err == nil {
			storage.mu.Lock()
			storage.gauges[metricName] = value
			storage.mu.Unlock()
			w.WriteHeader(http.StatusOK)
		}
	case "counter":
		var value int64
		value, err = strconv.ParseInt(metricValue, 10, 64)
		if err == nil {
			storage.mu.Lock()
			storage.counters[metricName] += value
			storage.mu.Unlock()
			w.WriteHeader(http.StatusOK)
		}
	default:
		err = fmt.Errorf("invalid metric type")
	}

	if err != nil {
		http.Error(w, "Invalid metric value", http.StatusBadRequest)
	}
}
