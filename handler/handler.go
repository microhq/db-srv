package handler

import (
	"github.com/micro/db-srv/db"
	mdb "github.com/micro/db-srv/proto/db"
	"github.com/micro/go-micro/errors"

	"golang.org/x/net/context"
)

type DB struct{}

func validateDB(method string, d *mdb.Database) error {
	if d == nil {
		return errors.BadRequest("go.micro.srv.db."+method, "invalid database")
	}

	if len(d.Name) == 0 {
		return errors.BadRequest("go.micro.srv.db."+method, "database is blank")
	}
	if len(d.Table) == 0 {
		return errors.BadRequest("go.micro.srv.db."+method, "table is blank")
	}

	// TODO: check exists

	return nil
}

func (d *DB) Read(ctx context.Context, req *mdb.ReadRequest, rsp *mdb.ReadResponse) error {
	if err := validateDB("DB.Read", req.Database); err != nil {
		return err
	}

	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.db.DB.Read", "invalid id")
	}

	r, err := db.Read(req.Database, req.Id)
	if err != nil && err == db.ErrNotFound {
		return errors.NotFound("go.micro.srv.db.DB.Read", "not found")
	} else if err != nil {
		return errors.InternalServerError("go.micro.srv.db.DB.Read", err.Error())
	}

	rsp.Record = r

	return nil
}

func (d *DB) Create(ctx context.Context, req *mdb.CreateRequest, rsp *mdb.CreateResponse) error {
	if req.Record == nil {
		return errors.BadRequest("go.micro.srv.db.DB.Create", "invalid record")
	}

	if err := validateDB("DB.Create", req.Record.Database); err != nil {
		return err
	}

	if len(req.Record.Id) == 0 {
		return errors.BadRequest("go.micro.srv.db.DB.Create", "invalid id")
	}

	if err := db.Create(req.Record.Database, req.Record); err != nil {
		return errors.InternalServerError("go.micro.srv.db.DB.Create", err.Error())
	}

	return nil
}

func (d *DB) Update(ctx context.Context, req *mdb.UpdateRequest, rsp *mdb.UpdateResponse) error {
	if req.Record == nil {
		return errors.BadRequest("go.micro.srv.db.DB.Update", "invalid record")
	}

	if err := validateDB("DB.Update", req.Record.Database); err != nil {
		return err
	}

	if len(req.Record.Id) == 0 {
		return errors.BadRequest("go.micro.srv.db.DB.Update", "invalid id")
	}

	if err := db.Update(req.Record.Database, req.Record); err != nil && err == db.ErrNotFound {
		return errors.NotFound("go.micro.srv.db.DB.Update", "not found")
	} else if err != nil {
		return errors.InternalServerError("go.micro.srv.db.DB.Update", err.Error())
	}

	return nil
}

func (d *DB) Delete(ctx context.Context, req *mdb.DeleteRequest, rsp *mdb.DeleteResponse) error {
	if err := validateDB("DB.Delete", req.Database); err != nil {
		return err
	}

	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.db.DB.Delete", "invalid id")
	}

	if err := db.Delete(req.Database, req.Id); err != nil && err == db.ErrNotFound {
		return nil
	} else if err != nil {
		return errors.InternalServerError("go.micro.srv.db.DB.Delete", err.Error())
	}

	return nil
}

func (d *DB) Search(ctx context.Context, req *mdb.SearchRequest, rsp *mdb.SearchResponse) error {
	if err := validateDB("DB.Search", req.Database); err != nil {
		return err
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	metadata := map[string]interface{}{}

	for k, v := range req.Metadata {
		metadata[k] = v
	}

	r, err := db.Search(req.Database, metadata, req.Limit, req.Offset)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.DB.Search", err.Error())
	}

	rsp.Records = r

	return nil
}
