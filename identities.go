package honuadatabase

import (
	"database/sql"
	"log"

	"github.com/JonasBordewick/honua-database/models"
)

func (hdb *HonuaDatabase) AddIdentity(identity *models.Identity) error {
	const query = "INSERT INTO identities(identifier, name) VALUES($1, $2);"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	_, err := hdb.db.Exec(query, identity.Id, identity.Name)
	if err != nil {
		log.Printf("An error occured during adding a new identity(identifier=%s, name=%s) to identities: %s\n", identity.Id, identity.Name, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) DeleteIdentity(identifier string) error {
	const query = "DELETE FROM identities WHERE identifier = $1"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	_, err := hdb.db.Exec(query, identifier)
	if err != nil {
		log.Printf("An error occured during deleting the identity %s: %s\n", identifier, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) ExistIdentity(identifier string) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM identities WHERE identifier = $1) THEN true ELSE false END;"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during checking if the identity with id %s exists: %s\n", identifier, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the identity with id %s exists: %s\n", identifier, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) GetIdentity(identifier string) (*models.Identity, error) {
	const query = "SELECT * FROM identities WHERE identifier = $1;"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting identity = %s: %s\n", identifier, err.Error())
		return nil, err
	}

	var result *models.Identity

	for rows.Next() {
		result, err = hdb.make_identity(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting identity = %s: %s\n", identifier, err.Error())
			return nil, err
		}
	}
	rows.Close()

	return result, err
}

func (hdb *HonuaDatabase) GetIdentities() ([]*models.Identity, error) {
	const query = "SELECT * FROM identities"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query)
	if err != nil {
		log.Printf("An error occured during getting all identities: %s\n", err.Error())
		return nil, err
	}

	var result []*models.Identity = []*models.Identity{}

	for rows.Next() {
		identity, err := hdb.make_identity(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all identities: %s\n", err.Error())
			return nil, err
		}
		result = append(result, identity)
	}
	rows.Close()

	return result, err
}

func (hdb *HonuaDatabase) make_identity(rows *sql.Rows) (*models.Identity, error) {
	var id string
	var name string

	err := rows.Scan(&id, &name)
	if err != nil {
		return nil, err
	}

	var result *models.Identity = &models.Identity{
		Id: id,
		Name: name,
	}

	return result, nil
}
