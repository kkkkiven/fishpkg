// Copyright (c) 2020. Homeland Interactive Technology Ltd. All rights reserved.

package db

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"sync"

	"git.yuetanggame.com/zfish/fishpkg/util"
)

const (
	_ = iota
	TypeInsert
	TypeDelete
	TypeUpdate
	TypeSelect
	TypeInsertUpdate
)

const WrapSymbol = "`"

var (
	dbType = "mysql"

	sbPool = &sync.Pool{
		New: func() interface{} {
			return &SB{
				buf:  GetBufferPool(),
				args: make([]interface{}, 0),
			}
		},
	}

	ErrTableEmpty = errors.New("table cannot be empty")
	ErrValueEmpty = errors.New("values cannot be empty")
	ErrDeleteAll  = errors.New("delete all data is not safe")
	ErrUpdateAll  = errors.New("update all data is not safe")
)

// SetDBType 数据连接的数据库类型
func SetDBType(typ string) {
	dbType = typ
}

// RawVal 字面值
type RawVal string

// IncVal 增量值, 这个类型的作用是对整型字段进行自增(减), 其需要置于 SBValues 中
type IncVal struct {
	Val       int64  // 增量
	BaseField string // 为空表示对当前字段累加
}

// Insert 构建INSERT语句
func Insert(ignore ...bool) *SB {
	var i bool
	if len(ignore) == 1 && ignore[0] {
		i = true
	}
	sb := sbPool.Get().(*SB)
	sb.typ = TypeInsert
	sb.db = DefaultDB
	sb.ignore = i
	sb.values = GetValues()
	return sb
}

// Delete 构建DELETE语句
func Delete() *SB {
	sb := sbPool.Get().(*SB)
	sb.typ = TypeDelete
	sb.db = DefaultDB
	return sb
}

// Update 构建UPDATE语句
func Update() *SB {
	sb := sbPool.Get().(*SB)
	sb.typ = TypeUpdate
	sb.db = DefaultDB
	sb.values = GetValues()
	return sb
}

// InsertUpdate 构建InsertUpdate语句, 仅针对MySQL有效, 内部使用ON DUPLICATE KEY UPDATE方式实现
func InsertUpdate() *SB {
	sb := sbPool.Get().(*SB)
	sb.typ = TypeInsertUpdate
	sb.db = DefaultDB
	sb.values = GetValues()
	sb.duplicateUpdateValues = GetValues()
	return sb
}

// Select 构建SELECT语句
func Select(str ...string) *SB {
	fields := "*"
	if len(str) == 1 {
		fields = str[0]
	}
	sb := sbPool.Get().(*SB)
	sb.typ = TypeSelect
	sb.db = DefaultDB
	sb.field = fields
	return sb
}

// SB SQL语句构造结构
type SB struct {
	ctx                                      context.Context
	db                                       *Database
	tx                                       *sql.Tx
	buf                                      *bytes.Buffer
	field, table, where, group, order, limit string
	values                                   *Values // INSERT、UPDATE 语句的值数据
	duplicateUpdateValues                    *Values // INSERT ... ON DUPLICATE KEY UPDATE 语句中, 遇重需更新的值数据
	args                                     []interface{}
	typ                                      int8
	retain                                   bool // 执行后不放入对象池
	ignore                                   bool
	fullSQL                                  bool
	unsafe                                   bool // 是否进行安全检查, 专门针对无限定的UPDATE和DELETE进行二次验证
}

func (q *SB) Release() {
	q.ctx = nil
	q.db = nil
	q.tx = nil
	q.buf.Reset()
	q.field = ""
	q.table = ""
	q.where = ""
	q.group = ""
	q.order = ""
	q.limit = ""
	if q.duplicateUpdateValues != nil {
		if q.duplicateUpdateValues != q.values {
			PutValues(q.duplicateUpdateValues)
		}
		q.duplicateUpdateValues = nil
	}
	if q.values != nil {
		PutValues(q.values)
		q.values = nil
	}
	q.args = q.args[:0]
	q.typ = 0
	q.retain = false
	q.ignore = false
	q.fullSQL = false
	q.unsafe = false
	sbPool.Put(q)
}

func (q *SB) release() {
	if !q.retain {
		q.Release()
	}
}

func (q *SB) Retain(v ...bool) *SB {
	q.retain = len(v) == 0 || v[0]
	return q
}

// Ctx 设置上下文
func (q *SB) Ctx(ctx context.Context) *SB {
	q.ctx = ctx
	return q
}

