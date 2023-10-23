package model

import "database/sql"

func MigrateGpsRouteTable(db *sql.DB) error {
	createGpsData := `CREATE TABLE IF NOT EXISTS gps_route_data (
		id SERIAL NOT NULL,
		route_id VARCHAR(100) NOT NULL,
		user_id  VARCHAR(100) NOT NULL,
		total_length int NOT NULL,
   		 total_time data NOT NULL,
        total_calories int NOT NULL,
        location VARCHAR(100)  NULL,
     sports_type int NOT NULL,
     start_time datetime NOT NULL,
     EndTime datetime NOT NULL,
     upload_time datetime NOT NULL,
     locus_url varchar(100) NULL,
     locus_url2 varchar(100)  NULL,
     steps int  NULL,
     file_url VARCHAR(200) NOT NULL,
		PRIMARY KEY (id)
	);`
	_, err := db.Exec(createGpsData)
	return err
}
