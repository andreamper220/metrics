package shared

// Metric определяет структуру JSON-метрики для передачи и/или получения.
type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// Metrics определяет структуру нескольких JSON-метрик для передачи и/или получения.
type Metrics []Metric
