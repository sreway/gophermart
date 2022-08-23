package httpserver

type Option func(*Server)

func Addr(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}
