package poker_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ryanyogan/poker"
)

var dummySpyAlerter = &SpyBlindAlerter{}

type scheduleAlert struct {
	at     time.Duration
	amount int
}

func (s scheduleAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduleAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduleAlert{at, amount})
}

func TestCLI(t *testing.T) {
	t.Run("record ryan win from user input", func(t *testing.T) {
		in := strings.NewReader("Ryan wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummySpyAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Ryan")
	})

	t.Run("record cori win from user input", func(t *testing.T) {
		in := strings.NewReader("Cori wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummySpyAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cori")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Ryan wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		if len(blindAlerter.alerts) != 11 {
			t.Fatal("expected a blind alert to be scheduled")
		}
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Ryan wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		testCases := []scheduleAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range testCases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled for %v", i, blindAlerter.alerts)
				}
			})
		}
	})
}
