package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDraftTradeRoute(t *testing.T) {
	handler := handleRequests()
	for _, test := range []struct {
		method string
		body   string
		want   int
	}{
		{http.MethodPost, "{", http.StatusBadRequest},
		{http.MethodPost, "{}", http.StatusBadRequest},
		{http.MethodGet, "", http.StatusNotFound},
	} {
		req := httptest.NewRequest(test.method, "/api/trades/draft/accept/trade/process", strings.NewReader(test.body))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		if recorder.Code != test.want {
			t.Errorf("%s %q: got %d, want %d", test.method, test.body, recorder.Code, test.want)
		}
	}
}
