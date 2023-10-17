package cloudsql

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"
)

var (
	db   *sql.DB
	once sync.Once
)

// getDB lazily instantiates a database connection pool. Users of Cloud Run or
// Cloud Functions may wish to skip this lazy instantiation and connect as soon
// as the function is loaded. This is primarily to help testing.
//获取建立好的链接
func GetDB() *sql.DB {
	once.Do(func() {
		db = mustConnect()
	})
	return db
}

//创建一个表
func migrateDB(db *sql.DB) error {
	createVotes := `CREATE TABLE IF NOT EXISTS votes (
		id SERIAL NOT NULL,
		created_at datetime NOT NULL,
		candidate VARCHAR(6) NOT NULL,
		PRIMARY KEY (id)
	);`
	_, err := db.Exec(createVotes)
	return err
}

//创建一个表
func migrateGpsTable(db *sql.DB) error {
	createGpsData := `CREATE TABLE IF NOT EXISTS gps_data (
		id SERIAL NOT NULL,
		user_id datetime NOT NULL,
		total_length int NOT NULL,
		PRIMARY KEY (id)
	);`
	_, err := db.Exec(createGpsData)
	return err
}

// 建立链接
func mustConnect() *sql.DB {
	var (
		db  *sql.DB
		err error
	)

	// Use a TCP socket when INSTANCE_HOST (e.g., 127.0.0.1) is defined
	if os.Getenv("INSTANCE_HOST") != "" {
		db, err = connectTCPSocket()
		if err != nil {
			log.Fatalf("connectTCPSocket: unable to connect: %s", err)
		}
	}
	// Use a Unix socket when INSTANCE_UNIX_SOCKET (e.g., /cloudsql/proj:region:instance) is defined.
	if os.Getenv("INSTANCE_UNIX_SOCKET") != "" {
		db, err = connectUnixSocket()
		if err != nil {
			log.Fatalf("connectUnixSocket: unable to connect: %s", err)
		}
	}

	// Use the connector when INSTANCE_CONNECTION_NAME (proj:region:instance) is defined.
	if os.Getenv("INSTANCE_CONNECTION_NAME") != "" {
		if os.Getenv("DB_USER") == "" && os.Getenv("DB_IAM_USER") == "" {
			log.Fatal("Warning: One of DB_USER or DB_IAM_USER must be defined")
		}
		// Use IAM Authentication (recommended) if DB_IAM_USER is set
		if os.Getenv("DB_IAM_USER") != "" {
			db, err = connectWithConnectorIAMAuthN()
		} else {
			db, err = connectWithConnector()
		}
		if err != nil {
			log.Fatalf("connectConnector: unable to connect: %s", err)
		}
	}

	if db == nil {
		log.Fatal("Missing database connection type. Please define one of INSTANCE_HOST, INSTANCE_UNIX_SOCKET, or INSTANCE_CONNECTION_NAME")
	}

	if err := migrateDB(db); err != nil {
		log.Fatalf("unable to create table: %s", err)
	}
	if err := migrateGpsTable(db); err != nil {
		log.Fatalf("unable to create table: %s", err)
	}

	return db
}

func configureConnectionPool(db *sql.DB) {
	// [START cloud_sql_mysql_databasesql_limit]
	// Set maximum number of connections in idle connection pool.
	db.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	db.SetMaxOpenConns(7)
	// [END cloud_sql_mysql_databasesql_limit]

	// [START cloud_sql_mysql_databasesql_lifetime]
	// Set Maximum time (in seconds) that a connection can remain open.
	db.SetConnMaxLifetime(1800 * time.Second)
	// [END cloud_sql_mysql_databasesql_lifetime]

	// [START cloud_sql_mysql_databasesql_backoff]
	// database/sql does not support specifying backoff
	// [END cloud_sql_mysql_databasesql_backoff]
	// [START cloud_sql_mysql_databasesql_timeout]
	// The database/sql package currently doesn't offer any functionality to
	// configure connection timeout.
	// [END cloud_sql_mysql_databasesql_timeout]
}
