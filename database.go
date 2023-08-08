package honuadatabase

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type HonuaDatabase struct {
	db *sql.DB
	//mutex       sync.Mutex
	pathToFiles string
}

var instance *HonuaDatabase

// Gibt die aktuelle Datenbank Instanz zurück
// Falls noch keine existiert, dann wird eine neue erstellt, dafür muss man die Parameter übergeben
func GetHonuaDatabaseInstance(user, password, host, port, dbname, pathToFiles string) *HonuaDatabase {
	if instance == nil {
		var connStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err) // If any error occure Panic
		}
		if err = db.Ping(); err != nil {
			panic(err) // If any error occure Panic
		}
		log.Println("The Database connection is established")
		instance = &HonuaDatabase{
			db: db,
			//mutex:       sync.Mutex{},
			pathToFiles: pathToFiles,
		}
		err = instance.CreateTables()
		if err != nil {
			panic(err) // If any error occure Panic
		}
		instance.Migrate()
	}
	return instance
}

// Führt im Ordner der Instanz die create.sql file aus
// Dazu wird die File eingelesen und in Statements unterteilt.
// Jedes Statement wird dann ausgeführt. Dadurch, dass die Create Table Statements
// mit einem IF NOT EXIST verbunden sind sollte es keine Fehler geben, wenn diese
// Methode öfter ausgeführt wird.
func (hdb *HonuaDatabase) CreateTables() error {
	stmts, err := read_and_parse_sql_file(fmt.Sprintf("%s/create.sql", hdb.pathToFiles))
	if err != nil {
		log.Printf("Error while reading file %s/create.sql: %s\n", hdb.pathToFiles, err.Error())
		return err
	}
	for _, stmt := range stmts {
		_, err := hdb.db.Exec(stmt)
		if err != nil {
			log.Printf("Error while executing statement %s: %s\n", stmt, err.Error())
			return err
		}
	}

	exists, err := hdb.exists_metadata(fmt.Sprintf("%s/create.sql", hdb.pathToFiles))
	if err != nil {
		return err
	}

	if !exists {
		hdb.write_metadata(fmt.Sprintf("%s/create.sql", hdb.pathToFiles))
	}

	return nil
}

func (hd *HonuaDatabase) CloseDatabase() {
	hd.db.Close()
	instance = nil
	log.Println("The Database Connection is closed")
}
