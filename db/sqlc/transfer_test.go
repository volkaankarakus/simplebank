package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

// ** for testing we should create an transfer first
func createRandomTransfer(t *testing.T) Transfer {

	// ** Create Account
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// ** Create Account Get Tests
	account1GetTest, error1 := testQueries.GetAccount(context.Background(),account1.ID)
	require.NoError(t, error1)
	require.NotEmpty(t, account1GetTest)
	require.Equal(t, account1GetTest.ID, account1.ID)

	account2GetTest, error2 := testQueries.GetAccount(context.Background(),account2.ID)
	require.NoError(t, error2)
	require.NotEmpty(t, account2GetTest)
	require.Equal(t, account2GetTest.ID, account2.ID)

	// ** Create Transfer
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: util.RandomMoney(),
	}

	transfer, transferError := testQueries.CreateTransfer(context.Background(),arg)
	require.NoError(t,transferError)
	require.NotEmpty(t,transfer)
	require.Equal(t,arg.ToAccountID,transfer.ToAccountID)
	require.Equal(t,arg.FromAccountID,transfer.FromAccountID)
	require.Equal(t,arg.Amount,transfer.Amount)

	return transfer
}


// ** Test Create Transfer
func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

// ** Test Get Transfer
func TestGetTransfer(t *testing.T){
	transfer1 := createRandomTransfer(t)
	transfer2, transferError := testQueries.GetTransfer(context.Background(),transfer1.ID)

	require.NoError(t,transferError)
	require.NotEmpty(t,transfer2)

	require.Equal(t,transfer1.ID,transfer2.ID)
	require.Equal(t,transfer1.FromAccountID,transfer2.FromAccountID)
	require.Equal(t,transfer1.ToAccountID,transfer2.ToAccountID)
	require.Equal(t,transfer1.Amount,transfer2.Amount)
	require.WithinDuration(t,transfer1.CreatedAt,transfer2.CreatedAt,time.Second)

}

// ** Test Update Transfer
func TestUpdateTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	arg := UpdateTransferParams{
		ID: transfer1.ID,
		Amount: util.RandomAmount(),
	}

	transfer2, transferError := testQueries.UpdateTransfer(context.Background(),arg)
	require.NoError(t,transferError)
	require.NotEmpty(t,transfer2)

	require.Equal(t,transfer1.ID,transfer2.ID)
	require.Equal(t,transfer1.ToAccountID,transfer2.ToAccountID)
	require.Equal(t,transfer1.FromAccountID,transfer2.FromAccountID)
	require.WithinDuration(t,transfer1.CreatedAt,transfer2.CreatedAt,time.Second)

	require.Equal(t,arg.Amount,transfer2.Amount)
}

// ** Test Delete Transfer
func TestDeleteTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transferError := testQueries.DeleteTransfer(context.Background(),transfer1.ID)
	require.NoError(t,transferError)

	// ** is transfer really deleted?
	transfer2, transferError2 := testQueries.GetTransfer(context.Background(),transfer1.ID)
	require.Error(t,transferError2)
	require.EqualError(t,transferError2,sql.ErrNoRows.Error())
	require.Empty(t,transfer2)
}

