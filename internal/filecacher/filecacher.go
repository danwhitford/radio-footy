package filecacher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Getter interface {
	GetUrl(url string) ([]byte, error)
}

type HttpGetter struct {
	client *http.Client
}

func (getter HttpGetter) getWithBackoff(url string) (*http.Response, error) {
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

func (cached HttpGetter) GetUrl(url string) ([]byte, error) {
	resp, err := cached.getWithBackoff(url)
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
	return body, nil
}

type StringGetter struct {
	Contents string
}

func (getter StringGetter) GetUrl(_ string) ([]byte, error) {
	return []byte(getter.Contents), nil
}
