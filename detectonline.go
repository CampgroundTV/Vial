package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func main() {
	streamKey := "your_stream_key_here"

	info, err := CheckStreamStatus(streamKey)
	if err != nil {
		fmt.Printf("Error checking stream status: %v\n", err)
		return
	}

	if info.IsLive {
		fmt.Printf("Stream %s is live with title '%s' in category '%s'\n", info.StreamKey, info.Title, info.Category)
	} else {
		fmt.Printf("Stream %s is not live\n", info.StreamKey)
	}
}
