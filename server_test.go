package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := League{
			{"Ryan", 10},
			{"Cori", 11},
			{"Bill", 12},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := NewLeagueRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := GetLeagueFromResponse(t, response.Body)
		AssertStatus(t, response.Code, http.StatusOK)
		AssertLeague(t, got, wantedLeague)
		AssertContentType(t, response, jsonContentType)
	})
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Ryan":  20,
			"Harry": 10,
		},
		nil,
		nil,
	}

	server := NewPlayerServer(&store)

	t.Run("returns Ryan's score", func(t *testing.T) {
		request := NewGetScoreRequest("Ryan")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Harry's score", func(t *testing.T) {
		request := NewGetScoreRequest("Harry")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := NewGetScoreRequest("Gopher")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestScoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Ryan", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		AssertStatus(t, response.Code, http.StatusAccepted)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}

	server := NewPlayerServer(&store)
	t.Run("it records wins when POST", func(t *testing.T) {
		player := "Not Ryan"
		request := NewPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		AssertStatus(t, response.Code, http.StatusAccepted)
		AssertPlayerWin(t, &store, "Not Ryan")
	})
}
