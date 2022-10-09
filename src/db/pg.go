package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"httpServer/src/common"
	"log"
	"os"
	"strings"
)

// Instance type keeps a pool of connections to db, and have methods to work with it.
type Instance struct {
	Db *pgxpool.Pool
}

// Ins helps as get our db logic across all codebase.
var Ins Instance

func getConnStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PWD"),
		os.Getenv("DB_NAME"))
}

// InitDb is an initialization func for Ins variable.
// Creates a pool of connection to db and ping it.
func InitDb(ctx context.Context) *pgxpool.Pool {
	poolConfig, _ := pgxpool.ParseConfig(getConnStr())
	poolConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("db ping failed: %v", err)
	}
	return pool
}

// GetFactByID returns a common.Fact from db with given id. If there is no such id, returns pgx.ErrNoRows
func (ins Instance) GetFactByID(ctx context.Context, id int) (common.Fact, error) {
	var fact common.Fact
	row := ins.Db.QueryRow(ctx,
		`SELECT f.id, f.title, f.description, STRING_AGG(l.link, ',') FROM Facts AS f
INNER JOIN Links AS l ON f.id = l.fact_id
WHERE f.id = $1
GROUP BY f.id;`, id)
	var linkList string
	err := row.Scan(&fact.ID, &fact.Title, &fact.Desc, &linkList)
	if err != nil {
		return common.Fact{}, err
	}
	fact.Links = strings.Split(linkList, ",")
	return fact, nil
}

// UpdFact updates facts with given id.
func (ins Instance) UpdFact(ctx context.Context, fact common.Fact, id int) error {
	tx, err := ins.Db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = ins.Db.Exec(ctx,
		`UPDATE Facts SET title = $1, description = $2 WHERE id = $3;`,
		fact.Title,
		fact.Desc,
		id)
	if err != nil {
		return err
	}
	_, err = ins.Db.Exec(ctx,
		`DELETE FROM Links WHERE fact_id = $1;`, id)
	if err != nil {
		return err
	}
	for _, v := range fact.Links {
		_, err = ins.Db.Exec(ctx,
			`INSERT INTO Links(fact_id, link) VALUES ($1, $2)`, id, v)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

// InsertFacts inserts common.FactsArr into db and returns slice of ids
func (ins Instance) InsertFacts(ctx context.Context, facts common.FactsArr) ([]int, error) {
	tx, err := ins.Db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	res := make([]int, 0, len(facts.Facts))
	id := 0
	for _, v := range facts.Facts {
		err = ins.Db.QueryRow(ctx,
			`INSERT INTO Facts(title, description) VALUES ($1, $2) RETURNING id;`,
			v.Title,
			v.Desc).Scan(&id)
		if err != nil {
			return nil, err
		}
		for _, l := range v.Links {
			_, err = ins.Db.Exec(ctx,
				`INSERT INTO Links(fact_id, link) VALUES ($1, $2)`, id, l)
			if err != nil {
				return nil, err
			}
		}
		res = append(res, id)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetMaxID returns current max id in db
func (ins Instance) GetMaxID(ctx context.Context) (int, error) {
	id := 0
	err := ins.Db.QueryRow(ctx, `SELECT MAX(id) FROM Facts;`).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
