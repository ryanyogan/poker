package poker

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
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

		assertLeague(t, got, want)

		// read again, confirm we remove the file and write a new one
		got = store.GetLeague()
		assertLeague(t, got, want)
		assertNoError(t, err)
	})

	t.Run("get player score", func(t *testing.T) {
		got := store.GetPlayerScore("Ryan")
		want := 10
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		store.RecordWin("Ryan")

		got := store.GetPlayerScore("Ryan")
		want := 11
		assertScoreEquals(t, got, want)
		assertNoError(t, err)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		store.RecordWin("New Player")

		got := store.GetPlayerScore("New Player")
		want := 1
		assertScoreEquals(t, got, want)
		assertNoError(t, err)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Cori", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetLeague()
		want := []Player{
			{"Cori", 33},
			{"Ryan", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()
		assertLeague(t, got, want)
	})
}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but received one, %v", err)
	}
}
