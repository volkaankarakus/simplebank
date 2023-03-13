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

// ** TransferTx performs a money transfer from one account to another.
// ** It creates a transfer record, add account entries,
// **  and update account's balance within a single database transaction

// ** First, let's define the TransferTxParams
// ** This struct contains all necessary input parameters to
// **   transfer money between 2 accounts
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// ** TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{arg.FromAccountID, true},
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: sql.NullInt64{arg.ToAccountID, true},
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// TODO : update account balance ***********
		return nil

	})
	return result, err
}
