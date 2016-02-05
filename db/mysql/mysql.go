package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/db-srv/db"
	mdb "github.com/micro/db-srv/proto/db"
	"github.com/micro/go-micro/registry"
)

type mysqlDriver struct{}

type mysqlDB struct {
	url string

	sync.RWMutex
	name    string
	table   string
	conn    *sql.DB
	queries map[string]*sql.Stmt
}

var (
	DBUser = "root"
	DBPass = ""
)

func init() {
	db.Drivers["mysql"] = new(mysqlDriver)
}

func (d *mysqlDriver) NewDB(nodes ...*registry.Node) (db.DB, error) {
	if len(nodes) == 0 {
		return nil, db.ErrNotAvailable
	}

	url := fmt.Sprintf("tcp(%s:%d)/", nodes[0].Address, nodes[0].Port)

	// add credentials
	// TODO: take database credentials
	if len(DBUser) > 0 && len(DBPass) > 0 {
		url = fmt.Sprintf("%s:%s@%s", DBUser, DBPass, url)
	} else if len(DBUser) > 0 {
		url = fmt.Sprintf("%s@%s", DBUser, url)
	}

	// test the connection
	conn, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return &mysqlDB{
		url:     url,
		queries: make(map[string]*sql.Stmt),
	}, nil
}

func (d *mysqlDB) Init(mdb *mdb.Database) error {
	d.Lock()
	defer d.Unlock()

	// Create a conn to initialise the database
	conn, err := sql.Open("mysql", d.url)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create the database
	// TODO: we really shouldn't create the databases automatically
	// Should be external managed process.
	// Maybe a handler on this service
	if _, err := conn.Exec("CREATE DATABASE IF NOT EXISTS " + mdb.Name); err != nil {
		return err
	}

	// create connection
	dconn, err := sql.Open("mysql", d.url+mdb.Name)
	if err != nil {
		return err
	}

	d.conn = dconn

	if _, err = d.conn.Exec(fmt.Sprintf(mysqlSchema, mdb.Table)); err != nil {
		return err
	}

	for query, statement := range mysqlQueries {
		prepared, err := d.conn.Prepare(fmt.Sprintf(statement, mdb.Name, mdb.Table))
		if err != nil {
			return err
		}
		d.queries[query] = prepared
	}

	d.name = mdb.Name
	d.table = mdb.Table

	return nil
}

func (d *mysqlDB) Close() error {
	d.Lock()
	defer d.Unlock()
	return d.conn.Close()
}

func (d *mysqlDB) Read(id string) (*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	r := &mdb.Record{}
	row := d.queries["read"].QueryRow(id)

	var meta []byte
	if err := row.Scan(&r.Id, &r.Created, &r.Updated, &meta, &r.Bytes); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(meta, &r.Metadata); err != nil {
		return nil, err
	}

	return r, nil
}

func (d *mysqlDB) Create(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()

	meta, err := json.Marshal(r.Metadata)
	if err != nil {
		return err
	}
	r.Created = time.Now().Unix()
	r.Updated = time.Now().Unix()

	_, err = d.queries["create"].Exec(r.Id, r.Created, r.Updated, string(meta), []byte(r.Bytes))
	return err
}

func (d *mysqlDB) Update(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()

	meta, err := json.Marshal(r.Metadata)
	if err != nil {
		return err
	}
	r.Updated = time.Now().Unix()

	_, err = d.queries["update"].Exec(r.Updated, string(meta), []byte(r.Bytes), r.Id)

	return nil
}

func (d *mysqlDB) Delete(id string) error {
	d.RLock()
	defer d.RUnlock()
	_, err := d.queries["delete"].Exec(id)
	return err
}

func (d *mysqlDB) Search(md map[string]string, limit, offset int64) ([]*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	var rows *sql.Rows
	var err error

	if len(md) > 0 {
		// THIS IS SUPER CRUFT
		// TODO: DONT DO THIS
		// Note: Tried to use mariadb dynamic columns. They suck.
		var query string
		var args []interface{}

		// create statement for each key-val pair
		for k, v := range md {
			if len(query) == 0 {
				query += " "
			} else {
				query += "AND metadata like ? "
			}
			args = append(args, fmt.Sprintf(`%%"%s":"%s"%%`, k, v))
		}

		// append limit offset
		args = append(args, limit, offset)
		query += " limit ? offset ?"
		query = fmt.Sprintf(searchMetadataQ, d.name, d.table) + query

		// doe the query
		rows, err = d.conn.Query(query, args...)
	} else {
		rows, err = d.queries["search"].Query(limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*mdb.Record

	for rows.Next() {
		r := &mdb.Record{}
		var meta []byte
		if err := rows.Scan(&r.Id, &r.Created, &r.Updated, &meta, &r.Bytes); err != nil {
			if err == sql.ErrNoRows {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		if err := json.Unmarshal(meta, &r.Metadata); err != nil {
			return nil, err
		}
		records = append(records, r)

	}
	if rows.Err() != nil {
		return nil, err
	}

	return records, nil
}
