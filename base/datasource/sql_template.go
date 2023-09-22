package datasource

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
	"text/template"
	"used-car-deal-gobackend/base/logger"
	"used-car-deal-gobackend/base/types"
)

const sqlFileDir = "./conf/sql/"

type SqlTemplate struct {
	tmpl     *template.Template
	fileName string
	txMgr    *TransactionManger
}

var log = logger.GetInstance()

func NewSqlTemplate(txMgr *TransactionManger, fileName string) (*SqlTemplate, error) {
	tmpl, err := template.New("sqlTemplate").Funcs(
		template.FuncMap{
			"val":    valueFunc,
			"rawVal": rawValueFunc,
			"join":   joinFunc,
		}).ParseFiles(sqlFileDir + fileName)
	if err != nil {
		return nil, fmt.Errorf("解释SQL模板文件错误！%v", err)
	}
	return &SqlTemplate{tmpl: tmpl, fileName: fileName, txMgr: txMgr}, nil
}

//sqlString 根据模板名称获取模板，并根据输入的数据执行模板
func (s *SqlTemplate) sqlString(tName string, param any) (string, error) {
	tmpl := s.tmpl.Lookup(tName)
	if tmpl == nil {
		return "", fmt.Errorf("找不到对应的sql模板:%v", tName)
	}
	sb := &strings.Builder{}
	err := tmpl.Execute(sb, param)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

//Get
//返回两个值，第一个值是整数，0表示数据库没有对应的数据，1 表示有数据，-1 表示存在错误
func (s *SqlTemplate) Get(ctx *gin.Context, dest any, tName string, param any) (int, error) {
	conn, err := s.txMgr.GetConn(ctx)
	if err != nil {
		return -1, err
	}
	sqlStr, err := s.sqlString(tName, param)
	if err != nil {
		return -1, err
	}
	log.Debug("SQL is: %s\n", sqlStr)
	err = conn.Get(dest, sqlStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return -1, err
	}
	return 1, nil
}

func (s *SqlTemplate) Select(ctx *gin.Context, dest any, tmpl string, param any) error {
	conn, err := s.txMgr.GetConn(ctx)
	if err != nil {
		return err
	}
	sql, err := s.sqlString(tmpl, param)
	if err != nil {
		return err
	}
	log.Debug("SQL is: %s\n", sql)
	err = conn.Select(dest, sql)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlTemplate) Exec(ctx *gin.Context, tmpl string, param any) (sql.Result, error) {
	conn, err := s.txMgr.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	sql, err := s.sqlString(tmpl, param)
	log.Debug("SQL is: %s\n", sql)
	r, err := conn.Exec(sql)
	return r, err
}

//valueFunc 解释模板函数
//用于把值转换成SQL语句中的字段值
var valueFunc = func(arg any) (any, error) {

	if arg == nil || (reflect.ValueOf(arg).CanAddr() && reflect.ValueOf(arg).IsNil()) {
		return "null", nil
	}
	switch arg := arg.(type) {
	case *string:
		return "'" + strings.ReplaceAll(*arg, "'", "\\'") + "'", nil
	case string:
		return "'" + strings.ReplaceAll(arg, "'", "\\'") + "'", nil
	case types.NullMarshalInterface:
		v, err := arg.MarshalJSON()
		return string(v[:]), err
	default:
		return arg, nil
	}

}

//rawValueFunc 用于处理SQL中的in语句的拼接
//参数必须是interface类型，基本类型会抛出错误
var rawValueFunc = func(arg any) (any, error) {

	if arg == nil || reflect.ValueOf(arg).IsNil() {
		return "null", nil
	}
	switch arg := arg.(type) {
	case *string:
		return strings.ReplaceAll(*arg, "'", "\\'"), nil
	case string:
		return strings.ReplaceAll(arg, "'", "\\'"), nil
	case types.NullMarshalInterface:
		v, err := arg.MarshalJSON()
		return string(v[:]), err
	default:
		return arg, nil
	}

}

//参数必须是interface类型，基本类型会抛出错误
var joinFunc = func(arg any) (any, error) {
	val := reflect.ValueOf(arg)
	if arg == nil || val.IsNil() {
		return "", nil
	}

	if val.Kind() != reflect.Slice {
		panic("参数不是slice类型")
	}
	len := val.Len()
	if len <= 0 {
		return "", nil
	}
	sb := &strings.Builder{}
	for i := 0; i < len; i++ {
		sb.WriteString(",")
		v := val.Index(i)
		// s,_:=rawValueFunc(v)
		sb.WriteString(fmt.Sprint(v))
	}
	s := sb.String()
	return s[1:], nil

}
