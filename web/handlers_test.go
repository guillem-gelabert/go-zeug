package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	dto "github.com/guillem-gelabert/go-zeug/web/dtos"
)

var mockSessionDTO []dto.CardDTO = []dto.CardDTO{
	{ID: 0, Article: "das", Substantive: "Haus", WordID: 0, UserID: 1, Stage: "SEEN"},
	{ID: 1, Article: "das", Substantive: "Jahr", WordID: 1, UserID: 1, Stage: "UNSEEN"},
	{ID: 2, Article: "das", Substantive: "Prozent", WordID: 2, UserID: 1, Stage: "SEEN"},
	{ID: 3, Article: "der", Substantive: "Euro", WordID: 3, UserID: 1, Stage: "UNSEEN"},
	{ID: 4, Article: "die", Substantive: "Zeit", WordID: 4, UserID: 1, Stage: "SEEN"},
	{ID: 5, Article: "die", Substantive: "Kategorie", WordID: 5, UserID: 1, Stage: "UNSEEN"},
	{ID: 6, Article: "die", Substantive: "Stadt", WordID: 6, UserID: 1, Stage: "SEEN"},
	{ID: 7, Article: "das", Substantive: "Ende", WordID: 7, UserID: 1, Stage: "UNSEEN"},
	{ID: 8, Article: "die", Substantive: "Frau", WordID: 8, UserID: 1, Stage: "SEEN"},
	{ID: 9, Article: "das", Substantive: "Leben", WordID: 9, UserID: 1, Stage: "UNSEEN"},
	{ID: 10, Article: "das", Substantive: "Leben", WordID: 10, UserID: 1, Stage: "SEEN"},
}

func TestGetSession(t *testing.T) {
	app := newTestApplication(t)

	testCases := []struct {
		desc           string
		ID             int
		expectedStatus int
		expectedBody   *[]dto.CardDTO
	}{
		{
			desc:           "GET /cards with valid ID",
			ID:             1,
			expectedStatus: http.StatusOK,
			expectedBody:   &mockSessionDTO,
		},
		{
			desc:           "GET /cards with invalid ID",
			ID:             2,
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			desc:           "GET /cards with valid ID but no cards scheduled",
			ID:             3,
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			app.loggedIn.ID = tC.ID
			ts := newTestServer(t, app.routes())
			defer ts.Close()

			code, _, body := ts.get(t, "/cards")
			if code != tC.expectedStatus {
				t.Errorf("Expected %d; got %d", tC.expectedStatus, code)
			}

			if tC.expectedBody == nil {
				if !reflect.DeepEqual([]byte(""), body) {
					t.Errorf("Expected \"\"; got %q", body)
				}
				return
			}

			var actual *[]dto.CardDTO
			err := json.Unmarshal(body, &actual)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, tC.expectedBody) {
				t.Errorf("Expected %#v; got %#v", tC.expectedBody, actual)
			}
		})
	}
}
