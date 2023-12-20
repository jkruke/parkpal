CREATE TABLE IF NOT EXISTS parking_lots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    latitude DOUBLE,
    longitude DOUBLE,
    bikeCount INTEGER,
    congestionRate DOUBLE,
    totalSpace INTEGER
);
