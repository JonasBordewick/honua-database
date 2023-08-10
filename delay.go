package honuadatabase

import (
	"errors"
	"fmt"
	"log"

	"github.com/JonasBordewick/honua-database/models"
)

func (hdb *HonuaDatabase) GetDelay(identifier string, delayID int) (*models.Delay, error) {
	exist, err := hdb.ExistDelay(identifier, delayID)

	if err != nil {
		log.Printf("An error occured during getting delay %d of %s: %s\n", delayID, identifier, err.Error())
		return nil, err
	}

	if !exist {
		log.Printf("The delay %d of %s does not exist.\n", delayID, identifier)
		return nil, fmt.Errorf("the delay %d of %s does not exist", delayID, identifier)
	}

	const query = "SELECT * from delays WHERE id=$1 AND identity=$2;"

	rows, err := hdb.db.Query(query, delayID, identifier)
	if err != nil {
		log.Printf("An error occured during getting delay %d of %s: %s\n", delayID, identifier, err.Error())
		return nil, err
	}

	var result *models.Delay

	for rows.Next() {
		var id int
		var identity string
		var hours int32
		var minutes int32
		var seconds int32
		err := rows.Scan(&id, &identity, &hours, &minutes, &seconds)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting delay %d of %s: %s\n", delayID, identifier, err.Error())
			return nil, err
		}

		result = &models.Delay{
			Id:      id,
			Hours:   hours,
			Minutes: minutes,
			Seconds: seconds,
		}
	}

	rows.Close()

	if result == nil {
		log.Printf("An error occured during getting delay %d of %s: %s\n", delayID, identifier, err.Error())
		return nil, fmt.Errorf("an error occured during getting delay %d of %s: %s", delayID, identifier, err.Error())
	}

	return result, nil
}

func (hdb *HonuaDatabase) AddDelay(identity string, delay *models.Delay) (int, error) {
	const query = "INSERT INTO delays(id, identity, hours, minutes, seconds) VALUES ($1, $2, $3, $4, $5);"
	id, err := hdb.get_delay_id(identity)
	if err != nil {
		log.Printf("An error occured during adding a new delay: %s\n", err.Error())
		return -1, err
	}

	_, err = hdb.db.Exec(query, id, identity, delay.Hours, delay.Minutes, delay.Seconds)
	if err != nil {
		log.Printf("An error occured during adding a new delay: %s\n", err.Error())
		return -1, err
	}
	return id, nil
}

func (hdb *HonuaDatabase) EditDelay(identity string, delay *models.Delay) error {
	exist, err := hdb.ExistDelay(identity, delay.Id)
	if err != nil {
		log.Printf("Error during editing delay %d of %s: %s\n", delay.Id, identity, err.Error())
		return err
	}
	if !exist {
		log.Printf("The delay %d of %s does not exist\n", delay.Id, identity)
		return fmt.Errorf("the delay %d of %s does not exist", delay.Id, identity)
	}
	const query = "UPDATE delays SET hours=$1, minutes=$2, seconds=$3 WHERE id=$4 AND identity=$5;"
	_, err = hdb.db.Exec(query, delay.Hours, delay.Minutes, delay.Seconds, delay.Id, identity)
	if err != nil {
		log.Printf("Error during editing delay %d of %s: %s\n", delay.Id, identity, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) DeleteDelay(identity string, delayID int) error {
	const query = "DELETE FROM delays WHERE id=$1 AND identity=$2;"

	_, err := hdb.db.Exec(query, delayID, identity)
	if err != nil {
		log.Printf("An error occured during deleting the delay with id = %d of identity %s: %s\n", delayID, identity, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) ExistDelay(identifier string, delayID int) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM delays WHERE identity = $1 AND id = $2) THEN true ELSE false END;"
	rows, err := hdb.db.Query(query, identifier, delayID)
	if err != nil {
		log.Printf("An error occured during checking if the delay %d exists in %s: %s\n", delayID, identifier, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the delay %d exists in %s: %s\n", delayID, identifier, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) get_delay_id(identifier string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM delays WHERE identity = $1) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting id of delay in %s: %s\n", identifier, err.Error())
		return -1, err
	}

	var exist_delay bool = false

	for rows.Next() {
		err = rows.Scan(&exist_delay)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of delay in %s: %s\n", identifier, err.Error())
			return -1, err
		}
	}

	rows.Close()

	if !exist_delay {
		return 0, nil
	}

	query = "SELECT MAX(id) FROM delays WHERE identity = $1;"

	rows, err = hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting id of delay in %s: %s\n", identifier, err.Error())
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of delay in %s: %s\n", identifier, err.Error())
			return -1, err
		}
	}
	rows.Close()

	if id == -1 {
		return -1, errors.New("something went wrong during getting id of delay")
	}

	id = id + 1

	return id, nil
}
