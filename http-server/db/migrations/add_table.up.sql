CREATE TABLE IF NOT EXISTS parking_lots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    latitude DOUBLE,
    longitude DOUBLE,
    congestionRate DOUBLE,
    totalSpace INTEGER
);

CREATE TABLE IF NOT EXISTS license_plates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    parkinglot_id INTEGER,
    FOREIGN KEY (parkinglot_id) REFERENCES parking_lots(id) ON DELETE CASCADE
);
