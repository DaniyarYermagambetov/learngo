package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Типы метрик
const (
	Gauge   = "gauge"
	Counter = "counter"
)

// GaugeStorage хранит метрики типа gauge
type GaugeStorage struct {
	sync.RWMutex
	data map[string]float64
}

// CounterStorage хранит метрики типа counter
type CounterStorage struct {
	sync.RWMutex
	data map[string]int64
}

// NewGaugeStorage создает новый экземпляр GaugeStorage
func NewGaugeStorage() *GaugeStorage {
	return &GaugeStorage{
		data: make(map[string]float64),
	}
}

// NewCounterStorage создает новый экземпляр CounterStorage
func NewCounterStorage() *CounterStorage {
	return &CounterStorage{
		data: make(map[string]int64),
	}
}

// MemStorage хранит все метрики
type MemStorage struct {
	gaugeStorage   *GaugeStorage
	counterStorage *CounterStorage
}

// NewMemStorage создает новый экземпляр MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeStorage:   NewGaugeStorage(),
		counterStorage: NewCounterStorage(),
	}
}

// SetGauge устанавливает значение метрики типа gauge
func (s *GaugeStorage) SetGauge(name string, value float64) {
	s.Lock()
	defer s.Unlock()
	s.data[name] = value
}

// AddCounter добавляет значение к метрике типа counter
func (s *CounterStorage) AddCounter(name string, value int64) {
	s.Lock()
	defer s.Unlock()
	s.data[name] += value
}

// GetGauge возвращает значение метрики типа gauge
func (s *GaugeStorage) GetGauge(name string) (float64, bool) {
	s.RLock()
	defer s.RUnlock()
	value, exists := s.data[name]
	return value, exists
}

// GetCounter возвращает значение метрики типа counter
func (s *CounterStorage) GetCounter(name string) (int64, bool) {
	s.RLock()
	defer s.RUnlock()
	value, exists := s.data[name]
	return value, exists
}

func updateMetricHandler(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Разбираем URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Invalid URL format", http.StatusNotFound)
			return
		}

		metricType := parts[2]
		metricName := parts[3]
		metricValueStr := parts[4]

		if metricName == "" {
			http.Error(w, "Metric name is required", http.StatusNotFound)
			return
		}

		// Преобразуем значение метрики
		metricValue, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}

		switch metricType {
		case Gauge:
			storage.gaugeStorage.SetGauge(metricName, metricValue)
		case Counter:
			intValue := int64(metricValue) // Преобразуем значение в int64 для хранения в counter
			storage.counterStorage.AddCounter(metricName, intValue)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	storage := NewMemStorage()
	http.HandleFunc("/update/", updateMetricHandler(storage))
	serverAddr := "localhost:8080"
	println("Starting server on", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		panic(err)
	}
}
