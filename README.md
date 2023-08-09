## gorm where 条件生成

```shell
go get -u github.com/dabao-zhao/where-builder
```

```go
var where []where_builder.Expr

where = append(where, where_builder.Eq{"status": 1})
query, args := where_builder.ToWhere(where)

res := conn.Where(query, args...).First(&data)
```