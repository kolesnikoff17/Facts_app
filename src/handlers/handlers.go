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

type errorHandler struct {
	http.ResponseWriter
}

func (h errorHandler) errorHandle(status int, err error, body string) {
	log.Println(err)
	h.WriteHeader(status)
	if status != 500 {
		_, err = io.WriteString(h, body)
		if err != nil {
			log.Println(err)
		}
	}
}

func Router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Path == "/fact/" {
			rand.Seed(time.Now().UnixNano())
			maxId, err := db.Ins.GetMaxId(r.Context())
			if err != nil {
				log.Printf("getMaxId err: %s", err)
			}
			r.URL.Path += strconv.Itoa(rand.Intn(maxId))
		}
		getFact(w, r)
	case "POST":
		postFact(w, r)
	case "PUT":
		putFact(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getFact(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r.URL.Path)
	if err != nil {
		errorHandler{w}.errorHandle(http.StatusBadRequest, err, "Wrong id format")
		return
	}
	fact, err := db.Ins.GetFactById(r.Context(), id)
	if err == pgx.ErrNoRows {
		errorHandler{w}.errorHandle(http.StatusBadRequest, err, "No such id")
		return
	} else if err != nil {
		errorHandler{w}.errorHandle(http.StatusInternalServerError, err, "")
		return
	}
	err = json.NewEncoder(w).Encode(fact)
	if err != nil {
		errorHandler{w}.errorHandle(http.StatusInternalServerError, err, "")
		return
	}
}

func postFact(w http.ResponseWriter, r *http.Request) {
	facts := make([]common.Fact, 1)
	err := json.NewDecoder(r.Body).Decode(&facts)
	if err != nil {
		errorHandler{w}.errorHandle(http.StatusBadRequest, err, "Wrong fact format")
		return
	}
	idList, err := db.Ins.InsertFacts(r.Context(), facts)
	if err != nil {
		errorHandler{w}.errorHandle(http.StatusInternalServerError, err, "")
		return
	}
	err = json.NewEncoder(w).Encode(struct {
		Ids []int `json:"id"`
	}{Ids: idList})
	if err != nil {
		errorHandler{w}.errorHandle(http.StatusInternalServerError, err, "")
		return
	}
}

func putFact(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r.URL.Path)
	if err != nil {
		errorHandler{w}.errorHandle(http.StatusBadRequest, err, "Wrong id format")
		return
	}
	var fact common.Fact
	err = json.NewDecoder(r.Body).Decode(&fact)
	if err != nil {
		errorHandler{w}.errorHandle(http.StatusBadRequest, err, "Wrong fact format")
		return
	}
	// todo add validation (probably while unmarshalling)
	err = db.Ins.UpdFact(r.Context(), fact, id)
	if err != nil {
		errorHandler{w}.errorHandle(http.StatusInternalServerError, err, "")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func parseId(s string) (int, error) {
	id := 0
	var c rune
	n, err := fmt.Scanf(s, "/fact/%d%c", &id, &c)
	if n == 1 && err == io.EOF && id > 0 {
		return id, nil
	}
	return 0, err
}
