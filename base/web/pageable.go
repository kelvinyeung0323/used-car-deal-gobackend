package web

import (
	"fmt"
	"reflect"
	"strings"
)

type Page[T any] struct {
	PageSize int64
	PageNum  int64
	Total    int64
	Data     []T
}

type IPageable interface {
	ISortable
	GetPageSize() int64
	Offset() int64
}
type ISortable interface {
	SortOrderStmt() string
}

type Pageable[T any] struct {
	Sortable[T]
	PageNum  int64 `form:"pageNum"`
	PageSize int64 `form:"pageSize"`
}

type Sortable[T any] struct {
	Sorts     string `form:"sorts"`
	sortClass *T
}

// Offset 返回偏移量，即返回从第几行开始
func (p *Pageable[any]) Offset() int64 {
	return (p.PageNum - 1) * p.PageSize
}
func (p *Pageable[any]) GetPageSize() int64 {
	return p.PageSize
}

// SortOrderStmt 生成order by statement
func (p *Sortable[any]) SortOrderStmt() string {
	var stmt string
	objType := reflect.TypeOf(p.sortClass)

	sorts := strings.Split(p.Sorts, ",")
	sortLen := len(sorts)
	for i, s := range sorts {
		s1 := strings.Split(s, ".")
		col := s1[0]
		if field, ok := structFieldByName(objType, col); ok {
			colName := field.Tag.Get("db")
			stmt += colName
			if len(s1) > 1 {
				order := strings.ToLower(s1[1])
				if order == "desc" || order == "asc" {
					stmt += " " + order
				} else {
					panic(fmt.Errorf("排序标识不正确"))
				}
			}
			if i+1 < sortLen {
				stmt += ","
			}
		}
	}
	if stmt != "" {
		stmt = "order by " + stmt + " "
	}
	return stmt
}
func structFieldByName(p reflect.Type, searchName string) (*reflect.StructField, bool) {
	log.Debugf("type name:%v", p.Name())
	el := p.Elem()
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		if strings.ToLower(f.Name) == strings.ToLower(searchName) {
			return &f, true
		}
	}
	return nil, false
}
