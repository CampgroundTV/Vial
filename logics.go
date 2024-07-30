package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// StreamInfo represents the information of a stream
type StreamInfo struct {
	StreamKey string `json:"stream_key"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	IsLive    bool   `json:"is_live"`
}

func CheckStreamStatus(streamKey string) (StreamInfo, error) {
	url := fmt.Sprintf("http://example.com/stream_info?key=%s", streamKey)
	resp, err := http.Get(url)
	if err != nil {
		return StreamInfo{}, err
	}
	defer resp.Body.Close()

	var info StreamInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return StreamInfo{}, err
	}
	return info, nil
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
	streamKey := "your_stream_key_here"
	db := NewFederationDB()
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	http.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
		streams := db.GetStreams()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(streams)
	})
	go func() {
		fmt.Println("Starting server on :8080")
		http.ListenAndServe(":8080", nil)
	}()

	for {
		select {
		case <-ticker.C:
			info, err := CheckStreamStatus(streamKey)
			if err != nil {
				fmt.Printf("Error checking stream status: %v\n", err)
				continue
			}
			if info.IsLive {
				db.UpdateStream(info)
				fmt.Printf("Stream %s is live with title '%s' in category '%s'\n", info.StreamKey, info.Title, info.Category)
			} else {
				db.RemoveStream(streamKey)
				fmt.Printf("Stream %s is no longer live\n", streamKey)
			}
		}
	}
}
