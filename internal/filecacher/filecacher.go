package filecacher

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Getter interface {
	Get(url string) (*http.Response, error)
}

type HttpGetter struct {
	client *http.Client
}

func (getter HttpGetter) Get(url string) (*http.Response, error) {
	backoff := 1
	limit := 64000
	for {
		res, err := getter.client.Get(url)
		if err != nil {
			return nil, err
		}
		if res.StatusCode == 200 {
			return res, nil
		}
		// Sleep for backoff seconds
		fmt.Printf("Got status code %d, sleeping for %d seconds\n", res.StatusCode, backoff)
		time.Sleep(time.Second * time.Duration(backoff))
		backoff *= 2
		if backoff > limit {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, fmt.Errorf("error reading body from error state: %w", err)
			}
			return nil,
				fmt.Errorf("backoff limit reached. Error getting URL: %s. Code: %d. Body: %s", url, res.StatusCode, string(body))
		}
	}
}

func NewHttpGetter() HttpGetter {
	// Create http client with custom timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	return HttpGetter{client: client}
}

func getAndSave(url, fname string, getter Getter) ([]byte, error) {
	resp, err := getter.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading body from error state: %w", err)
		}
		return nil, fmt.Errorf("error getting url: %s. code: %d. body: %s", url, resp.StatusCode, string(body))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("Writing file...")
	err = os.WriteFile(filepath.Join(".cache", fname), body, 0644)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetUrl(url string, getter Getter) ([]byte, error) {
	h := sha1.Sum([]byte(url))
	fname := hex.EncodeToString(h[:])
	data, err := os.ReadFile(fmt.Sprintf(".cache/%s", fname))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("file open error was not recognised. This is bad. url: %s. error: %w", url, err)
		} else {
			return getAndSave(url, fname, getter)
		}
	} else {
		info, err := os.Stat(fmt.Sprintf(".cache/%s", fname))
		if err != nil {
			return nil, err
		}
		if info.ModTime().Before(time.Now().Add(-1 * time.Hour * 24)) {
			return getAndSave(url, fname, getter)
		}
		return data, nil
	}
}
