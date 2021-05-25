package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Cori", "Wins": 33}]`)
	defer cleanDatabase()

	store := NewFileSystemPlayerStore(database)

	t.Run("league from a reader", func(t *testing.T) {
		got := store.GetLeague()
		want := League{
			{"Ryan", 10},
			{"Cori", 33},
		}

		assertLeague(t, got, want)

		// read again, confirm we remove the file and write a new one
		got = store.GetLeague()
		assertLeague(t, got, want)
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
	})

	t.Run("store wins for new players", func(t *testing.T) {
		store.RecordWin("New Player")

		got := store.GetPlayerScore("New Player")
		want := 1
		assertScoreEquals(t, got, want)
	})
}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
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
