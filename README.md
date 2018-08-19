# DB Server

**Experimental. Run away. RUN AWAY**

The DB server is an experimental proxy layer for backend databases. It allows the platform services to be leveraged 
as part of the system. Dynamic config, auth, routing, selector, etc.

The DB uses the registry to find supported databases for proxying. It expects the metadata key "driver" to be set 
to the type of database. It should match one of the supported drivers.

## Database Drivers

- mysql (mariadb)

future
- elasticsearch
- cassandra

...

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start one of the supported databases (mysql, elasticsearch, ...) and register with the registry.
	```
	Example. Register location of the **foo** database hosted by mysql

	$ micro register service '{"name": "go.micro.db.foo", "version": "0.0.1", "nodes": [{"id": "foo-1", "address": "127.0.0.1", "port": 3306, "metadata": {"driver": "mysql"}}]}'
	```
4. Download and start the service

	```shell
	go get github.com/microhq/db-srv
	db-srv --database_service_namespace=go.micro.db
	```

	OR as a docker container

	```shell
	docker run microhq/db-srv --database_service_namespace=go.micro.db --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
DB server implements the following RPC Methods

DB
- Read
- Create
- Update
- Delete
- Search

### DB.Create

```
micro query go.micro.srv.db DB.Create '{"database": {"name": "foo", "table": "bar"}, "record": {"id": "e7add322-e069-44c2-b920-c4fbfd62e6b5", "metadata": {"key": "value"}}}'
```

### DB.Read

```
micro query go.micro.srv.db DB.Read '{"database": {"name": "foo", "table": "bar"}, "id": "e7add322-e069-44c2-b920-c4fbfd62e6b5"}'

{
	"record": {
		"created": 1.454704366e+09,
		"id": "e7add322-e069-44c2-b920-c4fbfd62e6b5",
		"metadata": {
			"key": "value"
		},
		"updated": 1.454704366e+09
	}
}
```

### DB.Search

```
micro query go.micro.srv.db DB.Search '{"database": {"name": "foo", "table": "bar"}}'

{
	"records": [
		{
			"created": 1.454704366e+09,
			"id": "e7add322-e069-44c2-b920-c4fbfd62e6b5",
			"metadata": {
				"key": "value"
			},
			"updated": 1.454704366e+09
		}
	]
}
```

### DB.Update

```
micro query go.micro.srv.db DB.Update '{"database": {"name": "foo", "table": "bar"}, "record": {"id": "e7add322-e069-44c2-b920-c4fbfd62e6b5", "metadata": {"key": "value", "key2": "value2"}}}'
```

### DB.Delete

```
micro query go.micro.srv.db DB.Delete '{"database": {"name": "foo", "table": "bar"}, "id": "e7add322-e069-44c2-b920-c4fbfd62e6b5"}'
```
