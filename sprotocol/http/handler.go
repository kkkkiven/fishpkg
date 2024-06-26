package http

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"

	"git.yuetanggame.com/zfish/fishpkg/utils"

	xutils "git.yuetanggame.com/zfish/fishpkg/servicesdk/pkg/utils"
	"git.yuetanggame.com/zfish/fishpkg/sprotocol/tracer"

	"git.yuetanggame.com/zfish/fishpkg/logs"
	"github.com/pkg/errors"
)

const (
	X_TRACE_ID     = "X-TRACE-ID"
	X_SPAN_ID      = "X-SPAN-ID"
	X_PROJECT_NAME = "X-PROJECT-NAME"
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

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request, Params)

func (this HandlerFunc) Serve(ctx context.Context, rw http.ResponseWriter, r *http.Request, params Params) {
	this(ctx, rw, r, params)
}

var (
	tree *node
	mu   sync.RWMutex
)

func init() {
	tree = new(node)
}

type Wrapper struct {
	Groups  []string
	Handler func(context.Context, http.ResponseWriter, *http.Request, Params)
}

// AddHandler 添加hander
func AddHandler(path string, handler func(context.Context, http.ResponseWriter, *http.Request, Params), groups ...string) {
	mu.Lock()
	defer mu.Unlock()

	w := &Wrapper{}
	w.Groups = append(w.Groups, groups...)
	w.Handler = handler

	tree.addRoute(path, w)
}

// routerHandler 路由：/server/iface
func routerHandler(so *Socket, msg *Message) error {
	p := msg.GetStringExData()
	w, ps, _ := tree.getValue(p)
	if w == nil || w.Handler == nil {
		rsp := NewResponseMessage()
		rsp.SetRequestID(msg.GetRequestID())
		rsp.SetStatus(RC_HANDLER_NOT_FOUND)
		so.Send(rsp)
		return errors.New(fmt.Sprintf("Handler[%v] not found", p))
	}

	// 调用函数
	go doCall(w, ps, so, msg)

	return nil
}

// doCall 执行调用过程
func doCall(w *Wrapper, params Params, so *Socket, msg *Message) {
	var (
		h   Handler
		r   *http.Request
		rw  *ResponseWriter
		req *HTTPRequest = &HTTPRequest{}
		err error
	)

	// 构造调用链
	h = chain(HandlerFunc(w.Handler), w.Groups...)

	// 调用过程
	if _, err = req.UnmarshalMsg(msg.GetBody()); err != nil {
		rsp := NewResponseMessage()
		rsp.SetRequestID(msg.GetRequestID())
		rsp.SetStatus(RC_SYS_ERR)
		so.Send(rsp)
		logs.Errorf("- %v - Decode http request err: %v", so.GetConn().RemoteAddr().String(), err.Error())
		return
	}

	if r, err = req.ReadRequest(); err != nil {
		rsp := NewResponseMessage()
		rsp.SetRequestID(msg.GetRequestID())
		rsp.SetStatus(RC_SYS_ERR)
		so.Send(rsp)
		logs.Errorf("- %v - Decode http request err: %v", so.GetConn().RemoteAddr().String(), err.Error())
		return
	}

	rw = NewResponseWriter()
	ctx, span := genSpan(r)
	defer func() {
		body, _ := rw.MarshalMsg(nil)
		rsp := NewResponseMessage()
		rsp.SetRequestID(msg.GetRequestID())
		rsp.SetBody(body)

		if err := recover(); err != nil {
			rsp.SetStatus(RC_HANDLER_PANIC)

			hint := fmt.Sprintf("Panic: %+v\n%s", err, string(debug.Stack()))
			span.Tag("msg", hint)
			span.Tag("code", http.StatusInternalServerError)
			span.End()

			logs.Errorf("- %v - %s", so.GetConn().RemoteAddr().String(), hint)
		} else {
			span.Tag("code", rw.Status)
			span.End()
		}

		so.Send(rsp)
	}()

	h.Serve(ctx, rw, r, params)
}

func genSpan(r *http.Request) (ctx context.Context, span *tracer.Span) {
	var tid, pid uint64

	tids := r.Header.Get(strings.ToUpper(X_TRACE_ID))
	if tids != "" {
		if len(tids) < 32 {
			tids = xutils.Get16MD5Encode(tids)
		} else {
			tids = tids[8:24]
		}

		tidb, err := hex.DecodeString(tids)
		if err == nil && len(tidb) > 0 {
			tid = binary.BigEndian.Uint64(tidb)
		}
	}

	pids := r.Header.Get(strings.ToUpper(X_SPAN_ID))
	if pids != "" {
		if len(pids) < 32 {
			pids = xutils.Get16MD5Encode(pids)
		} else {
			pids = pids[8:24]
		}

		pidb, err := hex.DecodeString(pids)
		if err == nil && len(pidb) > 0 {
			pid = binary.BigEndian.Uint64(pidb)
		}
	}

	if tid != 0 {
		span, _ = tracer.CreateProbSpan(int64(tid), int64(pid))
		span.Tag("method", r.Method)
		span.Tag("path", r.RequestURI)
		addrs := strings.Split(r.RemoteAddr, ":")
		if len(addrs) == 2 {
			project := r.Header.Get(strings.ToUpper(X_PROJECT_NAME))
			if project == "" {
				project = "nginx"
			}

			span.SetRemoteEndpoint(project, 0, 0, addrs[0], uint16(utils.AtoUi(addrs[1])))
		}

		_, ctx = tracer.CreatePropagateSpan(int64(tid), span.GetSpanID())
	}

	return context.TODO(), nil
}
