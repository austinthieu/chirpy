package main

import (
	"net/http"
)

func (cfg *apiConfig) handleResetUsers(rw http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to reset the database: " + err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
