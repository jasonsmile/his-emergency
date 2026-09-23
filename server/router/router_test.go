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
	if len(routes) != 16 {
		t.Fatalf("expected health, login and registration routes, got %v", routes)
	}
	routeSet := make(map[string]bool, len(routes))
	for _, route := range routes {
		routeSet[route.Method+" "+route.Path] = true
	}
	for _, route := range []string{
		http.MethodGet + " /health",
		http.MethodPost + " /api/login",
		http.MethodGet + " /api/registration/searchPatients",
		http.MethodGet + " /api/registration/getPatientDetail",
		http.MethodGet + " /api/registration/getDepts",
		http.MethodGet + " /api/registration/getDoctors",
		http.MethodGet + " /api/registration/getSchedules",
		http.MethodPost + " /api/registration/createRegistration",
		http.MethodGet + " /api/registration/getEncounterList",
		http.MethodGet + " /api/registration/getEncounterDetail",
		http.MethodPost + " /api/registration/cancelEncounter",
		http.MethodPost + " /api/charge/getPendingCharges",
		http.MethodPost + " /api/charge/createCharge",
		http.MethodPost + " /api/charge/refundCharge",
		http.MethodGet + " /api/charge/getChargeRecords",
		http.MethodGet + " /api/charge/getDailyReport",
	} {
		if !routeSet[route] {
			t.Fatalf("missing route %s; got %v", route, routes)
		}
	}
	for _, tc := range []struct {
		path   string
		status int
		body   string
	}{
		{"/health", http.StatusOK, `{"status":"ok"}`},
		{"/api/patients", http.StatusNotFound, "404 page not found"},
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
