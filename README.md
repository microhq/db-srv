# DB Server

**Experimental. Run away. RUN AWAY**

The DB server is an experimental proxy layer for backend databases. It allows the platform services to be leveraged 
as part of the system. Dynamic config, auth, routing, selector, etc.

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start a mysql database

4. Download and start the service

	```shell
	go get github.com/micro/db-srv
	db-srv --database_url="root:root@tcp(192.168.99.100:3306)/db"
	```

	OR as a docker container

	```shell
	docker run microhq/db-srv --database_url="root:root@tcp(192.168.99.100:3306)/db" --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
DB server implements the following RPC Methods

DB
- Read
- Create
- Update
- Delete
- Search