func (q *SB) sqlSelect() error {
	q.buf.WriteString("SELECT ")
	q.buf.WriteString(q.field)
	if q.table != "" {
		q.buf.WriteString(" FROM ")
		q.buf.WriteString(q.table)
	}
	if q.where != "" {
		q.buf.WriteString(" WHERE ")
		q.buf.WriteString(q.where)
	}
	if q.group != "" {
		q.buf.WriteString(" GROUP BY ")
		q.buf.WriteString(q.group)
	}
	if q.order != "" {
		q.buf.WriteString(" ORDER BY ")
		q.buf.WriteString(q.order)
	}
	if q.limit != "" && (q.db.Type == "" || q.db.Type == "mysql") {
		q.buf.WriteString(" LIMIT ")
		q.buf.WriteString(q.limit)
	}
	return nil
}

func (q *SB) sqlInsert() error {
	if q.values.Len() == 0 {
		return ErrValueEmpty
	}
	if q.ignore {
		q.buf.WriteString("INSERT IGNORE INTO ")
	} else {
		q.buf.WriteString("INSERT INTO ")
	}
	q.buf.WriteString(q.table)
	fields, placeholder := q.processInsertValues()
	q.buf.WriteString(" (")
	q.buf.WriteString(fields.String())
	q.buf.WriteString(") VALUES (")
	q.buf.WriteString(placeholder.String())
	q.buf.WriteString(")")

	PutBufferPool(fields)
	PutBufferPool(placeholder)
	return nil
}

func (q *SB) sqlInsertUpdate() error {
	q.buf.WriteString("INSERT INTO ")
	q.buf.WriteString(q.table)
	fields, placeholder := q.processInsertValues()
	q.buf.WriteString(" (")
	q.buf.WriteString(fields.String())
	q.buf.WriteString(") VALUES (")
	q.buf.WriteString(placeholder.String())

	PutBufferPool(fields)
	PutBufferPool(placeholder)

	q.buf.WriteString(") ON DUPLICATE KEY UPDATE ")
	placeholder = q.processUpdateValues(q.duplicateUpdateValues)
	q.buf.WriteString(placeholder.String())

	PutBufferPool(placeholder)
	return nil
}

func (q *SB) sqlUpdate() error {
	q.buf.WriteString("UPDATE ")
	q.buf.WriteString(q.table)
	q.buf.WriteString(" SET ")
	placeholder := q.processUpdateValues(q.values)
	q.buf.WriteString(placeholder.String())

	PutBufferPool(placeholder)

	if q.where != "" {
		q.buf.WriteString(" WHERE ")
		q.buf.WriteString(q.where)
	} else {
		if !q.unsafe {
			return ErrUpdateAll
		}
	}
	return nil
}

func (q *SB) sqlDelete() error {
	q.buf.WriteString("DELETE FROM ")
	q.buf.WriteString(q.table)
	if q.where != "" {
		q.buf.WriteString(" WHERE ")
		q.buf.WriteString(q.where)
	} else {
		if !q.unsafe {
			return ErrDeleteAll
		}
	}
	return nil
}

// ToSQL 构建SQL语句
// param: returnFullSQL 是否返回完整的sql语句(即:绑定参数之后的语句)
func (q *SB) ToSQL(returnFullSQL ...bool) (str string, err error) {
	if q.table == "" {
		err = ErrTableEmpty
		return
	}

	if q.buf.Len() == 0 {
		switch q.typ {
		case TypeSelect:
			err = q.sqlSelect()
		case TypeInsert:
			err = q.sqlInsert()
		case TypeInsertUpdate:
			err = q.sqlInsertUpdate()
		case TypeUpdate:
			err = q.sqlUpdate()
		case TypeDelete:
			err = q.sqlDelete()
		}
	}

	str = q.buf.String()
	if len(returnFullSQL) == 1 && returnFullSQL[0] {
		str, err = FullSQL(str, q.args...)
	}
	return str, err
}

// processInsertValues 构建插入数据占位符
func (q *SB) processInsertValues() (*bytes.Buffer, *bytes.Buffer) {
	fields, placeholder := GetBufferPool(), GetBufferPool()
	first := true
	for i, k := range q.values.Keys {
		if first {
			first = false
		} else {
			fields.WriteString(",")
			placeholder.WriteString(",")
		}
		fields.WriteString(WrapSymbol)
		fields.WriteString(k)
		fields.WriteString(WrapSymbol)
		placeholder.WriteString("?")

		q.args = append(q.args, q.values.Vals[i])
	}
	return fields, placeholder
}

