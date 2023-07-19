package honuadatabase

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	get_metadata    = "SELECT filepath FROM metadata ORDER BY id ASC;"
	add_metadata    = "INSERT INTO metadata(filepath) VALUES ($1);"
	exists_metadata = "SELECT CASE WHEN EXISTS ( SELECT * FROM metadata WHERE filepath = $1) THEN true ELSE false END"
)

func (hdb *HonuaDatabase) exists_metadata(filepath string) (bool, error) {
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	rows, err := hdb.db.Query(exists_metadata, filepath)
	if err != nil {
		log.Printf("An error occured during checking if the metadata %s exists: %s\n", filepath, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the metadata %s exists: %s\n", filepath, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) get_all_done_migrations() []string {
	var migrations []string

	rows, err := hdb.db.Query(get_metadata)
	if err != nil {
		rows.Close()
		return migrations
	}

	for rows.Next() {
		var migration string
		err := rows.Scan(&migration)
		if err != nil {
			rows.Close()
			return migrations
		}
		migrations = append(migrations, migration)
	}

	rows.Close()

	return migrations
}

// reads all the migrations sql files which are in the folder /app/database/files
func (hdb *HonuaDatabase) read_migrations() []string {
	var files []string

	files, err := filepath.Glob(fmt.Sprintf("%s/*.sql", hdb.pathToFiles))
	if err != nil {
		log.Printf("Error running readMigrations %s\n", err.Error())
	}
	return files
}

// adds a migration to the migration table
func (hdb *HonuaDatabase) write_metadata(migration string) {
	_, err := hdb.db.Exec(add_metadata, migration)
	if err != nil {
		log.Printf("Error running writeMetadata %s\n", err.Error())
	}
}

// public Method to start hte Migration
func (hdb *HonuaDatabase) Migrate() {
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	// get all migrations that where done in the past
	var done []string = hdb.get_all_done_migrations()
	// get all migrations which are in the folder /app/database/files
	var migrations []string = hdb.read_migrations()

	// array to store all the migrations that were not already done in the past
	var todo []string = []string{}

	// Iterate through the migration list and add those migrations that have not already been done to the todo list
	for _, migration := range migrations {
		if !string_array_contains_string(migration, done) {
			todo = append(todo, migration)
		}
	}

	// Iterate through the todo list and parse the migrations to string statements. And Execute each statement
	// After that write the migration to the metadata table
	for _, migration := range todo {
		if strings.Contains(migration, "create.sql") {
			log.Printf("Migrate Database skip file %s\n", migration)
			continue
		}
		log.Printf("Migrate Database with file %s\n", migration)
		stmts, err := read_and_parse_sql_file(migration)
		if err != nil {
			log.Printf("Error while Migrating with file %s: %s\n", migration, err.Error())
			continue
		}
		for _, stmt := range stmts {
			_, err := hdb.db.Exec(stmt)
			if err != nil {
				log.Printf("Error while Migrating with file %s: %s\n", migration, err.Error())
				continue
			}
		}
		hdb.write_metadata(migration)
	}
}

// reads the file in with the given filepath and parse it to a string array
func read_and_parse_sql_file(filepath string) ([]string, error) {
	var statements []string

	// Open File
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Scanner zum Zeilenweisen Lesen der Datei erstellen
	scanner := bufio.NewScanner(file)

	// Variable zum Zwischenspeichern von mehrzeiligen Statements
	var statementBuilder strings.Builder

	// Zeilenweise Datei lesen
	for scanner.Scan() {
		line := scanner.Text()

		// Wenn die Zeile mit einem Kommentar beginnt, überspringen
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}

		// Wenn die Zeile ein Teil eines mehrzeiligen Statements ist,
		// an den Builder anhängen
		if strings.HasSuffix(strings.TrimSpace(line), ";") && statementBuilder.Len() > 0 {
			statementBuilder.WriteString(" ")
			statementBuilder.WriteString(strings.TrimSpace(line))
			statement := statementBuilder.String()
			statements = append(statements, statement)
			statementBuilder.Reset()
		} else {
			// Ansonsten die Zeile an den Builder anhängen
			statementBuilder.WriteString(" ")
			statementBuilder.WriteString(strings.TrimSpace(line))
		}
	}

	// Letztes Statement hinzufügen, falls mehrzeiliges Statement am Ende
	if statementBuilder.Len() > 0 {
		statement := statementBuilder.String()
		statements = append(statements, strings.TrimSpace(statement))
	}

	// Fehler beim Scanner prüfen
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return statements, nil
}

// Determine if a string s is in the string array a
func string_array_contains_string(s string, a []string) bool {
	for _, k := range a {
		if s == k {
			return true
		}
	}
	return false
}