package honuadatabase

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/JonasBordewick/honua-database/models"
)

func (hdb *HonuaDatabase) AddHassService(service *models.HassService, identity string) error {
	const query = `
INSERT INTO hass_services(
	identity, domain, name
) VALUES ($1, $2, $3);
`

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	_, err := hdb.db.Exec(query, identity, service.Domain, service.Name)

	if err != nil {
		log.Printf("An error occured during adding a new Homeassistant Service to table hass_services: %s\n", err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) GetIDofHassService(identity, domain string) (int, error) {
	const query = "SELECT id FROM hass_services WHERE identity=$1 AND domain=$2;"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, identity, domain)
	if err != nil {
		log.Printf("An error occured during getting the id of homeassistant service with identity = %s and domain = %s: %s\n", identity, domain, err.Error())
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting the id of homeassistant service with identity = %s and domain = %s: %s\n", identity, domain, err.Error())
			return -1, err
		}
	}

	rows.Close()

	if id == -1 {
		return -1, fmt.Errorf("no element in database found where identity = %s and domain = %s", identity, domain)
	}

	return id, nil
}

func (hdb *HonuaDatabase) ToggleHassService(identity, domain string) error {
	const query = "UPDATE hass_services SET enabled = NOT enabled WHERE identity=$1 AND domain=$2;"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	_, err := hdb.db.Exec(query, identity, domain)
	if err != nil {
		log.Printf("An error occured during changing the enabled state to the opposite from homeassistant service of identity %s with domain = %s: %s\n", identity, domain, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) DeleteHassService(identity, domain string) error {
	const query = "DELETE FROM hass_services WHERE identity=$1 AND domain=$2;"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	_, err := hdb.db.Exec(query, identity, domain)
	if err != nil {
		log.Printf("An error occured during deleting the homeassistant service of identity %s with domain = %s: %s\n", identity, domain, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) ExistsHassService(identity, domain string) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM hass_services WHERE identity = $1 AND domain = $2) THEN true ELSE false END"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, identity, domain)
	if err != nil {
		log.Printf("An error occured during checking if the entity %s exists in %s: %s\n", identity, domain, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the entity %s exists in %s: %s\n", identity, domain, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) GetAllowedHassServicesOfEntity(identity, entityId string) ([]*models.HassService, error) {
	const query = `
	SELECT services.domain, services.name, services.enabled 
	FROM hass_services as services, allowed_services as a 
	WHERE services.id = a.service_id AND a.entity_id = (SELECT id FROM entities WHERE identity=$1 AND entity_id=$2);
	`

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, identity, entityId)
	if err != nil {
		log.Printf("An error occured during getting all allowed homeassistant services of entity %s in %s: %s\n", entityId, identity, err.Error())
		return nil, err
	}

	var result []*models.HassService = []*models.HassService{}

	for rows.Next() {
		service, err := hdb.make_hass_service(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all allowed homeassistant services of entity %s in %s: %s\n", entityId, identity, err.Error())
			return nil, err
		}
		result = append(result, service)
	}
	rows.Close()

	return result, err
}

func (hdb *HonuaDatabase) make_hass_service(rows *sql.Rows) (*models.HassService, error) {
	var domain string
	var name string
	var enabled bool

	err := rows.Scan(&domain, &name, &enabled)
	if err != nil {
		return nil, err
	}

	var result *models.HassService = &models.HassService{
		Domain:  domain,
		Name:    name,
		Enabled: enabled,
	}

	return result, nil
}
