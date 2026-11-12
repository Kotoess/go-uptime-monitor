package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckURLSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result, err := checkURL(server.Client(), server.URL)

	if err != nil {
		t.Fatalf("не ожидалась ошибка запроса: %v", err)
	}

	if !result.Success {
		t.Errorf("ожидался успешный результат для HTTP %d", http.StatusOK)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf(
			"ожидался HTTP-код %d, получен %d",
			http.StatusOK,
			result.StatusCode,
		)
	}

	if result.URL != server.URL {
		t.Errorf("ожидался URL %q, получен %q", server.URL, result.URL)
	}
}

func TestCheckURLServerErrorAsDown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	result, err := checkURL(server.Client(), server.URL)

	if err != nil {
		t.Fatalf("не ожидалась ошибка запроса: %v", err)
	}

	if result.Success {
		t.Errorf("HTTP %d должен считаться недоступностью", http.StatusServiceUnavailable)
	}

	if result.StatusCode != http.StatusServiceUnavailable {
		t.Errorf(
			"ожидался HTTP-код %d, получен %d",
			http.StatusServiceUnavailable,
			result.StatusCode,
		)
	}
}
