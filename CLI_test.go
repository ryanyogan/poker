package poker_test

import (
	"strings"
	"testing"

	"github.com/ryanyogan/poker"
)

func TestCLI(t *testing.T) {
	t.Run("record ryan win from user input", func(t *testing.T) {
		in := strings.NewReader("Ryan wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Ryan")
	})

	t.Run("record cori win from user input", func(t *testing.T) {
		in := strings.NewReader("Cori wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cori")
	})
}
