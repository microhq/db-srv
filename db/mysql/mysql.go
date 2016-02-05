package mysql

import (
	"github.com/micro/db-srv/db"
	mdb "github.com/micro/db-srv/proto/db"
	"github.com/micro/go-micro/registry"
)

func init() {
	db.Drivers["mysql"] = new(mysqlDriver)
}

type mysqlDriver struct{}

type mysqlDB struct{}

func (d *mysqlDriver) NewDB(nodes ...*registry.Node) (db.DB, error) {
	return nil, nil
}

func (d *mysqlDB) Init(mdb *mdb.Database) error {
	return nil
}

func (d *mysqlDB) Close() error {
	return nil
}

func (d *mysqlDB) Read(id string) (*mdb.Record, error) {
	return nil, nil
}

func (d *mysqlDB) Create(r *mdb.Record) error {
	return nil
}

func (d *mysqlDB) Update(r *mdb.Record) error {
	return nil
}

func (d *mysqlDB) Delete(id string) error {
	return nil
}

func (d *mysqlDB) Search(md map[string]interface{}, limit, offset int64) ([]*mdb.Record, error) {
	return nil, nil
}
