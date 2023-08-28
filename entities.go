package honuadatabase

import (
	"database/sql"
	"errors"
	"log"

	"github.com/JonasBordewick/honua-database/models"
)

// Make sure that before calling this method, you have already been locked. This method does not lock
func (hdb *HonuaDatabase) get_entity_id(identifier string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM entities WHERE identity = $1) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting id of entity in %s: %s\n", identifier, err.Error())
		return -1, err
	}

	var exist_identity bool = false

	for rows.Next() {
		err = rows.Scan(&exist_identity)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of entity in %s: %s\n", identifier, err.Error())
			return -1, err
		}
	}

	rows.Close()

	if !exist_identity {
		return 0, nil
	}

	query = "SELECT MAX(id) FROM entities WHERE identity = $1;"

	rows, err = hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting id of entity in %s: %s\n", identifier, err.Error())
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of entity in %s: %s\n", identifier, err.Error())
			return -1, err
		}
	}
	rows.Close()

	if id == -1 {
		return -1, errors.New("something went wrong during getting id of entity")
	}

	id = id + 1

	return id, nil
}

func (hdb *HonuaDatabase) GetEntity(identity string, id int) (*models.Entity, error) {
	const query = "SELECT * FROM entities WHERE id=$1 AND identity=$2;"

	rows, err := hdb.db.Query(query, id, identity)
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
	id, identity, entity_id, name,
	is_device, allow_rules, has_attribute,
	attribute, is_victron_sensor, sensor_type, has_numeric_state
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
`
	var attributeString sql.NullString = sql.NullString{
		Valid:  entity.Attribute != "",
		String: entity.Attribute,
	}

	id, err := hdb.get_entity_id(entity.IdentityId)
	if err != nil {
		log.Printf("An error occured during adding a new entitiy to table entities: %s\n", err.Error())
		return err
	}

	log.Printf("ID %d\n", id)

	_, err = hdb.db.Exec(query, id, entity.IdentityId, entity.EntityId, entity.Name, entity.IsDevice, entity.AllowRules, entity.HasAttribute, attributeString, entity.IsVictronSensor, entity.SensorType, entity.HasNumericState)

	if err != nil {
		log.Printf("An error occured during adding a new entitiy to table entities: %s\n", err.Error())
	}
	return err
}

// Löscht eine Enität mit der ID im Parameter
func (hdb *HonuaDatabase) DeleteEntity(id int, identity string) error {
	const query = "DELETE FROM entities WHERE identity=$1 AND id = $2;"

	_, err := hdb.db.Exec(query, identity, id)
	if err != nil {
		log.Printf("An error occured during deleting the entity with id = %d: %s\n", id, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) EditEntity(identifier string, entity *models.Entity) error {
	const query = `
UPDATE entities
SET name = $1, is_device = $2, allow_rules = $3, has_attribute = $4, attribute = $5, is_victron_sensor = $6, sensor_type = $7, has_numeric_state = $8
WHERE identity = $9 AND entity_id = $10;
	`

	var attributeString sql.NullString = sql.NullString{
		Valid:  entity.Attribute != "",
		String: entity.Attribute,
	}

	entity.HasAttribute = attributeString.Valid

	_, err := hdb.db.Exec(query, entity.Name, entity.IsDevice, entity.AllowRules, entity.HasAttribute, attributeString, entity.IsVictronSensor, entity.SensorType, entity.HasNumericState, entity.IdentityId, entity.EntityId)

	if err != nil {
		log.Printf("An error occured during editity entitiy: %s\n", err.Error())
	}
	return err
}

// Checkt, ob eine Entität existiert die einen bestimmten Identifier und eine EntityID hat
func (hdb *HonuaDatabase) ExistEntity(identifier string, id int, hasAttribute bool, attribute string) (bool, error) {

	if hasAttribute {
		const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM entities WHERE identity = $1 AND id = $2 AND has_attribute AND attribute = $3) THEN true ELSE false END"

		rows, err := hdb.db.Query(query, identifier, id, attribute)
		if err != nil {
			log.Printf("An error occured during checking if the entity %d exists in %s: %s\n", id, identifier, err.Error())
			return false, err
		}

		var state bool = false

		for rows.Next() {
			err = rows.Scan(&state)
			if err != nil {
				rows.Close()
				log.Printf("An error occured during checking if the entity %d exists in %s: %s\n", id, identifier, err.Error())
				return false, err
			}
		}

		rows.Close()

		return state, nil
	}

	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM entities WHERE identity = $1 AND id = $2 AND NOT has_attribute) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identifier, id)
	if err != nil {
		log.Printf("An error occured during checking if the entity %d exists in %s: %s\n", id, identifier, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the entity %d exists in %s: %s\n", id, identifier, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) GetIdOfEntity(identifier, entityId string) (int, error) {
	const query = "SELECT id FROM entities WHERE identity = $1 AND entity_id = $2"

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

func (hdb *HonuaDatabase) GetEntitiesWithoutRule(identifier string) ([]*models.Entity, error) {

	existRules, err := hdb.ExistRules(identifier)
	if err != nil {
		log.Printf("An error occured during getting all entities without rule of identity = %s: %s\n", identifier, err.Error())
		return nil, err
	}
	if !existRules {
		entities, err := hdb.GetEntities(identifier)
		if err != nil {
			log.Printf("An error occured during getting all entities without rule of identity = %s: %s\n", identifier, err.Error())
			return nil, err
		}
		return entities, nil
	}
	const query = "SELECT * FROM entities WHERE identity=$1 AND id NOT IN (SELECT entity_id FROM rules WHERE identity=$1);"
	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting all entities without rule of identity = %s: %s\n", identifier, err.Error())
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

func (hdb *HonuaDatabase) GetVictronEntities(identifier string) ([]*models.Entity, error) {
	const query = "SELECT * FROM entities WHERE identity = $1 AND is_victron_sensor;"

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
	var sensorType models.SensorType
	var hasNumericState bool
	var rulesEnabled bool

	err := rows.Scan(&id, &identity, &entityID, &name, &isDevice, &allowRules, &hasAttribute, &attribute, &isVictronSensor, &hasNumericState, &rulesEnabled, &sensorType)
	if err != nil {
		return nil, err
	}

	var result *models.Entity = &models.Entity{
		Id:              id,
		IdentityId:      identity,
		Name:            name,
		EntityId:        entityID,
		IsDevice:        isDevice,
		AllowRules:      allowRules,
		HasAttribute:    hasAttribute,
		IsVictronSensor: isVictronSensor,
		SensorType: sensorType,
		HasNumericState: hasNumericState,
		RulesEnabled:    rulesEnabled,
	}

	if hasAttribute && attribute.Valid {
		result.Attribute = attribute.String
	} else {
		result.HasAttribute = false
		result.Attribute = ""
	}
	return result, nil
}
