package where_builder

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestEq(t *testing.T) {
	eq := Eq{"cate": "123"}
	query, args := eq.ToWhere()
	assert.Equal(t, "cate = ?", query)
	assert.Equal(t, []interface{}{"123"}, args)

	eq = Eq{
		"cate": "123",
		"id":   "456",
	}
	query, args = eq.ToWhere()
	assert.Equal(t, "cate = ? AND id = ?", query)
	assert.Equal(t, []interface{}{"123", "456"}, args)
}

func TestGtOrEq(t *testing.T) {
	eq := GtOrEq{"cate": "123"}
	query, args := eq.ToWhere()
	assert.Equal(t, "cate >= ?", query)
	assert.Equal(t, []interface{}{"123"}, args)

	eq = GtOrEq{
		"cate": "123",
		"id":   "456",
	}
	query, args = eq.ToWhere()
	assert.Equal(t, "cate >= ? AND id >= ?", query)
	assert.Equal(t, []interface{}{"123", "456"}, args)
}

func TestLtOrEq(t *testing.T) {
	eq := LtOrEq{"cate": "123"}
	query, args := eq.ToWhere()
	assert.Equal(t, "cate <= ?", query)
	assert.Equal(t, []interface{}{"123"}, args)

	eq = LtOrEq{
		"cate": "123",
		"id":   "456",
	}
	query, args = eq.ToWhere()
	assert.Equal(t, "cate <= ? AND id <= ?", query)
	assert.Equal(t, []interface{}{"123", "456"}, args)
}

func TestIn(t *testing.T) {
	eq := In{"cate": []string{"123", "456"}}
	query, args := eq.ToWhere()
	assert.Equal(t, "cate IN (?)", query)
	assert.Equal(t, []interface{}{[]string{"123", "456"}}, args)

	eq = In{
		"cate": []string{"123", "456"},
		"id":   "456",
	}
	query, args = eq.ToWhere()
	assert.Equal(t, "cate IN (?) AND id IN (?)", query)
	assert.Equal(t, []interface{}{[]string{"123", "456"}, "456"}, args)
}

func TestOr(t *testing.T) {
	eq := Or{Eq{"name": "1"}, Eq{"name": "2"}}
	query, args := eq.ToWhere()
	assert.Equal(t, "(name = ? OR name = ?)", query)
	assert.Equal(t, []interface{}{"1", "2"}, args)
}

func TestAnd(t *testing.T) {
	eq := And{Eq{"name": "1"}, Eq{"name": "2"}}
	query, args := eq.ToWhere()
	assert.Equal(t, "(name = ? AND name = ?)", query)
	assert.Equal(t, []interface{}{"1", "2"}, args)
}

func TestToWhere(t *testing.T) {
	exprs := []Expr{
		In{"cate": []string{"123", "456"}},
		Or{Eq{"name": "1"}, Eq{"name": "2"}},
		GtOrEq{"cate": "123"},
	}
	query, args := ToWhere(exprs)
	assert.Equal(t, "cate IN (?) AND (name = ? OR name = ?) AND cate >= ?", query)
	assert.Equal(t, []interface{}{[]string{"123", "456"}, "1", "2", "123"}, args)

	exprs = []Expr{
		In{"cate": []string{"123", "456"}},
	}
	query, args = ToWhere(exprs)
	assert.Equal(t, "cate IN (?)", query)
	assert.Equal(t, []interface{}{[]string{"123", "456"}}, args)
}

func TestGorm(t *testing.T) {
	db, _, _ := getDBMock()

	var where []Expr
	where = append(where, Eq{"Eq": 0})
	where = append(where, Neq{"Neq": 0})
	where = append(where, GtOrEq{"GtOrEq": 0})
	where = append(where, Gt{"Gt": 0})
	where = append(where, Lt{"Lt": 0})
	where = append(where, LtOrEq{"LtOrEq": 0})
	where = append(where, In{"In": []string{"1", "2", "3"}})
	where = append(where, Like{"Like": "test"})
	where = append(where, Or{Eq{"Or": 1}, Eq{"Or": 2}})
	where = append(where, And{Eq{"And1": 1}, Eq{"And2": 2}})

	query, args := ToWhere(where)
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Table("test").Where(query, args...).Find(nil)
	})

	assert.Equal(t, "SELECT * FROM `test` WHERE Eq = 0 AND Neq != 0 AND GtOrEq >= 0 AND Gt > 0 AND Lt < 0 AND LtOrEq <= 0 AND In IN ('1','2','3') AND Like LIKE 'test' AND (Or = 1 OR Or = 2) AND (And1 = 1 AND And2 = 2)", sql)
}

func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	// mock一个*sql.DB对象，不需要连接真实的数据库
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	//测试时不需要真正连接数据库
	gdb, _ := gorm.Open(mysql.New(mysql.Config{
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	return gdb, mock, nil
}
