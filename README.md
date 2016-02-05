# DB Server

**Experimental. Run away. RUN AWAY**

The DB server is an experimental proxy layer for backend databases. It allows the platform services to be leveraged 
as part of the system. Dynamic config, auth, routing, selector, etc.

The DB uses the registry to find supported databases for proxying. It expects the metadata key "driver" to be set 
to the type of database. It should match one of the supported drivers.

## Database Drivers

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

4. Download and start the service

	```shell
	go get github.com/micro/db-srv
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
