package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/celest1al/simplebank-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) (Account, Entry) {
	account := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	return account, entry
}

func deleteTestEntry(t *testing.T, accountId int64, entryId int64) {
	testQueries.DeleteEntry(context.Background(), entryId)
	testQueries.DeleteAccount(context.Background(), accountId)
}

func TestCreateEntry(t *testing.T) {
	account, entry := createRandomEntry(t)

	deleteTestEntry(t, account.ID, entry.ID)
}

func TestGetEntry(t *testing.T) {
	account, entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

	deleteTestEntry(t, account.ID, entry1.ID)
}

func TestUpdateEntry(t *testing.T) {
	account, entry1 := createRandomEntry(t)

	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}

	entry2, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, arg.ID, entry2.ID)
	require.Equal(t, arg.Amount, entry2.Amount)

	deleteTestEntry(t, account.ID, entry1.ID)
}

func TestDeleteEntry(t *testing.T) {
	account, entry1 := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry1.ID)

	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)

	testQueries.DeleteAccount(context.Background(), account.ID)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  10,
		Offset: 1,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 10)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
