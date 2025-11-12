package main

// import (
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// )

// func TestWriteJSON(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	payload := map[string]string{"message": "hello"}
// 	writeJSON(w, http.StatusOK, payload)

// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
// 	}

// 	contentType := res.Header.Get("Content-Type")
// 	if contentType != "application/json" {
// 		t.Errorf("expected Content-Type application/json, got %s", contentType)
// 	}

// 	var body map[string]string
// 	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
// 		t.Fatalf("could not decode response body: %v", err)
// 	}
// 	if body["message"] != "hello" {
// 		t.Errorf("expected message 'hello', got '%s'", body["message"])
// 	}

// 	w = httptest.NewRecorder()
// 	writeJSON(w, http.StatusNoContent, nil)

// 	res = w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusNoContent {
// 		t.Errorf("expected status %d, got %d", http.StatusNoContent, res.StatusCode)
// 	}
// 	if w.Body.Len() != 0 {
// 		t.Errorf("expected empty body, got %s", w.Body.String())
// 	}
// }

// func TestReadJSON(t *testing.T) {
// 	type testStruct struct {
// 		Name string `json:"name"`
// 		Age  int    `json:"age"`
// 	}
// 	jsonStr := `{"name":"John","age":30}`
// 	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonStr))

// 	var dst testStruct
// 	err := readJSON(req, &dst)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}
// 	if dst.Name != "John" {
// 		t.Errorf("expected name 'John', got '%s'", dst.Name)
// 	}
// 	if dst.Age != 30 {
// 		t.Errorf("expected age 30, got %d", dst.Age)
// 	}

// 	jsonStr = `{"name":"John","age":30,}`
// 	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonStr))
// 	err = readJSON(req, &dst)
// 	if err == nil {
// 		t.Fatal("expected an error for invalid JSON, got nil")
// 	}

// 	jsonStr = `{"name":"Jane","age":25,"city":"New York"}`
// 	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonStr))
// 	err = readJSON(req, &dst)
// 	if err == nil {
// 		t.Fatal("expected an error for unknown fields, got nil")
// 	}
// }
