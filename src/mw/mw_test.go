package mw

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fact/1", nil)
	w := httptest.NewRecorder()
	var mock http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(w, "aboba")
	}
	handler := Logging(mock)
	fn := handler.(http.HandlerFunc)
	fn(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "aboba" {
		t.Errorf("log test: expected aboba, got %v", data)
	}
}

func TestPanicRecovery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fact/1", nil)
	w := httptest.NewRecorder()
	var mock http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		panic("aboba")
	}
	handler := PanicRecovery(mock)
	fn := handler.(http.HandlerFunc)
	fn(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if res.StatusCode != 500 {
		t.Errorf("panic recovery test: expected %v, got %v", http.StatusInternalServerError, data)
	}
	if string(data) != http.StatusText(500)+"\n" {
		t.Errorf("panic recovery test: expected %v, got %v", http.StatusText(500), string(data))
	}
}
