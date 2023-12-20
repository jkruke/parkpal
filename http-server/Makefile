createdb:
	sqlite3 parking_lots.db;

migratecreate: 
	migrate create -ext sql -dir db/migrations add_parking_lots_table

migrateup:
	sqlite3 parking_lots.db < ./db/migrations/20231220092537_add_parking_lots_table.up.sql

migratedown:
	sqlite3 parking_lots.db < ./db/migrations/20231220092537_add_parking_lots_table.down.sql

dbinsert:
	sqlite3 parking_lots.db < ./db/migrations/insert.sql
