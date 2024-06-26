// Copyright (c) 2020. Homeland Interactive Technology Ltd. All rights reserved.

// Package db 数据库工具包
package db

import (
	"context"
	"database/sql"

	"git.yuetanggame.com/zfish/fishpkg/util"
	"golang.org/x/crypto/ssh"
)

// DefaultDB 默认实例
var DefaultDB *Database

func getDatabase(d ...*Database) *Database {
	if len(d) == 1 {
		return d[0]
	}
	return DefaultDB
}

// Begin 开启事务
func Begin(d ...*Database) (*sql.Tx, error) {
	return BeginContext(context.Background(), d...)
}

// BeginContext 开启事务
func BeginContext(ctx context.Context, d ...*Database) (*sql.Tx, error) {
	return getDatabase(d...).BeginContext(ctx)
}

// RawExec 执行语句
func RawExec(query string, args ...interface{}) (sql.Result, error) {
	return RawExecContext(context.Background(), query, args...)
}

// RawExecContext 执行语句
func RawExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return DefaultDB.ExecContext(ctx, query, args...)
}

// RawInsert 执行 INSERT 语句并返回最后生成的自增ID
// 返回0表示没有出错, 但没生成自增ID
// 返回-1表示出错
func RawInsert(query string, args ...interface{}) (int64, error) {
	return RawInsertContext(context.Background(), query, args...)
}

// RawInsertContext 执行 INSERT 语句并返回最后生成的自增ID
// 返回0表示没有出错, 但没生成自增ID
// 返回-1表示出错
func RawInsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	ret, err := DefaultDB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	last, err := ret.LastInsertId()
	if err != nil {
		return -1, err

	}
	return last, nil
}

// RawUpdate 执行 UPDATE 语句并返回受影响的行数
// 返回0表示没有出错, 但没有被更新的行
// 返回-1表示出错
func RawUpdate(query string, args ...interface{}) (int64, error) {
	return RawUpdateContext(context.Background(), query, args...)
}

// RawUpdateContext 执行 UPDATE 语句并返回受影响的行数
// 返回0表示没有出错, 但没有被更新的行
// 返回-1表示出错
func RawUpdateContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	ret, err := DefaultDB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	aff, err := ret.RowsAffected()
	if err != nil {
		return -1, err
	}
	return aff, nil
}

// RawDelete 执行 DELETE 语句并返回受影响的行数
// 返回0表示没有出错, 但没有被删除的行
// 返回-1表示出错
func RawDelete(query string, args ...interface{}) (int64, error) {
	return RawDeleteContext(context.Background(), query, args...)
}

// RawDeleteContext 执行 DELETE 语句并返回受影响的行数
// 返回0表示没有出错, 但没有被删除的行
// 返回-1表示出错
func RawDeleteContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	return DefaultDB.DeleteContext(ctx, query, args...)
}

// RawQuery 查询单条记录
func RawQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return RawQueryContext(context.Background(), query, args...)
}

// RawQueryContext 查询单条记录
func RawQueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return DefaultDB.QueryContext(ctx, query, args...)
}

// RawQueryRow 查询单条记录
func RawQueryRow(query string, args ...interface{}) *sql.Row {
	return RawQueryRowContext(context.Background(), query, args...)
}

// RawQueryRowContext 查询单条记录
func RawQueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return DefaultDB.QueryRowContext(ctx, query, args...)
}

// RawSelect 查询不定字段的结果集
func RawSelect(query string, args ...interface{}) (Results, error) {
	return RawSelectContext(context.Background(), query, args...)
}

// RawSelectContext 查询不定字段的结果集
func RawSelectContext(ctx context.Context, query string, args ...interface{}) (Results, error) {
	rows, err := DefaultDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToResults(rows)
}

// RawSelectOne 查询一行不定字段的结果
func RawSelectOne(query string, args ...interface{}) (OneRow, error) {
	return RawSelectOneContext(context.Background(), query, args...)
}

// RawSelectOneContext 查询一行不定字段的结果
func RawSelectOneContext(ctx context.Context, query string, args ...interface{}) (OneRow, error) {
	ret, err := DefaultDB.SelectContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) > 0 {
		return ret[0], nil
	}
	return make(OneRow), nil
}

// NewDatabase 创建数据库实例
func NewDatabase() *Database {
	d := &Database{
		Type: "mysql",
	}
	d.interceptors = append(d.interceptors, globalInterceptors...)
	return d.AddInterceptor(nil)
}

