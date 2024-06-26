// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package db

import (
	"context"
	"database/sql"
)

type ctxKey string

const (
	ctxKeyTx ctxKey = "_godb_tx"
)

var _ Tx = &dbtx{}

type Tx interface {
	RawTx() *sql.Tx

	Insert(ignore ...bool) *SB
	Delete() *SB
	Update() *SB
	InsertUpdate() *SB
	Select(str ...string) *SB
	SelectResults(query string, args ...interface{}) (Results, error)
	SelectOne(query string, args ...interface{}) (OneRow, error)

	// 以下为标准库 *sql.Tx 自带方法

	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Stmt(stmt *sql.Stmt) *sql.Stmt
	StmtContext(ctx context.Context, stmt *sql.Stmt) *sql.Stmt
}

// Transaction 执行事务代理
// f 返回错误时会自动 rollback
// 支持嵌套事务, 但最终提交已最外层的事务为准, 嵌套过程必须传递 context.Context 到嵌套的事务中
func Transaction(
	ctx context.Context,
	f func(context.Context, Tx) error,
	d ...*Database) (err error) {

	tx, err := getTx(ctx, d...)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.tx.Rollback()
		} else if err != nil {
			_ = tx.tx.Rollback()
		} else {
			tx.depth--
			if tx.depth > 0 {
				return
			}
			err = tx.tx.Commit()
		}
	}()

	tx.depth++
	err = f(tx.ctx, tx)
	return
}

// ContextTransaction 执行事务代理, 与 Transaction 逻辑一致, 主要区别是通过 context 传递 tx
// 支持嵌套事务, 但最终提交已最外层的事务为准, 嵌套过程必须传递 context.Context 到嵌套的事务中
func ContextTransaction(
	ctx context.Context,
	f func(ctx context.Context) error,
	d ...*Database) (err error) {

	tx, err := getTx(ctx, d...)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.tx.Rollback()
		} else if err != nil {
			_ = tx.tx.Rollback()
		} else {
			tx.depth--
			if tx.depth > 0 {
				return
			}
			err = tx.tx.Commit()
		}
	}()

	tx.depth++
	err = f(tx.ctx)
	return
}

func getTx(ctx context.Context, d ...*Database) (*dbtx, error) {
	if v := ctx.Value(ctxKeyTx); v != nil {
		if t, ok := v.(*dbtx); ok {
			return t, nil
		}
	}

	tx, err := Begin(d...)
	if err != nil {
		return nil, err
	}

	t := &dbtx{
		db: getDatabase(d...),
		tx: tx,
	}
	t.ctx = context.WithValue(ctx, ctxKeyTx, t)
	return t, nil
}

func hasTx(ctx context.Context) *dbtx {
	if v := ctx.Value(ctxKeyTx); v != nil {
		if t, ok := v.(*dbtx); ok {
			return t
		}
	}
	return nil
}

type dbtx struct {
	ctx   context.Context
	db    *Database
	tx    *sql.Tx
	depth int
}

func (t *dbtx) RawTx() *sql.Tx {
	return t.tx
}

func (t *dbtx) Insert(ignore ...bool) *SB {
	return Insert(ignore...).Tx(t.tx)
}

func (t *dbtx) Delete() *SB {
	return Delete().Tx(t.tx)
}

func (t *dbtx) Update() *SB {
	return Update().Tx(t.tx)
}

func (t *dbtx) InsertUpdate() *SB {
	return InsertUpdate().Tx(t.tx)
}

func (t *dbtx) Select(str ...string) *SB {
	return Select(str...).Tx(t.tx)
}

func (t *dbtx) SelectResults(query string, args ...interface{}) (Results, error) {
	rows, err := t.db.txQueryContextHandler(t.ctx, t.tx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToResults(rows)
}

func (t *dbtx) SelectOne(query string, args ...interface{}) (OneRow, error) {
	ret, err := t.SelectResults(query, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) > 0 {
		return ret[0], nil
	}
	return make(OneRow), nil
}

func (t *dbtx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.ExecContext(t.ctx, query, args...)
}

func (t *dbtx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.db.txExecContextHandler(ctx, t.tx, query, args...)
}

func (t *dbtx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.QueryContext(t.ctx, query, args...)
}

func (t *dbtx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.db.txQueryContextHandler(ctx, t.tx, query, args...)
}

func (t *dbtx) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.QueryRowContext(t.ctx, query, args...)
}

func (t *dbtx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.db.txQueryRowContextHandler(ctx, t.tx, query, args...)
}

func (t *dbtx) Prepare(query string) (*sql.Stmt, error) {
	return t.PrepareContext(t.ctx, query)
}

func (t *dbtx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return t.db.txPrepareContextHandler(ctx, t.tx, query)
}

func (t *dbtx) Stmt(stmt *sql.Stmt) *sql.Stmt {
	return t.StmtContext(t.ctx, stmt)
}

func (t *dbtx) StmtContext(ctx context.Context, stmt *sql.Stmt) *sql.Stmt {
	return t.db.txStmtContextHandler(ctx, t.tx, stmt)
}
