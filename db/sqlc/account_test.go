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
func createRandomAccount(t *testing.T) Account {
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

// ** Test Create Account
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

// ** Test Get Account
func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, error := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, error)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	// ** for the timestamp fields like created_at, beside require.Equal()
	// **  ->  require.WithinDuration(). to check that 2 timestamp are different by at most some delta duration
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second) // delta to be 1 second
}


// ** Test Update Account
func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	account2, error := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, error)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

	require.Equal(t, arg.Balance, account2.Balance)
}

// ** Test Delete Account
func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	error := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, error)

	// ** is account really deleted?
	account2, error2 := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, error2)
	require.EqualError(t, error2, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}


// ** Test List Accounts
func TestListAccounts(t *testing.T) {
	// ** create 10 random account
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{ // ** which means skip the first 5 records and return the next 5.
		Limit:  5,
		Offset: 5,
	}

	// ** when we run the test, there will be at least 10 accounts in the database, so with these parameters we expect to get 5 records
	accounts, error := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, error)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