// Database 数据容器抽象对象定义
type Database struct {
	Type         string // 用来给 SqlBuilder 进行一些特殊的判断 (空值或 mysql 皆表示这是一个 MySQL 实例)
	DB           *sql.DB
	SSHClient    *ssh.Client
	interceptors []Interceptor

	beginContextHandler     BeginContextFunc
	execContextHandler      ExecContextFunc
	insertContextHandler    InsertContextFunc
	updateContextHandler    UpdateContextFunc
	deleteContextHandler    DeleteContextFunc
	queryContextHandler     QueryContextFunc
	queryRowContextHandler  QueryRowContextFunc
	selectContextHandler    SelectContextFunc
	selectOneContextHandler SelectOneContextFunc

	txExecContextHandler     TxExecContextFunc
	txQueryContextHandler    TxQueryContextFunc
	txQueryRowContextHandler TxQueryRowContextFunc
	txPrepareContextHandler  TxPrepareContextFunc
	txStmtContextHandler     TxStmtContextFunc
}

func (dba *Database) beginContext(ctx context.Context) (*sql.Tx, error) {
	return dba.DB.BeginTx(ctx, nil)
}

func (dba *Database) execContext(
	ctx context.Context, query string, args ...interface{}) (sql.Result, error) {

	tx := hasTx(ctx)
	if tx != nil {
		return tx.tx.ExecContext(ctx, query, args...)
	}

	return dba.DB.ExecContext(ctx, query, args...)
}

func (dba *Database) execAffectedContext(
	ctx context.Context, query string, args ...interface{}) (int64, error) {

	ret, err := dba.execContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	aff, err := ret.RowsAffected()
	if err != nil {
		return -1, err
	}
	return aff, nil
}

func (dba *Database) insertContext(
	ctx context.Context, query string, args ...interface{}) (int64, error) {

	ret, err := dba.execContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	last, err := ret.LastInsertId()
	if err != nil {
		return -1, err

	}
	return last, nil
}

func (dba *Database) updateContext(
	ctx context.Context, query string, args ...interface{}) (int64, error) {

	return dba.execAffectedContext(ctx, query, args...)
}

func (dba *Database) deleteContext(
	ctx context.Context, query string, args ...interface{}) (int64, error) {

	return dba.execAffectedContext(ctx, query, args...)
}

