package filecacher

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func setup() {
	os.Mkdir(".cache", 0755)
}

func shutdown() {
	err := os.RemoveAll(".cache")
	if err != nil {
		panic(err)
	}
}

type DummyGetter struct {
	used bool
}

func (getter *DummyGetter) Get(url string) (*http.Response, error) {
	if getter.used {
		return nil, fmt.Errorf("getter has been used")
	}
	getter.used = true
	response := http.Response{
		Body:       io.NopCloser(strings.NewReader("Hello, world!")),
		StatusCode: 200,
	}
	return &response, nil
}

func TestCreatesCache(t *testing.T) {
	setup()
	defer shutdown()

	getter := DummyGetter{}
	b, err := GetUrl("https://www.example.com", &getter)
	if err != nil {
		t.Fatalf("err was %v", err)
	}
	if b == nil {
		t.Fatal("b was nil")
	}

	b, err = os.ReadFile(".cache/740e7397907c0b004010d92b33d283e98f74063d")
	if err != nil {
		t.Fatalf("err was %v", err)
	}
	if b == nil {
		t.Fatal("b was nil")
	}
}

func TestUsesCache(t *testing.T) {
	setup()
	defer shutdown()

	getter := DummyGetter{}
	b, err := GetUrl("https://www.example.com", &getter)
	if err != nil {
		t.Fatalf("err was %v", err)
	}
	if b == nil {
		t.Fatal("b was nil")
	}

	b, err = GetUrl("https://www.example.com", &getter)
	if err != nil {
		t.Fatalf("err was %v", err)
	}
	if b == nil {
		t.Fatal("b was nil")
	}
}
