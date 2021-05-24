package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Cori", "Wins": 33}]`)

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Ryan", 10},
			{"Cori", 33},
		}

		assertLeague(t, got, want)
	})
}
