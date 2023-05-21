package filecacher

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Getter interface {
	Get(url string) (*http.Response, error)
}

type HttpGetter struct{}
func (getter HttpGetter) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func getAndSave(url, fname string, getter Getter) ([]byte, error) {
	resp, err := getter.Get(url)
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

func GetUrl(url string, getter Getter) ([]byte, error) {
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
			return getAndSave(url, fname, getter)
		}
	} else {
		log.Println("File exists, checking mod time")
		info, err := os.Stat(fmt.Sprintf(".cache/%s", fname))
		if err != nil {
			return nil, err
		}
		if info.ModTime().Before(time.Now().Add(-1 * time.Hour * 24)) {
			log.Println("File is older than 24 hours, fetching and caching...")
			return getAndSave(url, fname, getter)
		}
		log.Println("File is recent, returning cached data")
		return data, nil
	}
}
