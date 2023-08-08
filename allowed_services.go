package honuadatabase

import (
	"fmt"
	"log"
)

func (hdb *HonuaDatabase) AllowService(identity, domain, entityId string) error {
	exists, err := hdb.ExistsHassService(identity, domain)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the homeassistant service with identity %s and domain %s does not exist", identity, domain)
	}

	exists, err = hdb.ExistEntity(identity, entityId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the entity with identity %s and entityId %s does not exist", identity, entityId)
	}

	sId, err := hdb.GetIDofHassService(identity, domain)
	if err != nil {
		return err
	}

	eId, err := hdb.GetIdOfEntity(identity, entityId)
	if err != nil {
		return err
	}

	const query = "INSERT INTO allowed_services(identity, entity_id, service_id) VALUES ($1, $2, $3)"

	_, err = hdb.db.Exec(query, identity, eId, sId)

	if err != nil {
		log.Printf("An error occured during adding a new Homeassistant Service to table hass_services: %s\n", err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) DisallowService(identity, domain, entityId string) error {
	allowed, err := hdb.IsServiceAllowed(identity, domain, entityId)
	if err != nil {
		return err
	}
	
	if !allowed {
		return fmt.Errorf("homeassistant service %s is not allowed for %s in %s", domain, entityId, identity)
	}

	sId, err := hdb.GetIDofHassService(identity, domain)
	if err != nil {
		return err
	}

	eId, err := hdb.GetIdOfEntity(identity, entityId)
	if err != nil {
		return err
	}
	
	const query = "DELETE FROM allowed_services WHERE identity=$1 AND entity_id=$2 AND service_id=$3;"

	_, err = hdb.db.Exec(query, identity, eId, sId)

	if err != nil {
		log.Printf("An error occured during deleting a Homeassistant Service from table hass_services: %s\n", err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) IsServiceAllowed(identity, domain, entityId string) (bool, error) {
	exists, err := hdb.ExistsHassService(identity, domain)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("the homeassistant service with identity %s and domain %s does not exist", identity, domain)
	}

	exists, err = hdb.ExistEntity(identity, entityId)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("the entity with identity %s and entityId %s does not exist", identity, entityId)
	}

	sId, err := hdb.GetIDofHassService(identity, domain)
	if err != nil {
		return false, err
	}

	eId, err := hdb.GetIdOfEntity(identity, entityId)
	if err != nil {
		return false, err
	}

	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM allowed_services WHERE identity = $1 aND entity_id = $2 AND service_id = $3) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identity, eId, sId)
	if err != nil {
		log.Printf("An error occured during checking if the service %s is allowed for %s in %s: %s\n", domain, entityId, identity, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the service %s is allowed for %s in %s: %s\n", domain, entityId, identity, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}
