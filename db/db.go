package db

import (
	"errors"
	"strings"

	mdb "github.com/micro/db-srv/proto/db"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
)

// Database Driver
type Driver interface {
	NewDB(nodes ...*registry.Node) (DB, error)
}

// An initialised DB connection
// Must call Init to load Database/Table
type DB interface {
	// Initialise the database
	// If the database doesn't exist
	// will throw an error
	Init(mdb *mdb.Database) error
	// Close the connection
	Close() error
	// Query commands
	Read(id string) (*mdb.Record, error)
	Create(mdb *mdb.Record) error
	Update(mdb *mdb.Record) error
	Delete(id string) error
	Search(md map[string]string, limit, offset int64) ([]*mdb.Record, error)
}

type db struct {
	selector  selector.Selector
	namespace string
	drivers   map[string]Driver
	driverKey string
}

var (
	DefaultDB *db

	// Prefix for lookup in registry
	DBServiceNamespace = "go.micro.db"

	// Used to lookup the metadata for the driver
	DBDriverKey = "driver"

	// supported drivers: mysql, elastic, cassandra
	Drivers = map[string]Driver{}

	// Errors
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotAvailable  = errors.New("not available")
)

func newDB(s selector.Selector) *db {
	return &db{
		selector:  s,
		namespace: DBServiceNamespace,
		drivers:   Drivers,
		driverKey: DBDriverKey,
	}
}

func (d *db) name(db *mdb.Database) string {
	return strings.Join([]string{d.namespace, db.Name}, ".")
}

func (d *db) lookup(db *mdb.Database) (DB, error) {
	next, err := d.selector.Select(d.name(db))
	if err != nil {
		return nil, err
	}

	var id string

	// TODO: create a node list rather than connecting to one
	for {
		node, err := next()
		if err != nil {
			return nil, err
		}

		// seen all?
		if node.Id == id {
			return nil, ErrNotAvailable
		}

		id = node.Id

		// is the driver set?
		dv, ok := node.Metadata[d.driverKey]
		if !ok {
			continue
		}

		// is the driver supported?
		dr, ok := d.drivers[dv]
		if !ok {
			continue
		}

		conn, err := dr.NewDB(node)
		if err != nil {
			return nil, err
		}

		if err := conn.Init(db); err != nil {
			return nil, err
		}

		return conn, nil
	}

	return nil, ErrNotAvailable
}

func (d *db) Read(db *mdb.Database, id string) (*mdb.Record, error) {
	dr, err := d.lookup(db)
	if err != nil {
		return nil, err
	}
	defer dr.Close()
	return dr.Read(id)
}

func (d *db) Create(db *mdb.Database, r *mdb.Record) error {
	dr, err := d.lookup(db)
	if err != nil {
		return err
	}
	defer dr.Close()
	return dr.Create(r)
}

func (d *db) Update(db *mdb.Database, r *mdb.Record) error {
	dr, err := d.lookup(db)
	if err != nil {
		return err
	}
	defer dr.Close()
	return dr.Update(r)
}

func (d *db) Delete(db *mdb.Database, id string) error {
	dr, err := d.lookup(db)
	if err != nil {
		return err
	}
	defer dr.Close()
	return dr.Delete(id)
}

func (d *db) Search(db *mdb.Database, md map[string]string, limit, offset int64) ([]*mdb.Record, error) {
	dr, err := d.lookup(db)
	if err != nil {
		return nil, err
	}
	defer dr.Close()
	return dr.Search(md, limit, offset)
}

func Init(s selector.Selector) error {
	DefaultDB = newDB(s)
	return nil
}

func Read(db *mdb.Database, id string) (*mdb.Record, error) {
	return DefaultDB.Read(db, id)
}

func Create(db *mdb.Database, r *mdb.Record) error {
	return DefaultDB.Create(db, r)
}

func Update(db *mdb.Database, r *mdb.Record) error {
	return DefaultDB.Update(db, r)
}

func Delete(db *mdb.Database, id string) error {
	return DefaultDB.Delete(db, id)
}

func Search(db *mdb.Database, md map[string]string, limit, offset int64) ([]*mdb.Record, error) {
	return DefaultDB.Search(db, md, limit, offset)
}
