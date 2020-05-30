package main

import "net/http"

type Server struct {
	*router
	middlewares  []Middleware
	startHanlder HandlerFunc
}

func NewServer() *Server {
	r := &router{make(map[string]map[string]HandlerFunc)}
	s := &Server{router: r}
	s.middlewares = []Middleware{
		logHandler,
		recoverHandler,
		staticHandler,
		parseFormHandler,
		parseJSONBodyHandler}
	return s
}

func (s *Server) Run(addr string) {
	s.startHanlder = s.router.handler()

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.startHanlder = s.middlewares[i](s.startHanlder)
	}

	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		Params:         make(map[string]interface{}),
		ResponseWriter: w,
		Request:        r,
	}
	for k, v := range r.URL.Query() {
		c.Params[k] = v[0]
	}
	s.startHanlder(c)
}

func (s *Server) Use(middleware ...Middleware) {
	s.middlewares = append(s.middlewares, middleware...)
}
