package where_builder

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type Expr interface {
	ToWhere() (query interface{}, args []interface{})
}

type Eq map[string]interface{}

func (b Eq) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	sortedKeys := getSortedKeys(b)
	for _, key := range sortedKeys {
		val := b[key]
		exprs = append(exprs, fmt.Sprintf("%s = ?", key))
		args = append(args, val)
	}

	query = strings.Join(exprs, " AND ")
	return
}

type Neq map[string]interface{}

func (b Neq) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	sortedKeys := getSortedKeys(b)
	for _, key := range sortedKeys {
		val := b[key]
		exprs = append(exprs, fmt.Sprintf("%s != ?", key))
		args = append(args, val)
	}

	query = strings.Join(exprs, " AND ")
	return
}

type GtOrEq map[string]interface{}

func (b GtOrEq) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	sortedKeys := getSortedKeys(b)
	for _, key := range sortedKeys {
		val := b[key]
		exprs = append(exprs, fmt.Sprintf("%s >= ?", key))
		args = append(args, val)
	}

	query = strings.Join(exprs, " AND ")
	return
}

type Gt map[string]interface{}

func (b Gt) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	sortedKeys := getSortedKeys(b)
	for _, key := range sortedKeys {
		val := b[key]
		exprs = append(exprs, fmt.Sprintf("%s > ?", key))
		args = append(args, val)
	}

	query = strings.Join(exprs, " AND ")
	return
}

type Lt map[string]interface{}

func (b Lt) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	sortedKeys := getSortedKeys(b)
	for _, key := range sortedKeys {
		val := b[key]
		exprs = append(exprs, fmt.Sprintf("%s < ?", key))
		args = append(args, val)
	}

	query = strings.Join(exprs, " AND ")
	return
}

type LtOrEq map[string]interface{}

func (b LtOrEq) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	sortedKeys := getSortedKeys(b)
	for _, key := range sortedKeys {
		val := b[key]
		exprs = append(exprs, fmt.Sprintf("%s <= ?", key))
		args = append(args, val)
	}

	query = strings.Join(exprs, " AND ")
	return
}

type In map[string]interface{}

func (b In) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	sortedKeys := getSortedKeys(b)
	for _, key := range sortedKeys {
		val := b[key]
		exprs = append(exprs, fmt.Sprintf("%s IN (?)", key))
		args = append(args, val)
	}

	query = strings.Join(exprs, " AND ")
	return
}

func getSortedKeys(exp map[string]interface{}) []string {
	sortedKeys := make([]string, 0, len(exp))
	for k := range exp {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

type Like map[string]interface{}

func (b Like) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	sortedKeys := getSortedKeys(b)
	for _, key := range sortedKeys {
		val := b[key]
		exprs = append(exprs, fmt.Sprintf("%s LIKE ?", key))
		args = append(args, val)
	}

	query = strings.Join(exprs, " AND ")
	return
}

type Or []Expr

func (b Or) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	for _, expr := range b {
		expr, arg := expr.ToWhere()
		exprs = append(exprs, fmt.Sprintf("%s", expr))
		args = append(args, arg...)
	}

	query = "(" + strings.Join(exprs, " OR ") + ")"
	return
}

type And []Expr

func (b And) ToWhere() (query interface{}, args []interface{}) {
	var exprs []string

	for _, expr := range b {
		expr, arg := expr.ToWhere()
		exprs = append(exprs, fmt.Sprintf("%s", expr))
		args = append(args, arg...)
	}

	query = "(" + strings.Join(exprs, " AND ") + ")"
	return
}

func ToWhere(exprs []Expr) (string, []interface{}) {
	if len(exprs) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	var args []interface{}

	for _, e := range exprs {
		expr, arg := e.ToWhere()
		buf.WriteString(expr.(string))
		buf.WriteString(" AND ")
		args = append(args, arg...)
	}

	query := buf.String()
	return query[:len(query)-5], args
}
