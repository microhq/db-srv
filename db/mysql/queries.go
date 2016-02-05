package mysql

var (
	mysqlQueries = map[string]string{
		"delete": "DELETE from %s.%s where id = ? limit 1",
		"create": `INSERT into %s.%s (id, created, updated, metadata, bytes) values (?, ?, ?, ?, ?)`,
		"update": "UPDATE %s.%s set updated = ?, metadata = ?, bytes = ? where id = ?",
		"search": "SELECT id, created, updated, metadata, bytes from %s.%s limit ? offset ?",
		"read":   "SELECT id, created, updated, metadata, bytes from %s.%s where id = ? limit 1",
	}

	searchMetadataQ = "SELECT id, created, updated, metadata, bytes from %s.%s where metadata like ?"
)
