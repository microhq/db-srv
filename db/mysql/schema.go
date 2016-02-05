package mysql

var (
	mysqlSchema = `CREATE TABLE IF NOT EXISTS %s (
id varchar(36) primary key,
created integer,
updated integer,
metadata text,
bytes blob,
index(created),
index(updated));`
)
