createdb:
	sqlite3 parkpal.db;

migratecreate: 
	migrate create -ext sql -dir db/migrations add_table

migrateup:
	sqlite3 parkpal.db < ./db/migrations/add_table.up.sql

migratedown:
	sqlite3 parkpal.db < ./db/migrations/add_table.down.sql

dbinsert:
	sqlite3 parkpal.db < ./db/migrations/insert.sql
