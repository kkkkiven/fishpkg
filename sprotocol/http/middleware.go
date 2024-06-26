package http

import (
	"context"
	"net/http"
)

type Handler interface {
	Serve(context.Context, http.ResponseWriter, *http.Request, Params)
}

var globalMiddleware []func(Handler) Handler
var groupMiddlerware map[string][]func(Handler) Handler

// 初始化中间件
func init() {
	globalMiddleware = make([]func(Handler) Handler, 0)
	groupMiddlerware = make(map[string][]func(Handler) Handler, 0)
}

// Use 设置全局中间件
func Use(mws ...func(Handler) Handler) {
	globalMiddleware = append(globalMiddleware, mws...)
}

// Group 设置分组中间件
func Group(group string, mws ...func(Handler) Handler) {
	if _, ok := groupMiddlerware[group]; !ok {
		groupMiddlerware[group] = make([]func(Handler) Handler, 0)
	}

	groupMiddlerware[group] = append(groupMiddlerware[group], mws...)
}

// chain 构造调用链
func chain(endpoint Handler, groups ...string) Handler {
	var h Handler
	for i := 0; i < len(groups); i++ {
		if middleware, ok := groupMiddlerware[groups[i]]; ok {
			if h == nil {
				h = middleware[len(middleware)-1](endpoint)
				for j := len(middleware) - 2; j >= 0; j-- {
					h = middleware[j](h)
				}
			} else {
				for j := len(middleware) - 1; j >= 0; j-- {
					h = middleware[j](h)
				}
			}
		}
	}

	if len(globalMiddleware) != 0 {
		if h == nil {
			h = globalMiddleware[len(globalMiddleware)-1](endpoint)
			for i := len(globalMiddleware) - 2; i >= 0; i-- {
				h = globalMiddleware[i](h)
			}
		} else {
			for i := len(globalMiddleware) - 1; i >= 0; i-- {
				h = globalMiddleware[i](h)
			}
		}
	}

	if h == nil {
		return endpoint
	}
	return h
}
