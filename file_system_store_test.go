package poker

import (
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Cori", "Wins": 33}]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)

	t.Run("league from a reader", func(t *testing.T) {
		got := store.GetLeague()
		want := League{
			{"Cori", 33},
			{"Ryan", 10},
		}

		AssertLeague(t, got, want)

		// read again, confirm we remove the file and write a new one
		got = store.GetLeague()
		AssertLeague(t, got, want)
		AssertNoError(t, err)
	})

	t.Run("get player score", func(t *testing.T) {
		got := store.GetPlayerScore("Ryan")
		want := 10
		AssertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		store.RecordWin("Ryan")

		got := store.GetPlayerScore("Ryan")
		want := 11
		AssertScoreEquals(t, got, want)
		AssertNoError(t, err)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		store.RecordWin("New Player")

		got := store.GetPlayerScore("New Player")
		want := 1
		AssertScoreEquals(t, got, want)
		AssertNoError(t, err)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Cori", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		got := store.GetLeague()
		want := []Player{
			{"Cori", 33},
			{"Ryan", 10},
		}

		AssertLeague(t, got, want)

		got = store.GetLeague()
		AssertLeague(t, got, want)
	})
}
