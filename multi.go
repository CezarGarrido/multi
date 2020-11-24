package multi

import (
	"database/sql"
)

type multiTxFn func(*sql.Tx) (Result, error)

// Multi :
type Multi struct {
}

// New : return a new multi
func New() *Multi {
	return &Multi{}
}

// Run : ...
func (mt *Multi) Run(name string, db *sql.DB, fn multiTxFn) (interface{}, error) {
	tx, _ := db.Begin()
	value, err := fn(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return value, nil
}

// Result :
type Result map[string]interface{}