// processUpdateValues 构造更新数据占位符
func (q *SB) processUpdateValues(vals *Values) *bytes.Buffer {
	placeholder := GetBufferPool()
	first := true
	for i, k := range vals.Keys {
		if first {
			first = false
		} else {
			placeholder.WriteString(",")
		}
		placeholder.WriteString(WrapSymbol)
		placeholder.WriteString(k)
		placeholder.WriteString(WrapSymbol)
		placeholder.WriteString("=")

		switch iv := vals.Vals[i].(type) {
		case IncVal:
			if iv.BaseField == "" {
				placeholder.WriteString(k)
			} else {
				placeholder.WriteString(iv.BaseField)
			}
			if iv.Val >= 0 { // 正数需要一个加号
				placeholder.WriteString("+")
			}
			placeholder.WriteString(util.I64toa(iv.Val))
		case RawVal:
			placeholder.WriteString(string(iv))
		default:
			placeholder.WriteString("?")
			q.args = append(q.args, vals.Vals[i])
		}
	}
	return placeholder
}

// DB 设置数据库对象
func (q *SB) DB(db *Database) *SB {
	q.db = db
	return q
}

// Tx 设置事务对象
func (q *SB) Tx(tx *sql.Tx) *SB {
	q.tx = tx
	return q
}

// From 设置FROM字句
func (q *SB) From(str string) *SB {
	q.table = str
	return q
}

// Table 设置表名
func (q *SB) Table(str string) *SB {
	return q.From(str)
}

// Where 设置WHERE字句
func (q *SB) Where(str string) *SB {
	q.where = str
	return q
}

// Group 设置GROUP字句
func (q *SB) Group(str string) *SB {
	q.group = str
	return q
}

// Order 设置GROUP字句
func (q *SB) Order(str string) *SB {
	q.order = str
	return q
}

// Limit 设置LIMIT字句
func (q *SB) Limit(count int64, offset ...int64) *SB {
	if len(offset) > 0 {
		q.limit = util.I64toa(offset[0]) + "," + util.I64toa(count)
	} else {
		q.limit = "0," + util.I64toa(count)
	}
	return q
}

// Unsafe 设置安全检查开关
func (q *SB) Unsafe(unsefe ...bool) *SB {
	if len(unsefe) == 1 && !unsefe[0] {
		q.unsafe = false
	} else {
		q.unsafe = true
	}
	return q
}

// Value 设置值数据
func (q *SB) Value(vals *Values) *SB {
	q.values = vals
	return q
}

// SBValue 设置值数据, 不推荐使用
func (q *SB) SBValue(vals SBValues) *SB {
	q.values = GetValues().AddSBValues(vals)
	return q
}

// DuplicateUpdateValue 设置遇重需更新的值数据
func (q *SB) DuplicateUpdateValue(vals *Values) *SB {
	q.duplicateUpdateValues = vals
	return q
}

// DuplicateUpdateSBValue 设置遇重需更新的值数据, 不推荐使用
func (q *SB) DuplicateUpdateSBValue(vals SBValues) *SB {
	q.duplicateUpdateValues = GetValues().AddSBValues(vals)
	return q
}

// AddValue 添加值
func (q *SB) AddValue(key string, val interface{}) *SB {
	q.values.Add(key, val)
	return q
}

// AddValues 添加多个值
func (q *SB) AddValues(key string, val ...interface{}) *SB {
	q.values.Adds(key, val...)
	return q
}

// SetValues 添加或更新值
func (q *SB) SetValues(key string, val interface{}) *SB {
	q.values.Set(key, val)
	return q
}

// DelValue 删除值
func (q *SB) DelValue(key string) *SB {
	q.values.Del(key)
	return q
}

// GetValues 返回当前 Builder 的 *Values
func (q *SB) GetValues() *Values {
	return q.values
}

// AddDuplicateUpdateValue 遇重需更新的值数据
func (q *SB) AddDuplicateUpdateValue(key string, val interface{}) *SB {
	q.duplicateUpdateValues.Add(key, val)
	return q
}

// AddDuplicateUpdateValues 添加多个值
func (q *SB) AddDuplicateUpdateValues(key string, val ...interface{}) *SB {
	q.duplicateUpdateValues.Adds(key, val...)
	return q
}

// SetDuplicateUpdateValues 添加或更新值
func (q *SB) SetDuplicateUpdateValues(key string, val interface{}) *SB {
	q.duplicateUpdateValues.Set(key, val)
	return q
}

// DelDuplicateUpdateValue 删除值
func (q *SB) DelDuplicateUpdateValue(key string) *SB {
	q.duplicateUpdateValues.Del(key)
	return q
}

// GetDuplicateUpdateValues 返回当前 Builder 的 *DuplicateUpdateValues
func (q *SB) GetDuplicateUpdateValues() *Values {
	return q.duplicateUpdateValues
}

// GetArgs 获取构造SQL后的参数
func (q *SB) GetArgs() []interface{} {
	return q.args
}

