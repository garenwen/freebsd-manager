package server

type Server interface {
	Start()
	Shutdown(code int, reason string)
}
