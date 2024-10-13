package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"time"
)

func getFileHash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func runPeriodically(interval time.Duration, f func()) {
	// Run the function immediately on startup
	f()

	// Run the function on the specified interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		f()
	}
}
