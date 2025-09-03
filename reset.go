package main

import "net/http"

func (cfg *apiConfig) handleResetServerHits(rw http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hits reset to 0"))
}