func (dba *Database) queryContext(
	ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {

	tx := hasTx(ctx)
	if tx != nil {
		return tx.tx.QueryContext(ctx, query, args...)
	}

	return dba.DB.QueryContext(ctx, query, args...)
}

func (dba *Database) queryRowContext(
	ctx context.Context, query string, args ...interface{}) *sql.Row {

	tx := hasTx(ctx)
	if tx != nil {
		return tx.tx.QueryRowContext(ctx, query, args...)
	}

	return dba.DB.QueryRowContext(ctx, query, args...)
}

func (dba *Database) selectContext(
	ctx context.Context, query string, args ...interface{}) (Results, error) {

	var (
		rows *sql.Rows
		err  error
	)
	tx := hasTx(ctx)
	if tx != nil {
		rows, err = tx.tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = dba.DB.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToResults(rows)
}

func (dba *Database) selectOneContext(
	ctx context.Context, query string, args ...interface{}) (OneRow, error) {

	ret, err := dba.selectContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) > 0 {
		return ret[0], nil
	}
	return make(OneRow), nil
}

func (*Database) txExecContext(
	ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {

	return tx.ExecContext(ctx, query, args...)
}

func (*Database) txQueryContext(
	ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error) {

	return tx.QueryContext(ctx, query, args...)
}

func (*Database) txQueryRowContext(
	ctx context.Context, tx *sql.Tx, query string, args ...interface{}) *sql.Row {

	return tx.QueryRowContext(ctx, query, args...)
}

func (*Database) txPrepareContext(ctx context.Context, tx *sql.Tx, query string) (*sql.Stmt, error) {
	return tx.PrepareContext(ctx, query)
}

func (*Database) txStmtContext(ctx context.Context, tx *sql.Tx, stmt *sql.Stmt) *sql.Stmt {
	return tx.StmtContext(ctx, stmt)
}

// ================================================================================================
// ================================================================================================
// ================================================================================================

// Close 关闭数据库连接
func (dba *Database) Close() error {
	defer func() {
		if dba.SSHClient != nil {
			_ = dba.SSHClient.Close()
		}
	}()
	return dba.DB.Close()
}

// AddInterceptor 添加拦截器, 越后添加的越先执行
func (dba *Database) AddInterceptor(i Interceptor) *Database {
	if i != nil {
		dba.interceptors = append(dba.interceptors, i)
	}

	dba.beginContextHandler = dba.beginContext
	dba.execContextHandler = dba.execContext
	dba.insertContextHandler = dba.insertContext
	dba.updateContextHandler = dba.updateContext
	dba.deleteContextHandler = dba.deleteContext
	dba.queryContextHandler = dba.queryContext
	dba.queryRowContextHandler = dba.queryRowContext
	dba.selectContextHandler = dba.selectContext
	dba.selectOneContextHandler = dba.selectOneContext

	dba.txExecContextHandler = dba.txExecContext
	dba.txQueryContextHandler = dba.txQueryContext
	dba.txQueryRowContextHandler = dba.txQueryRowContext
	dba.txPrepareContextHandler = dba.txPrepareContext
	dba.txStmtContextHandler = dba.txStmtContext

	for i := len(dba.interceptors) - 1; i >= 0; i-- {
		dba.beginContextHandler = dba.interceptors[i].BeginContext(dba.beginContextHandler)
		dba.execContextHandler = dba.interceptors[i].ExecContext(dba.execContextHandler)
		dba.insertContextHandler = dba.interceptors[i].InsertContext(dba.insertContextHandler)
		dba.updateContextHandler = dba.interceptors[i].UpdateContext(dba.updateContextHandler)
		dba.deleteContextHandler = dba.interceptors[i].DeleteContext(dba.deleteContextHandler)
		dba.queryContextHandler = dba.interceptors[i].QueryContext(dba.queryContextHandler)
		dba.queryRowContextHandler = dba.interceptors[i].QueryRowContext(dba.queryRowContextHandler)
		dba.selectContextHandler = dba.interceptors[i].SelectContext(dba.selectContextHandler)
		dba.selectOneContextHandler = dba.interceptors[i].SelectOneContext(dba.selectOneContextHandler)

		dba.txExecContextHandler = dba.interceptors[i].TxExecContext(dba.txExecContextHandler)
		dba.txQueryContextHandler = dba.interceptors[i].TxQueryContext(dba.txQueryContextHandler)
		dba.txQueryRowContextHandler = dba.interceptors[i].TxQueryRowContext(dba.txQueryRowContextHandler)
		dba.txPrepareContextHandler = dba.interceptors[i].TxPrepareContext(dba.txPrepareContextHandler)
		dba.txStmtContextHandler = dba.interceptors[i].TxStmtContext(dba.txStmtContextHandler)
	}
	return dba
}

// Begin 开启事务
func (dba *Database) Begin() (*sql.Tx, error) {
	return dba.beginContextHandler(context.Background())
}

// BeginContext 开启事务
func (dba *Database) BeginContext(ctx context.Context) (*sql.Tx, error) {
	return dba.beginContextHandler(ctx)
}

// Exec 执行语句
func (dba *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dba.execContextHandler(context.Background(), query, args...)
}

// ExecContext 执行语句
func (dba *Database) ExecContext(
	ctx context.Context, query string, args ...interface{}) (sql.Result, error) {

	return dba.execContextHandler(ctx, query, args...)
}

// Insert 执行 INSERT 语句并返回最后生成的自增ID
// 返回0表示没有出错, 但没生成自增ID
// 返回-1表示出错
func (dba *Database) Insert(query string, args ...interface{}) (int64, error) {
	return dba.insertContextHandler(context.Background(), query, args...)
}

// InsertContext 执行 INSERT 语句并返回最后生成的自增ID
// 返回0表示没有出错, 但没生成自增ID
// 返回-1表示出错
func (dba *Database) InsertContext(
	ctx context.Context, query string, args ...interface{}) (int64, error) {

	return dba.insertContextHandler(ctx, query, args...)
}

// Update 执行 UPDATE 语句并返回受影响的行数
// 返回0表示没有出错, 但没有被更新的行
// 返回-1表示出错
func (dba *Database) Update(query string, args ...interface{}) (int64, error) {
	return dba.updateContextHandler(context.Background(), query, args...)
}

// UpdateContext 执行 UPDATE 语句并返回受影响的行数
// 返回0表示没有出错, 但没有被更新的行
// 返回-1表示出错
func (dba *Database) UpdateContext(
	ctx context.Context, query string, args ...interface{}) (int64, error) {

	return dba.updateContextHandler(ctx, query, args...)
}

// Delete 执行 DELETE 语句并返回受影响的行数
// 返回0表示没有出错, 但没有被删除的行
// 返回-1表示出错
func (dba *Database) Delete(query string, args ...interface{}) (int64, error) {
	return dba.deleteContextHandler(context.Background(), query, args...)
}

// DeleteContext 执行 DELETE 语句并返回受影响的行数
// 返回0表示没有出错, 但没有被删除的行
// 返回-1表示出错
func (dba *Database) DeleteContext(
	ctx context.Context, query string, args ...interface{}) (int64, error) {

	return dba.deleteContextHandler(ctx, query, args...)
}

// Query 查询记录
func (dba *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return dba.queryContextHandler(context.Background(), query, args...)
}

// QueryContext 查询记录
func (dba *Database) QueryContext(
	ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {

	return dba.queryContextHandler(ctx, query, args...)
}

// QueryRow 查询单条记录
func (dba *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return dba.queryRowContextHandler(context.Background(), query, args...)
}

// QueryRowContext 查询单条记录
func (dba *Database) QueryRowContext(
	ctx context.Context, query string, args ...interface{}) *sql.Row {

	return dba.queryRowContextHandler(ctx, query, args...)
}

// Select 查询不定字段的结果集
func (dba *Database) Select(query string, args ...interface{}) (Results, error) {
	return dba.selectContextHandler(context.Background(), query, args...)
}

// SelectContext 查询不定字段的结果集
func (dba *Database) SelectContext(
	ctx context.Context, query string, args ...interface{}) (Results, error) {

	return dba.selectContextHandler(ctx, query, args...)
}

// SelectOne 查询一行不定字段的结果
func (dba *Database) SelectOne(query string, args ...interface{}) (OneRow, error) {
	return dba.selectOneContextHandler(context.Background(), query, args...)
}

// SelectOneContext 查询一行不定字段的结果
func (dba *Database) SelectOneContext(
	ctx context.Context, query string, args ...interface{}) (OneRow, error) {

	return dba.selectOneContextHandler(ctx, query, args...)
}

// TxExecContext 执行带事务带语句
func (dba *Database) TxExecContext(
	ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {

	return dba.txExecContextHandler(ctx, tx, query, args...)
}

// TxQueryContext 通过事务查询记录
func (dba *Database) TxQueryContext(
	ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error) {

	return dba.txQueryContextHandler(ctx, tx, query, args...)
}

// TxQueryRowContext 通过事务查询单条记录
func (dba *Database) TxQueryRowContext(
	ctx context.Context, tx *sql.Tx, query string, args ...interface{}) *sql.Row {

	return dba.txQueryRowContextHandler(ctx, tx, query, args...)
}

func rowsToResults(rows *sql.Rows) (Results, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	colNum := len(cols)
	rawValues := make([][]byte, colNum)

	// query.Scan 的参数，因为每次查询出来的列是不定长的，所以传入长度固定当次查询的长度
	scans := make([]interface{}, len(cols))

	// 将每行数据填充到[][]byte里
	for i := range rawValues {
		scans[i] = &rawValues[i]
	}

	results := make(Results, 0)
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]string)

		for k, raw := range rawValues {
			key := cols[k]
			/* if raw == nil {
				row[key] = "\\N"
			} else { */
			row[key] = string(raw)
			// }
		}
		results = append(results, row)
	}
	return results, nil
}

// Results 多行数据集结果
type Results []OneRow

// OneRow 单行查询结果
type OneRow map[string]string

// Set 设置值
func (row OneRow) Set(key, val string) {
	row[key] = val
}

// Exist 判断字段是否存在
func (row OneRow) Exist(field string) bool {
	if _, ok := row[field]; ok {
		return true
	}
	return false
}

// Get 获取指定字段的值
func (row OneRow) Get(field string) string {
	if v, ok := row[field]; ok {
		return v
	}
	return ""
}

// GetInt8 获取指定字段的 int8 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetInt8(field string) int8 {
	if v, ok := row[field]; ok {
		return util.Atoi8(v)
	}
	return 0
}

