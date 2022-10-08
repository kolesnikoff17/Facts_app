package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"httpServer/src/common"
	"reflect"
	"testing"
)

func isEqual(a, b common.Fact) bool {
	if a.Id != b.Id || a.Title != b.Title || a.Desc != b.Desc || len(a.Links) != len(b.Links) {
		return false
	}
	for i := 0; i < len(a.Links) && i < len(b.Links); i++ {
		if a.Links[i] != b.Links[i] {
			return false
		}
	}
	return true
}

func TestInitDb(t *testing.T) {
	ctx := context.Background()
	Ins.Db = InitDb(ctx)
}

func TestInstance_GetFactById(t *testing.T) {
	ctx := context.Background()
	factEx := common.Fact{Id: 1, Title: "aboba", Desc: "aboba_aboba", Links: []string{"aboba", "aboba_aboba"}}
	fact, err := Ins.GetFactById(ctx, 1)
	if err != nil {
		t.Fatalf("get fact by id test: expected err nil, got %v", err)
	}
	if !isEqual(factEx, fact) {
		t.Errorf("get fact by id test: expected fact %v, got %v", factEx, fact)
	}
	_, err = Ins.GetFactById(ctx, -1)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("get fact by id test: expected err %v, got %v", pgx.ErrNoRows, err)
	}
}

func TestInstance_GetMaxId(t *testing.T) {
	ctx := context.Background()
	expected := 3
	got, err := Ins.GetMaxId(ctx)
	if err != nil {
		t.Fatalf("get max id test: expected err nil, got %v", err)
	}
	if expected != got {
		t.Errorf("get max id test: expected %d, got %d", expected, got)
	}
}

func TestInstance_InsertFacts(t *testing.T) {
	ctx := context.Background()
	expected := []int{4, 5}
	facts := common.FactsArr{Facts: []common.Fact{{Title: "ab", Desc: "ab", Links: []string{"ab"}},
		{Title: "bd", Desc: "bd"}}}
	got, err := Ins.InsertFacts(ctx, facts)
	if err != nil {
		t.Fatalf("insert facts test: expected err nil, got %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("insert facts test: expected %v, got %v", expected, got)
	}
}

func TestInstance_InsertFacts2(t *testing.T) {
	ctx := context.Background()
	got, err := Ins.InsertFacts(ctx, common.FactsArr{})
	if err != nil {
		t.Fatalf("insert empty facts test: expected err nil, got %v", err)
	}
	if !reflect.DeepEqual(got, []int{}) {
		t.Errorf("insert empty facts test: expected %v, got %v", []int{}, got)
	}
}

func TestInstance_UpdFact(t *testing.T) {
	ctx := context.Background()
	fact := common.Fact{Id: 2, Title: "hhh", Desc: "hhh", Links: []string{"hhh"}}
	err := Ins.UpdFact(ctx, fact, 2)
	if err != nil {
		t.Fatalf("update fact test: expected err nil, got %v", err)
	}
}

func TestInstance_GetFactById2(t *testing.T) {
	ctx := context.Background()
	factEx := common.Fact{Id: 2, Title: "hhh", Desc: "hhh", Links: []string{"hhh"}}
	fact, err := Ins.GetFactById(ctx, 2)
	if err != nil {
		t.Fatalf("get fact by id after update test: expected err nil, got %v", err)
	}
	if !isEqual(factEx, fact) {
		t.Errorf("get fact by id after update test: expected fact %v, got %v", factEx, fact)
	}
}
