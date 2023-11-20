package server

import "net/http"

type Server struct {
	httpServer *http.Server
}

func (s *Server) RunServer(host string, port string, controller http.Handler) error {

	s.httpServer = &http.Server{
		Addr:    host + port,
		Handler: controller,
	}

	return s.httpServer.ListenAndServe()
}
