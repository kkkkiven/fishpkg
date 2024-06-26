package http

import (
	"context"
	"net/http"
	"unsafe"

	p "github.com/kkkkiven/fishpkg/sprotocol/http"
)

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

// AddHandler 注册路由与处理函数
func AddHandler(path, desc string, fn func(context.Context, http.ResponseWriter, *http.Request, Params), groups ...string) {
	fp := *(*p.HandlerFunc)(unsafe.Pointer(&fn))
	p.AddHandler(path, fp, groups...)
	iface := _Iface{
		Path: path,
		Desc: desc,
	}
	srv.AddIface(iface)
}

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request, Params)

func (this HandlerFunc) Serve(ctx context.Context, rw http.ResponseWriter, r *http.Request, params Params) {
	this(ctx, rw, r, params)
}

// Handler
type Handler interface {
	Serve(context.Context, http.ResponseWriter, *http.Request, Params)
}

// middlewareFunc
type middlewareFunc func(p.Handler) p.Handler

// AddMiddleware
func AddMiddleware(fn func(Handler) Handler) {
	fp := *(*middlewareFunc)(unsafe.Pointer(&fn))
	p.Use(fp)
}

// AddGroupMiddleware
func AddGroupMiddleware(name string, fn func(Handler) Handler) {
	fp := *(*middlewareFunc)(unsafe.Pointer(&fn))
	p.Group(name, fp)
}
