package config

import (
	"reflect"
	"testing"
)

func TestParseCORSOrigins(t *testing.T) {
	got := parseCORSOrigins("https://budget.example.com, https://other.example.com")
	want := []string{"https://budget.example.com", "https://other.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestParseCORSOriginsDefault(t *testing.T) {
	got := parseCORSOrigins("")
	if !reflect.DeepEqual(got, defaultCORSOrigins) {
		t.Fatalf("expected defaults, got %v", got)
	}
}
