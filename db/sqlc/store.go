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
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// ** TXKEY
var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	txName := ctx.Value(txKey)

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// TODO : update account balance
		// ** get accounts -> update its balance

		// ** money out of the FromAccount

		// fmt.Println(txName,"get account for update")
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName,"update account")
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount,result.ToAccount,err = addMoney(ctx,q,arg.FromAccountID,- arg.Amount,arg.ToAccountID,arg.Amount)
		} else {
			result.ToAccount,result.FromAccount, err = addMoney(ctx,q,arg.ToAccountID,arg.Amount,arg.FromAccountID,-arg.Amount)
			

		}

		// ** money into to ToAccount
		// fmt.Println(txName,"get account for update 2")
		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName,"update account 2")
		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		return nil

	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
