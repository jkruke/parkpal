const express = require('express');
const app = express();
const sqlite3 = require('sqlite3');

app.use(express.json());

// Connecting Database
let db = new sqlite3.Database("database.db", (err) => {
    if (err) {
        console.log("Error Occurred - " + err.message);
    } else {
        console.log("DataBase Connected");
    }
})

app.get('/parking-lot', (req, res) => {
    const id = req.query.id
    db.all(`SELECT * FROM ParkingLots WHERE id = ${id}`, (err, data) => {
        if (err) {
            console.error("Could not select from table", err)
            res.status(404).send(err.message)
            return
        }
        if (data.length === 0) {
            res.status(400).send("Unknown parking lot ID!")
            return
        }
        res.send(data)
    })
});

app.put('/parking-lot', (req, res) => {
    const parkingLot = req.body
    db.run(`UPDATE ParkingLots SET name = "${parkingLot.name}" WHERE id = ${parkingLot.id}`, (err) => {
        if (err) {
            console.error("Could not update table", err)
            res.status(400).send(err.message)
            return
        }
        res.send("Successfully updated name for parking lot " + parkingLot.id)
    })
})

app.post('/parking-lot', (req, res) => {
    const parkingLot = req.body
    db.run(`INSERT INTO ParkingLots (id, name, bikeCount) VALUES (${parkingLot.id}, "${parkingLot.name}", 0)`, (err) => {
        if (err) {
            console.error("Could not insert into table", err)
            res.status(400).send(err.message)
            return
        }
        res.send("Successfully added new parking lot: " + JSON.stringify(parkingLot))
    })
})

app.delete('/parking-lot', (req, res) => {
    const id = req.query.id
    db.run(`DELETE FROM ParkingLots WHERE id = ${id}`, (err) => {
        if (err) {
            console.error("Could not delete from table", err)
            res.status(400).send(err.message)
            return
        }
        res.send("Successfully deleted parking lot " + id)
    })
})

function dbInitialSetup() {
    db.all('SELECT * FROM ParkingLots;', (err, data) => {
        if (err) {
            const createQuery = 'CREATE TABLE ParkingLots ( id NUMBER PRIMARY KEY, name VARCHAR(100), bikeCount NUMBER);';
            db.run(createQuery, (err) => {
                if (err) {
                    console.error("Could not create table!", err)
                    return
                }
                console.log("Table Created");
                const insertQuery = 'INSERT INTO ParkingLots (id, name, bikeCount) VALUES (1 , "D8", 123);'
                db.run(insertQuery, (err) => {
                    if (err) {
                        console.error("Could not insert data!", err)
                        return
                    }
                    console.log("Insertion Done");
                });
            });
        }
    });
}

// app.use(express.json())
app.listen(8080, () => {
    console.log("server started");
    dbInitialSetup();
});
