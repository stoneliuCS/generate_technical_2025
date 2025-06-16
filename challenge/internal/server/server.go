package server

import (
	api "generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/utils"
	"net/http"
)

// Runs the server api with the given handler.
func RunServer(handler api.Handler, addr string) {
	srv_func := func() (*api.Server, error) { return api.NewServer(handler) }
	srv := utils.SafeCall(srv_func)
	serve_func := func() error {
		return http.ListenAndServe(addr, srv)
	}
	utils.SafeCallErrorSupplier(serve_func)
}
