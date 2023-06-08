package main

type MemStorage struct {
	metrics map[string]interface{}
	//mutex   sync.Mutex
}

func NewMetricsStore() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]interface{}),
	}
}
