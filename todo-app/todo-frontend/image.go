package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Image struct {
	mu      sync.Mutex
	name    string
	modTime time.Time
}

func GetImage(img *Image) (bool, *Image) {
	img.mu.Lock()
	defer img.mu.Unlock()

	// Cache Hit check
	if !img.modTime.IsZero() && time.Since(img.modTime) < 10*time.Minute {
		return true, img
	}

	// Cache miss -> Download
	filename, err := downloadImage()
	if err != nil {
		log.Printf("Download failed: %v", err)
		return true, img // Fallback to stale data if download fails
	}

	info, err := os.Stat(filename)
	if err != nil {
		log.Printf("Stat failed: %v", err)
		return true, img
	}

	img.name = filename
	img.modTime = info.ModTime()

	return false, img
}

func downloadImage() (string, error) {
	resp, err := http.Get(
		fmt.Sprintf("%s?t=%d", imageURL, time.Now().UnixMilli()),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	f, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}
	return imagePath, nil
}
