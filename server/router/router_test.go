package router

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthOnly(t *testing.T) {
	r := SetupRouter("release", (*sql.DB)(nil))
	routes := r.Routes()
	if len(routes) != 2 || routes[0].Method != http.MethodGet || routes[0].Path != "/health" || routes[1].Path != "/api/v1/patients" {
		t.Fatalf("expected health and patient query routes, got %v", routes)
	}
	for _, tc := range []struct {
		path   string
		status int
		body   string
	}{
		{"/health", http.StatusOK, `{"status":"ok"}`},
		{"/api/v1/system/info", http.StatusNotFound, "404 page not found"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if w.Code != tc.status || w.Body.String() != tc.body {
				t.Fatalf("got %d %q, want %d %q", w.Code, w.Body.String(), tc.status, tc.body)
			}
		})
	}
}
