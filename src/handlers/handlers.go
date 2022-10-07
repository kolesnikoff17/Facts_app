package handlers

import (
	"encoding/json"
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

type ServeErr func(w http.ResponseWriter, r *http.Request) *appErr

func (fn ServeErr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Println(e.err)
		http.Error(w, e.msg, e.status)
	}
}

type appErr struct {
	status int
	err    error
	msg    string
}

func RouterTree(w http.ResponseWriter, r *http.Request) *appErr {
	switch r.Method {
	case "GET":
		return getFact(w, r)
	case "PUT":
		return putFact(w, r)
	default:
		return &appErr{status: http.StatusMethodNotAllowed, err: nil, msg: "Allowed: GET, PUT"}
	}
}

func Router(w http.ResponseWriter, r *http.Request) *appErr {
	switch r.Method {
	case "GET":
		rand.Seed(time.Now().UnixNano())
		maxId, err := db.Ins.GetMaxId(r.Context())
		if err != nil {
			log.Printf("getMaxId err: %s", err)
		}
		r.URL.Path += "/" + strconv.Itoa(rand.Intn(maxId-1)+1)
		return getFact(w, r)
	case "POST":
		return postFact(w, r)
	default:
		return &appErr{status: http.StatusMethodNotAllowed, err: nil, msg: "Allowed: GET, POST"}
	}
}

func getFact(w http.ResponseWriter, r *http.Request) *appErr {
	id, err := parseId(r.URL.Path)
	if err != nil {
		return &appErr{status: http.StatusBadRequest, err: err, msg: "Wrong id format"}
	}
	fact, err := db.Ins.GetFactById(r.Context(), id)
	if err == pgx.ErrNoRows {
		return &appErr{status: http.StatusBadRequest, err: err, msg: "No such id"}
	} else if err != nil {
		return &appErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	err = json.NewEncoder(w).Encode(fact)
	if err != nil {
		return &appErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	return nil
}

func postFact(w http.ResponseWriter, r *http.Request) *appErr {
	facts := common.FactsArr{Facts: make([]common.Fact, 1)}
	err := json.NewDecoder(r.Body).Decode(&facts)
	if err != nil {
		return &appErr{status: http.StatusBadRequest, err: err, msg: "Wrong fact format"}
	}
	idList, err := db.Ins.InsertFacts(r.Context(), facts)
	if err != nil {
		return &appErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	err = json.NewEncoder(w).Encode(struct {
		Ids []int `json:"id"`
	}{Ids: idList})
	if err != nil {
		return &appErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	return nil
}

func putFact(w http.ResponseWriter, r *http.Request) *appErr {
	id, err := parseId(r.URL.Path)
	if err != nil {
		return &appErr{status: http.StatusBadRequest, err: err, msg: "Wrong id format"}
	}
	var fact common.Fact
	err = json.NewDecoder(r.Body).Decode(&fact)
	if err != nil {
		return &appErr{status: http.StatusBadRequest, err: err, msg: "Wrong fact format"}
	}
	if fact.Id != id {
		return &appErr{status: http.StatusBadRequest, err: err, msg: "Id mismatch"}
	}
	err = db.Ins.UpdFact(r.Context(), fact, id)
	if err != nil {
		return &appErr{status: http.StatusInternalServerError, err: err, msg: ""}
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func parseId(s string) (int, error) {
	id := 0
	var c rune
	n, err := fmt.Sscanf(s, "/fact/%d%c", &id, &c)
	if n == 1 && err == io.EOF && id > 0 {
		return id, nil
	}
	return 0, err
}
