package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"httpServer/src/common"
	"httpServer/src/db"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func isEqual(a, b common.Fact) bool {
	if a.ID != b.ID || a.Title != b.Title || a.Desc != b.Desc || len(a.Links) != len(b.Links) {
		return false
	}
	for i := 0; i < len(a.Links) && i < len(b.Links); i++ {
		if a.Links[i] != b.Links[i] {
			return false
		}
	}
	return true
}

func TestSetup(t *testing.T) {
	ctx := context.Background()
	db.Ins.Db = db.InitDb(ctx)
}

func TestRouterTreeGetId(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fact/1", nil)
	w := httptest.NewRecorder()
	ServeErr(RouterTree).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	fact := common.Fact{}
	err := json.NewDecoder(res.Body).Decode(&fact)
	if err != nil {
		t.Fatalf("router (tree) get method valid test: expected err nil, got %v", err)
	}
	factEx := common.Fact{ID: 1, Title: "aboba", Desc: "aboba_aboba", Links: []string{"aboba", "aboba_aboba"}}
	if !isEqual(fact, factEx) {
		t.Errorf("router (tree) get method valid test: expected %v, got %v", factEx, fact)
	}
}

func TestRouterTreeGetIdInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fact/666", nil)
	w := httptest.NewRecorder()
	ServeErr(RouterTree).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("router (tree) get method too large id test: expected err nil, got %v", err)
	}
	if expected := "No such id\n"; string(data) != expected {
		t.Errorf("router (tree) get method too large id test: expected %v, got %v", expected, string(data))
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("router (tree) get method too large id test: expected %v, got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestRouterTreeGetIdWrong(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fact/666asdas", nil)
	w := httptest.NewRecorder()
	ServeErr(RouterTree).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("router (tree) get method wrong id format test: expected err nil, got %v", err)
	}
	if expected := "Wrong id format\n"; string(data) != expected {
		t.Errorf("router (tree) get method wrong id format test: expected %v, got %v", expected, string(data))
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("router (tree) get method wrong id format test: expected %v, got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestRouterGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fact", nil)
	w := httptest.NewRecorder()
	ServeErr(Router).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	fact := common.Fact{}
	err := json.NewDecoder(res.Body).Decode(&fact)
	if err != nil {
		t.Errorf("router get method test: expected err nil, got %v", err)
		t.Errorf("expected statusCode 200, got %v", res.StatusCode)
	}
}

func TestRouterPost(t *testing.T) {
	var b bytes.Buffer
	facts := common.FactsArr{Facts: []common.Fact{{Title: "a", Desc: "a"},
		{Title: "b", Desc: "b", Links: []string{"ad"}}}}
	err := json.NewEncoder(&b).Encode(facts)
	if err != nil {
		t.Fatalf("router post method test: expected err nil, got %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/fact", &b)
	w := httptest.NewRecorder()
	ServeErr(Router).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("router post method test: expected %d, got %d", http.StatusOK, res.StatusCode)
	}
	expected := struct {
		Ids []int `json:"ids"`
	}{[]int{6, 7}}
	got := struct {
		Ids []int `json:"ids"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&got)
	if err != nil {
		t.Fatalf("router post method test: expected err nil, got %v", err)
	}
	fmt.Println(expected, got)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("router post method test: expected %v, got %v", expected, got)
	}
}

func TestRouterPostInvalid(t *testing.T) {
	var b bytes.Buffer
	facts := []struct {
		Name string `json:"name"`
	}{{Name: "asdda"}, {Name: "dasd"}}
	err := json.NewEncoder(&b).Encode(facts)
	if err != nil {
		t.Fatalf("router post method test: expected err nil, got %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/fact", &b)
	w := httptest.NewRecorder()
	ServeErr(Router).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("router post method test: expected err nil, got %v", err)
	}
	if expected := "Wrong fact format\n"; string(data) != expected {
		t.Errorf("router post method test: expected %v, got %v", expected, string(data))
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("router post method test: expected %v, got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestRouterTreePut(t *testing.T) {
	var b bytes.Buffer
	facts := common.Fact{ID: 1, Title: "a", Desc: "a"}
	err := json.NewEncoder(&b).Encode(facts)
	if err != nil {
		t.Fatalf("router (tree) put method test: expected err nil, got %v", err)
	}
	req := httptest.NewRequest(http.MethodPut, "/fact/1", &b)
	w := httptest.NewRecorder()
	ServeErr(RouterTree).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("router (tree) put method test: expected %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestRouterTreePutMismatch(t *testing.T) {
	var b bytes.Buffer
	facts := common.Fact{ID: 2, Title: "a", Desc: "a"}
	err := json.NewEncoder(&b).Encode(facts)
	if err != nil {
		t.Fatalf("router (tree) put method test: expected err nil, got %v", err)
	}
	req := httptest.NewRequest(http.MethodPut, "/fact/1", &b)
	w := httptest.NewRecorder()
	ServeErr(RouterTree).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("router (tree) put method test: expected err nil, got %v", err)
	}
	if expected := "ID mismatch\n"; expected != string(data) {
		t.Errorf("router (tree) put method test: expected %v, got %v", expected, string(data))
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("router (tree) put method test: expected %d, got %d", http.StatusBadRequest, res.StatusCode)
	}
}
