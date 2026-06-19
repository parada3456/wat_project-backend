package scraper

import (
	"crypto/md5"
	"fmt"
	"log"
)

// GenerateJobID creates a deterministic job ID from a source URL using MD5 hashing.
func GenerateJobID(url string) string {
	log.Println("debugprint: entering GenerateJobID")
	hash := md5.Sum([]byte(url))
	return fmt.Sprintf("%x", hash)
}
