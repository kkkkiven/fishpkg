// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package db

import (
	"context"
	"database/sql"
)

var (
	globalInterceptors = make([]Interceptor, 0)

	_ Interceptor = &EmptyInterceptor{}
)

// AddInterceptor 添加全局拦截器, 越后添加的越先执行
func AddInterceptor(i Interceptor) {
	globalInterceptors = append(globalInterceptors, i)
}

type (
	BeginContextFunc     func(ctx context.Context) (*sql.Tx, error)
	InsertContextFunc    func(ctx context.Context, query string, args ...interface{}) (int64, error)
	UpdateContextFunc    func(ctx context.Context, query string, args ...interface{}) (int64, error)
	DeleteContextFunc    func(ctx context.Context, query string, args ...interface{}) (int64, error)
	SelectContextFunc    func(ctx context.Context, query string, args ...interface{}) (Results, error)
	SelectOneContextFunc func(ctx context.Context, query string, args ...interface{}) (OneRow, error)

	ExecContextFunc     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContextFunc    func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContextFunc func(ctx context.Context, query string, args ...interface{}) *sql.Row

	TxExecContextFunc     func(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (sql.Result, error)
	TxQueryContextFunc    func(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error)
	TxQueryRowContextFunc func(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) *sql.Row

	TxPrepareContextFunc func(ctx context.Context, tx *sql.Tx, query string) (*sql.Stmt, error)
	TxStmtContextFunc    func(ctx context.Context, tx *sql.Tx, stmt *sql.Stmt) *sql.Stmt
)

// Interceptor 拦截器接口
type Interceptor interface {
	BeginContext(BeginContextFunc) BeginContextFunc
	InsertContext(InsertContextFunc) InsertContextFunc
	UpdateContext(UpdateContextFunc) UpdateContextFunc
	DeleteContext(DeleteContextFunc) DeleteContextFunc
	SelectContext(SelectContextFunc) SelectContextFunc
	SelectOneContext(SelectOneContextFunc) SelectOneContextFunc

	// 以下接口是 builder 中会用到的

	ExecContext(ExecContextFunc) ExecContextFunc
	QueryContext(QueryContextFunc) QueryContextFunc
	QueryRowContext(QueryRowContextFunc) QueryRowContextFunc

	TxExecContext(TxExecContextFunc) TxExecContextFunc
	TxQueryContext(TxQueryContextFunc) TxQueryContextFunc
	TxQueryRowContext(TxQueryRowContextFunc) TxQueryRowContextFunc

	TxPrepareContext(TxPrepareContextFunc) TxPrepareContextFunc
	TxStmtContext(TxStmtContextFunc) TxStmtContextFunc
}

type EmptyInterceptor struct {
}

func (i *EmptyInterceptor) BeginContext(f BeginContextFunc) BeginContextFunc {
	return f
}

func (i *EmptyInterceptor) ExecContext(f ExecContextFunc) ExecContextFunc {
	return f
}

func (i *EmptyInterceptor) InsertContext(f InsertContextFunc) InsertContextFunc {
	return f
}

func (i *EmptyInterceptor) UpdateContext(f UpdateContextFunc) UpdateContextFunc {
	return f
}

func (i *EmptyInterceptor) DeleteContext(f DeleteContextFunc) DeleteContextFunc {
	return f
}

func (i *EmptyInterceptor) QueryContext(f QueryContextFunc) QueryContextFunc {
	return f
}

func (i *EmptyInterceptor) QueryRowContext(f QueryRowContextFunc) QueryRowContextFunc {
	return f
}

func (i *EmptyInterceptor) SelectContext(f SelectContextFunc) SelectContextFunc {
	return f
}

func (i *EmptyInterceptor) SelectOneContext(f SelectOneContextFunc) SelectOneContextFunc {
	return f
}

func (i *EmptyInterceptor) TxExecContext(f TxExecContextFunc) TxExecContextFunc {
	return f
}

func (i *EmptyInterceptor) TxQueryContext(f TxQueryContextFunc) TxQueryContextFunc {
	return f
}

func (i *EmptyInterceptor) TxQueryRowContext(f TxQueryRowContextFunc) TxQueryRowContextFunc {
	return f
}

func (i *EmptyInterceptor) TxPrepareContext(f TxPrepareContextFunc) TxPrepareContextFunc {
	return f
}

func (i *EmptyInterceptor) TxStmtContext(f TxStmtContextFunc) TxStmtContextFunc {
	return f
}
