package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"httpServer/src/common"
	"httpServer/src/db"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// ServeErr is a wrapper type for error handling
type ServeErr func(w http.ResponseWriter, r *http.Request) *AppErr

// ServeHTTP made to implement http.Handler
func (fn ServeErr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Println(e.err)
		http.Error(w, e.msg, e.status)
	}
}

// AppErr is a custom error
type AppErr struct {
	status int
	err    error
	msg    string
}

// RouterTree is a method router for /fact/ request
func RouterTree(w http.ResponseWriter, r *http.Request) *AppErr {
	switch r.Method {
	case "GET":
		return getFact(w, r)
	case "PUT":
		return putFact(w, r)
	default:
		return &AppErr{status: http.StatusMethodNotAllowed, err: nil, msg: "Allowed: GET, PUT"}
	}
}

// Router is a method router for /fact request
func Router(w http.ResponseWriter, r *http.Request) *AppErr {
	switch r.Method {
	case "GET":
		rand.Seed(time.Now().Unix())
		maxID, err := db.Ins.GetMaxID(r.Context())
		if err != nil {
			log.Printf("getMaxId err: %s", err)
		}
		r.URL.Path += "/" + strconv.Itoa(rand.Intn(maxID)+1)
		return getFact(w, r)
	case "POST":
		return postFact(w, r)
	default:
		return &AppErr{status: http.StatusMethodNotAllowed, err: nil, msg: "Allowed: GET, POST"}
	}
}

func getFact(w http.ResponseWriter, r *http.Request) *AppErr {
	id, err := parseID(r.URL.Path)
	fmt.Println(id)
	if err != nil {
		return &AppErr{status: http.StatusBadRequest, err: err, msg: "Wrong id format"}
	}
	fact, err := db.Ins.GetFactByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		return &AppErr{status: http.StatusBadRequest, err: err, msg: "No such id"}
	} else if err != nil {
		return &AppErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	err = json.NewEncoder(w).Encode(fact)
	if err != nil {
		return &AppErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	return nil
}

func postFact(w http.ResponseWriter, r *http.Request) *AppErr {
	facts := common.FactsArr{Facts: make([]common.Fact, 1)}
	err := json.NewDecoder(r.Body).Decode(&facts)
	if err != nil {
		return &AppErr{status: http.StatusBadRequest, err: err, msg: "Wrong fact format"}
	}
	idList, err := db.Ins.InsertFacts(r.Context(), facts)
	if err != nil {
		return &AppErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	err = json.NewEncoder(w).Encode(struct {
		Ids []int `json:"id"`
	}{Ids: idList})
	if err != nil {
		return &AppErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	return nil
}

func putFact(w http.ResponseWriter, r *http.Request) *AppErr {
	id, err := parseID(r.URL.Path)
	if err != nil {
		return &AppErr{status: http.StatusBadRequest, err: err, msg: "Wrong id format"}
	}
	var fact common.Fact
	err = json.NewDecoder(r.Body).Decode(&fact)
	if err != nil {
		return &AppErr{status: http.StatusBadRequest, err: err, msg: "Wrong fact format"}
	}
	if fact.ID != id {
		return &AppErr{status: http.StatusBadRequest, err: err, msg: "ID mismatch"}
	}
	err = db.Ins.UpdFact(r.Context(), fact, id)
	if err != nil {
		return &AppErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func parseID(s string) (int, error) {
	id := 0
	var c rune
	n, err := fmt.Sscanf(s, "/fact/%d%c", &id, &c)
	if n > 1 && err != io.EOF {
		return 0, fmt.Errorf("wrong URL format: %w", err)
	}
	if id < 1 {
		return 0, errors.New("id less then one")
	}

	return id, nil
}
