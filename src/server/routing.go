package server

func RegisterRoute(p *Path, m Method, h Handler) *Router {
	return R.HandleFunc(p, m, h)
}

func Get(p *Path, h Handler) *Router {
	return RegisterRoute(p, GET, h)
}

func Post(p *Path, h Handler) *Router {
	return RegisterRoute(p, POST, h)
}

func Put(p *Path, h Handler) *Router {
	return RegisterRoute(p, PUT, h)
}

func Patch(p *Path, h Handler) *Router {
	return RegisterRoute(p, PATCH, h)
}

func Delete(p *Path, h Handler) *Router {
	return RegisterRoute(p, DELETE, h)
}
