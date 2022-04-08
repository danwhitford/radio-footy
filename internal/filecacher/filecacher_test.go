package filecacher

import (
	"os"
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

func TestCreatesCache(t *testing.T) {
	setup()
	defer shutdown()

	b, err := GetUrl("https://www.example.com")
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
	
	b, err := GetUrl("https://www.example.com")
	if err != nil {
		t.Fatalf("err was %v", err)
	}
	if b == nil {
		t.Fatal("b was nil")
	}

	b, err = GetUrl("https://www.example.com")
	if err != nil {
		t.Fatalf("err was %v", err)
	}
	if b == nil {
		t.Fatal("b was nil")
	}
}