package rest_server

type ServerOption func(*Server)

func WithPort(port int) ServerOption {
	return func(server *Server) {
		server.port = port
	}
}

func WithTracerName(name string) ServerOption {
	return func(server *Server) {
		server.tracerName = name
	}
}

func WithMode(mode string) ServerOption {
	return func(server *Server) {
		server.mode = mode
	}
}

func WithHealthy(healthy bool) ServerOption {
	return func(server *Server) {
		server.healthy = healthy
	}
}

func WithEnableProfiling(enableProfiling bool) ServerOption {
	return func(server *Server) {
		server.enableProfiling = enableProfiling
	}
}

func WithMiddlewares(middlewares []string) ServerOption {
	return func(server *Server) {
		server.middlewares = middlewares
	}
}

func WithJwt(jwt *JwtInfo) ServerOption {
	return func(server *Server) {
		server.jwt = jwt
	}
}

func WithTransName(name string) ServerOption {
	return func(server *Server) {
		server.transName = name
	}
}
