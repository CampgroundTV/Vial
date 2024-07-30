package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// StreamInfo represents the information of a stream
type StreamInfo struct {
	StreamKey string `json:"stream_key"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	IsLive    bool   `json:"is_live"`
}

// FederationDB is a simple in-memory database to store live stream info
type FederationDB struct {
	mu          sync.Mutex
	liveStreams map[string]StreamInfo
}

// NewFederationDB initializes a new FederationDB
func NewFederationDB() *FederationDB {
	return &FederationDB{
		liveStreams: make(map[string]StreamInfo),
	}
}

// UpdateStream updates the information of a stream in the database
func (db *FederationDB) UpdateStream(info StreamInfo) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.liveStreams[info.StreamKey] = info
}

// RemoveStream removes a stream key from the database
func (db *FederationDB) RemoveStream(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.liveStreams, key)
}

// GetStreams returns all live streams
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

	// Start HTTP server
	http.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
		streams := db.GetStreams()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(streams)
	})
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
