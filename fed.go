package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type StreamInfo struct {
	StreamKey string `json:"stream_key"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	IsLive    bool   `json:"is_live"`
}

type FederationDB struct {
	mu          sync.Mutex
	liveStreams map[string]StreamInfo
}


func NewFederationDB() *FederationDB {
	return &FederationDB{
		liveStreams: make(map[string]StreamInfo),
	}
}


func (db *FederationDB) UpdateStream(info StreamInfo) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.liveStreams[info.StreamKey] = info
}

func (db *FederationDB) RemoveStream(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.liveStreams, key)
}


func (db *FederationDB) GetStreams() []StreamInfo {
	db.mu.Lock()
	defer db.mu.Unlock()
	streams := make([]StreamInfo, 0, len(db.liveStreams))
	for _, info := range db.liveStreams {
		streams = append(streams, info)
	}
	return streams
}

func main() {
	db := NewFederationDB()

	http.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
		streams := db.GetStreams()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(streams)
	})
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
