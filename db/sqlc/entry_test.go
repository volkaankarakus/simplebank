package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

// ** for testing we should create an account first
func createRandomAccountForEntry(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, error := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, error)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

// ** test create account for entry
func TestCreateAccountForEntry(t *testing.T) {
	createRandomAccountForEntry(t)
}

// ** test get account for entry
func TestGetAccountForEntry(t *testing.T) {
	account1 := createRandomAccountForEntry(t)
	account2, error := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, error)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

// ** create random entry
func createRandomEntry(t *testing.T) Entry {
	account1 := createRandomAccountForEntry(t)
	arg := CreateEntryParams{
		AccountID: sql.NullInt64{account1.ID, true},
		Amount:    util.RandomAmount(),
	}

	entry, error := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, error)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

// ** test create entry
func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

// ** test get entry
func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, error := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, error)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

// ** test update entry
func TestUpdateEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomAmount(),
	}

	entry2, error := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, error)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

	require.Equal(t, arg.Amount, entry2.Amount)
}

// ** test delete entry
func TestDeleteEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	error := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, error)

	// ** is entry really deleted?
	entry2, error2 := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, error2)
	require.EqualError(t, error2, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

// ** test list entries
func TestListEntries(t *testing.T) {
	// ** create 10 random entries
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{ // ** skip the first 3 records and return the next 5
		Limit:  5,
		Offset: 3,
	}

	// ** when we run the test, there will be at least 10 entries in the database, so with these parameters we expect to get 5 records
	accounts, error := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, error)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
