package filecacher

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"log"
	"fmt"
)

func getAndSave(url, fname string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("Writing file...")
	err = os.WriteFile(filepath.Join(".cache", fname), body, 0644)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetUrl(url string) ([]byte, error) {
	log.Printf("Getting URL: %s\n", url)	
	h := sha1.Sum([]byte(url))
	fname := hex.EncodeToString(h[:])
	data, err := os.ReadFile(fmt.Sprintf(".cache/%s", fname))
	if err != nil {
		log.Println("Could not open file from cache")
		if !os.IsNotExist(err) {
			log.Println("File open error was not recognised. This is bad.")
			return nil, err
		} else {
			log.Println("File does not exist, fetching and caching...")
			return getAndSave(url, fname)
		}
	} else {
		log.Println("File exists, using cache")
		return data, nil
	}
}
