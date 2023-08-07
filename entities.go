package honuadatabase

import (
	"database/sql"
	"errors"
	"log"

	"github.com/JonasBordewick/honua-database/models"
)

func (hdb *HonuaDatabase) GetEntity(id int) (*models.Entity, error) {
	const query = "SELECT * FROM entities WHERE id=$1;"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, id)
	if err != nil {
		log.Printf("An error occured during getting entity: %s\n", err.Error())
		return nil, err
	}

	var result *models.Entity

	for rows.Next() {
		entity, err := hdb.make_entity(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting entity: %s\n", err.Error())
			return nil, err
		}
		result = entity
	}

	rows.Close()

	return result, nil
}

// Fügt eine neue Entität zur Datenbank hinzu
func (hdb *HonuaDatabase) AddEntity(entity *models.Entity) error {
	const query = `
INSERT INTO entities(
	identity, entity_id, name,
	is_device, allow_rules, has_attribute,
	attribute, is_victron_sensor, has_numeric_state
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
`

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	var attributeString sql.NullString = sql.NullString{
		Valid:  entity.Attribute != "",
		String: entity.Attribute,
	}

	_, err := hdb.db.Exec(query, entity.IdentityId, entity.EntityId, entity.Name, entity.IsDevice, entity.AllowRules, entity.HasAttribute, attributeString, entity.IsVictronSensor, entity.HasNumericState)

	if err != nil {
		log.Printf("An error occured during adding a new entitiy to table entities: %s\n", err.Error())
	}
	return err
}

// Löscht eine Enität mit der ID im Parameter
func (hdb *HonuaDatabase) DeleteEntity(id int) error {
	const query = "DELETE FROM entities WHERE id = $1;"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	_, err := hdb.db.Exec(query, id)
	if err != nil {
		log.Printf("An error occured during deleting the entity with id = %d: %s\n", id, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) EditEntity(identifier string, entity *models.Entity) error {
	const query = `
UPDATE entities
SET name = $1, is_device = $2, allow_rules = $3, has_attribute = $4, attribute = $5, is_victron_sensor = $6, has_numeric_state = $7
WHERE identity = $8 AND entity_id = $9;
	`
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	var attributeString sql.NullString = sql.NullString{
		Valid:  entity.Attribute != "",
		String: entity.Attribute,
	}

	entity.HasAttribute = attributeString.Valid

	_, err := hdb.db.Exec(query, entity.Name, entity.IsDevice, entity.AllowRules, entity.HasAttribute, attributeString, entity.IsVictronSensor, entity.HasNumericState, entity.IdentityId, entity.EntityId)

	if err != nil {
		log.Printf("An error occured during editity entitiy: %s\n", err.Error())
	}
	return err
}

// Checkt, ob eine Entität existiert die einen bestimmten Identifier und eine EntityID hat
func (hdb *HonuaDatabase) ExistEntity(identifier, entityId string) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM entities WHERE identity = $1 AND entity_id = $2) THEN true ELSE false END"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, identifier, entityId)
	if err != nil {
		log.Printf("An error occured during checking if the entity %s exists in %s: %s\n", identifier, entityId, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the entity %s exists in %s: %s\n", identifier, entityId, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) GetIdOfEntity(identifier, entityId string) (int, error) {
	const query = "SELECT id FROM entities WHERE identity = $1 AND entity_id = $2"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	rows, err := hdb.db.Query(query, identifier, entityId)
	if err != nil {
		log.Printf("An error occured during checking the id of entity (%s, %s): %s\n", identifier, entityId, err.Error())
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking the id of entity (%s, %s): %s\n", identifier, entityId, err.Error())
			return -1, err
		}
	}

	rows.Close()

	return id, nil
}

func (hdb *HonuaDatabase) GetEntities(identifier string) ([]*models.Entity, error) {
	const query = "SELECT * FROM entities WHERE identity = $1;"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting all entities of identity = %s: %s\n", identifier, err.Error())
		return nil, err
	}

	var result []*models.Entity = []*models.Entity{}

	for rows.Next() {
		entity, err := hdb.make_entity(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all entities of identity = %s: %s\n", identifier, err.Error())
			return nil, err
		}
		result = append(result, entity)
	}

	rows.Close()

	return result, nil
}

func (hdb *HonuaDatabase) GetEntitiesWhereRulesAreAllowed(identifier string) ([]*models.Entity, error) {
	const query = "SELECT * FROM entities WHERE identity = $1 AND allow_rules;"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting all entities of identity = %s: %s\n", identifier, err.Error())
		return nil, err
	}

	var result []*models.Entity = []*models.Entity{}

	for rows.Next() {
		entity, err := hdb.make_entity(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all entities of identity = %s: %s\n", identifier, err.Error())
			return nil, err
		}
		result = append(result, entity)
	}

	rows.Close()

	return result, nil

}

func (hdb *HonuaDatabase) GetEntitiesWithoutAnyRule() ([]*models.Entity, error) {
	return nil, errors.New("not implemented yet")
}

func (hdb *HonuaDatabase) GetVictronEntities(identifier string) ([]*models.Entity, error) {
	const query = "SELECT * FROM entities WHERE identity = $1 AND is_victron_sensor;"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting all entities of identity = %s: %s\n", identifier, err.Error())
		return nil, err
	}

	var result []*models.Entity = []*models.Entity{}

	for rows.Next() {
		entity, err := hdb.make_entity(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all entities of identity = %s: %s\n", identifier, err.Error())
			return nil, err
		}
		result = append(result, entity)
	}

	rows.Close()

	return result, nil
}

func (hdb *HonuaDatabase) make_entity(rows *sql.Rows) (*models.Entity, error) {
	var id int
	var identity string
	var entityID string
	var name string
	var isDevice bool
	var allowRules bool
	var hasAttribute bool
	var attribute sql.NullString
	var isVictronSensor bool
	var hasNumericState bool


	err := rows.Scan(&id, &identity, &entityID, &name, &isDevice, &allowRules, &hasAttribute, &attribute, &isVictronSensor, &hasNumericState)
	if err != nil {
		return nil, err
	}


	var result *models.Entity = &models.Entity{
		Id: id,
		IdentityId: identity,
		Name: name,
		EntityId: entityID,
		IsDevice: isDevice,
		AllowRules: allowRules,
		HasAttribute: hasAttribute,
		IsVictronSensor: isVictronSensor,
		HasNumericState: hasNumericState,
	}

	if hasAttribute && attribute.Valid {
		result.Attribute = attribute.String
	} else {
		result.HasAttribute = false
		result.Attribute = ""
	}
	return result, nil
}