// GetInt16 获取指定字段的 int16 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetInt16(field string) int16 {
	if v, ok := row[field]; ok {
		return util.Atoi16(v)
	}
	return 0
}

// GetInt 获取指定字段的 int 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetInt(field string) int {
	if v, ok := row[field]; ok {
		return util.Atoi(v)
	}
	return 0
}

// GetInt32 获取指定字段的 int32 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetInt32(field string) int32 {
	if v, ok := row[field]; ok {
		return util.Atoi32(v)
	}
	return 0
}

// GetInt64 获取指定字段的 int64 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetInt64(field string) int64 {
	if v, ok := row[field]; ok {
		return util.Atoi64(v)
	}
	return 0
}

// GetUint8 获取指定字段的 uint8 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetUint8(field string) uint8 {
	if v, ok := row[field]; ok {
		return util.Atoui8(v)
	}
	return 0
}

// GetUint16 获取指定字段的 uint16 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetUint16(field string) uint16 {
	if v, ok := row[field]; ok {
		return util.Atoui16(v)
	}
	return 0
}

// GetUint 获取指定字段的 uint 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetUint(field string) uint {
	if v, ok := row[field]; ok {
		return util.Atoui(v)
	}
	return 0
}

// GetUint32 获取指定字段的 uint32 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetUint32(field string) uint32 {
	if v, ok := row[field]; ok {
		return util.Atoui32(v)
	}
	return 0
}

// GetUint64 获取指定字段的 uint64 类型值, 注意, 如果该字段不存在则会返回0
func (row OneRow) GetUint64(field string) uint64 {
	if v, ok := row[field]; ok {
		return util.Atoui64(v)
	}
	return 0
}
