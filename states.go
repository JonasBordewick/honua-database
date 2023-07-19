package honuadatabase

import (
	"database/sql"
	"log"
	"time"

	"github.com/JonasBordewick/honua-database/models"
)

func (hdb *HonuaDatabase) AddState(state *models.State) error {
	const query = "INSERT INTO states (entity_id, state) VALUES ($1, $2);"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	_, err := hdb.db.Exec(query, state.EntityId, state.State)
	if err != nil {
		log.Printf("An error occured during adding a new state to table states: %s\n", err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) GetState(entityID int) (*models.State, error) {
	const query = "SELECT * FROM states WHERE id = (SELECT MAX(id) FROM states WHERE entity_id = $1);"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, entityID)
	if err != nil {
		rows.Close()
		log.Printf("An error occured during getting the latest state of entity with id = %d: %s\n", entityID, err.Error())
		return nil, err
	}

	var state *models.State

	for rows.Next() {
		state, err = hdb.make_state(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting the latest state of entity with id = %d: %s\n", entityID, err.Error())
			return nil, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) DeleteOldestState(entityID int) error {
	const query = "DELETE FROM states WHERE id = (SELECT MIN(id) FROM states WHERE entity_id = $1);"
	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()
	_, err := hdb.db.Exec(query, entityID)
	if err != nil {
		log.Printf("An error occured during deleting the oldest state of enitity with id = %d: %s\n", entityID, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) GetNumberOfStatesOfEntity(entityID int) (int, error) {
	const query = "SELECT COUNT(*) AS count FROM states WHERE entity_id = $1;"

	hdb.mutex.Lock()
	defer hdb.mutex.Unlock()

	rows, err := hdb.db.Query(query, entityID)
	if err != nil {
		rows.Close()
		log.Printf("An error occured during getting the number of states of entity with id = %d: %s\n", entityID, err.Error())
		return -1, err
	}

	var counter int = -1

	for rows.Next() {
		err = rows.Scan(&counter)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting the number of states of entity with id = %d: %s\n", entityID, err.Error())
			return -1, err
		}
	}

	rows.Close()

	return counter, nil
}

func (hdb *HonuaDatabase) make_state(rows *sql.Rows) (*models.State, error) {
	var id int
	var entityID int
	var state string
	var recordTime *time.Time
	err := rows.Scan(&id, &entityID, &state, &recordTime)
	if err != nil {
		return nil, err
	}

	return &models.State{
		Id:         id,
		EntityId:   entityID,
		State:      state,
		RecordTime: recordTime,
	}, nil
}