// FullSQL 是否直接将参数直接填充至语句中并生成SQL字符串
func (q *SB) FullSQL(yes ...bool) *SB {
	if len(yes) == 1 {
		q.fullSQL = yes[0]
	} else {
		q.fullSQL = true
	}
	return q
}

// SBResult Exec 方法的返回结果
type SBResult struct {
	Sql      string // 最后执行的SQL
	Err      error  // 错误提示信息
	LastID   int64  // 最后产生的ID
	Affected int64  // 受影响的行数
	Code     int    // 错误代码
	Success  bool   // 语句是否执行成功
}

// Exec 执行INSERT、DELETE、UPDATE语句
func (q *SB) Exec(args ...interface{}) *SBResult {
	defer q.release()
	sbRet := &SBResult{}
	sbRet.Sql, sbRet.Err = q.ToSQL()
	if sbRet.Err != nil {
		return sbRet
	}

	if q.ctx == nil {
		q.ctx = context.Background()
	}

	var ret sql.Result
	if q.fullSQL {
		var sqlStr string
		sqlStr, sbRet.Err = FullSQL(sbRet.Sql, append(q.args, args...)...)
		if sbRet.Err != nil {
			return sbRet
		}
		if q.tx != nil {
			ret, sbRet.Err = q.db.TxExecContext(q.ctx, q.tx, sqlStr)
		} else {
			ret, sbRet.Err = q.db.ExecContext(q.ctx, sqlStr)
		}
	} else {
		if q.tx != nil {
			ret, sbRet.Err = q.db.TxExecContext(q.ctx, q.tx, sbRet.Sql, append(q.args, args...)...)
		} else {
			ret, sbRet.Err = q.db.ExecContext(q.ctx, sbRet.Sql, append(q.args, args...)...)
		}
		if sbRet.Err != nil {
			return sbRet
		}
	}

	sbRet.Success = true
	switch q.typ {
	case TypeInsert:
		if dbType == "mysql" {
			last, err := ret.LastInsertId()
			if err == nil {
				sbRet.LastID = last
			}
		}
	case TypeDelete, TypeUpdate, TypeInsertUpdate:
		aff, err := ret.RowsAffected()
		if err == nil {
			sbRet.Affected = aff
		}
	}
	return sbRet
}

// Query 查询记录集
func (q *SB) Query(args ...interface{}) (Results, error) {
	defer q.release()
	s, err := q.ToSQL()
	if err != nil {
		return nil, err
	}

	if q.ctx == nil {
		q.ctx = context.Background()
	}

	var rows *sql.Rows
	if q.tx == nil {
		rows, err = q.db.QueryContext(q.ctx, s, args...)
	} else {
		rows, err = q.db.TxQueryContext(q.ctx, q.tx, s, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToResults(rows)
}

// QueryOne 查询单行数据
func (q *SB) QueryOne(args ...interface{}) (OneRow, error) {
	defer q.release()
	q.Limit(1, 0)
	s, err := q.ToSQL()
	if err != nil {
		return nil, err
	}

	if q.ctx == nil {
		q.ctx = context.Background()
	}

	var rows *sql.Rows
	if q.tx == nil {
		rows, err = q.db.QueryContext(q.ctx, s, args...)
	} else {
		rows, err = q.db.TxQueryContext(q.ctx, q.tx, s, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret, err := rowsToResults(rows)
	if err != nil {
		return nil, err
	}
	if len(ret) > 0 {
		return ret[0], nil
	}
	return OneRow{}, nil
}

// QueryAllRow 查询记录集
func (q *SB) QueryAllRow(args ...interface{}) (*sql.Rows, error) {
	defer q.release()
	s, e := q.ToSQL()
	if e != nil {
		return nil, e
	}

	if q.ctx == nil {
		q.ctx = context.Background()
	}

	if q.tx == nil {
		return q.db.QueryContext(q.ctx, s, args...)
	}
	return q.db.TxQueryContext(q.ctx, q.tx, s, args...)
}

// QueryRow 查询单行数据
func (q *SB) QueryRow(args ...interface{}) *sql.Row {
	defer q.release()
	q.Limit(1, 0)
	s, e := q.ToSQL()
	if e != nil {
		return nil
	}

	if q.ctx == nil {
		q.ctx = context.Background()
	}

	if q.tx == nil {
		return q.db.QueryRowContext(q.ctx, s, args...)
	}
	return q.db.TxQueryRowContext(q.ctx, q.tx, s, args...)
}

// Exists 判断记录是否存在
func (q *SB) Exists(args ...interface{}) (bool, error) {
	var i int8
	q.field = "1"
	err := q.QueryRow(args...).Scan(&i)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return i == 1, nil
	}
}
