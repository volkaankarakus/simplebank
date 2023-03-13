package db

import (
	"context"
	"database/sql"
	"fmt"
)

// ** Store provides all functions to execute DB queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// ** New Store creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// ** We will add a function to the Store to execute a generic database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)

	if err != nil {
		if rollbackError := tx.Rollback(); rollbackError != nil {
			return fmt.Errorf("tx error : %v, rollback error : %v", err, rollbackError)
		}
		return err
	}
	return tx.Commit()
}
