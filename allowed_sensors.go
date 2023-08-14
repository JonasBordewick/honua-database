package honuadatabase

import (
	"fmt"
	"log"
)

func (hdb *HonuaDatabase) AllowSensor(identity, deviceId, sensorId string) error {
	dId, err := hdb.GetIdOfEntity(identity, deviceId)
	if err != nil {
		return err
	}

	sId, err := hdb.GetIdOfEntity(identity, sensorId)
	if err != nil {
		return err
	}

	const query = "INSERT INTO allowed_sensors(identity, device_id, sensor_id) VALUES ($1, $2, $3);"

	_, err = hdb.db.Exec(query, identity, dId, sId)

	if err != nil {
		log.Printf("An error occured during allowing the sensor %s for %s: %s\n", deviceId, sensorId, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) DisallowSensor(identity, deviceId, sensorId string) error {
	allowed, err := hdb.IsSensorAllowed(identity, deviceId, sensorId)
	if err != nil {
		return err
	}
	
	if !allowed {
		return fmt.Errorf("sensor %s is not allowed for %s in %s", sensorId, deviceId, identity)
	}

	dId, err := hdb.GetIdOfEntity(identity, deviceId)
	if err != nil {
		return err
	}

	sId, err := hdb.GetIdOfEntity(identity, sensorId)
	if err != nil {
		return err
	}
	
	const query = "DELETE FROM allowed_sensors WHERE identity=$1 AND device_id=$2 AND sensor_id=$3;"

	_, err = hdb.db.Exec(query, identity, dId, sId)

	if err != nil {
		log.Printf("An error occured during deleting from allowed_sensors: %s\n", err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) IsSensorAllowed(identity, deviceId, sensorId string) (bool, error) {
	dId, err := hdb.GetIdOfEntity(identity, deviceId)
	if err != nil {
		return false, err
	}

	sId, err := hdb.GetIdOfEntity(identity, sensorId)
	if err != nil {
		return false, err
	}

	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM allowed_sensors WHERE identity=$1 AND device_id = $2 AND sensor_id = $3) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identity, dId, sId)
	if err != nil {
		log.Printf("An error occured during checking if the sensor %s is allowed for %s in %s: %s\n", sensorId, deviceId, identity, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the sensor %s is allowed for %s in %s: %s\n", sensorId, deviceId, identity, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}
