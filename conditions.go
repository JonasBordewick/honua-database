package honuadatabase

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/JonasBordewick/honua-database/models"
)

const add_condition_query = `
INSERT INTO conditions(
	id, identity, type, sensor_id, before,
	after, below, above,
	comparison_state, parent_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
`

func (hdb *HonuaDatabase) get_condition_id(identifier string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM conditions WHERE identity = $1) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting id of condition in %s: %s\n", identifier, err.Error())
		return -1, err
	}

	var exist_identity bool = false

	for rows.Next() {
		err = rows.Scan(&exist_identity)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of enconditiontity in %s: %s\n", identifier, err.Error())
			return -1, err
		}
	}

	rows.Close()

	if !exist_identity {
		return 0, nil
	}

	query = "SELECT MAX(id) FROM conditions WHERE identity = $1;"

	rows, err = hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting id of condition in %s: %s\n", identifier, err.Error())
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of condition in %s: %s\n", identifier, err.Error())
			return -1, err
		}
	}
	rows.Close()

	if id == -1 {
		return -1, errors.New("something went wrong during getting id of condition")
	}

	id = id + 1

	return id, nil
}

func (hdb *HonuaDatabase) AddCondition(identity string, condition *models.Condition) (int, error) {

	id, err := hdb.get_condition_id(identity)
	if err != nil {
		log.Printf("An error occured during adding a new condition: %s\n", err.Error())
		return -1, err
	}

	_, err = hdb.db.Exec(add_condition_query, id, identity, condition.Type, sql.NullInt32{}, sql.NullString{}, sql.NullString{}, sql.NullInt32{}, sql.NullInt32{}, sql.NullString{}, sql.NullInt32{})
	if err != nil {
		log.Printf("Error during adding new condition to table: %s\n", err.Error())
		return -1, err
	}

	for _, sub := range condition.SubConditions {
		err = hdb.add_subcondition(identity, sub, id)
		if err != nil {
			log.Printf("Error during adding new condition to table: %s\n", err.Error())
			return -1, err
		}
	}

	return id, nil
}

func (hdb *HonuaDatabase) add_subcondition(identity string, condition *models.Condition, parentID int) error {
	id, err := hdb.get_condition_id(identity)
	if err != nil {
		log.Printf("An error occured during adding a new condition: %s\n", err.Error())
		return err
	}
	log.Printf("Parent: %d || ID: %d\n", parentID, id)
	if condition.Type == models.NUMERICSTATE {
		var below sql.NullInt32 = sql.NullInt32{}
		var above sql.NullInt32 = sql.NullInt32{}

		if condition.Below != nil {
			below = sql.NullInt32{Valid: condition.Below.Valid, Int32: int32(condition.Below.Value)}
		}

		if condition.Above != nil {
			above = sql.NullInt32{Valid: condition.Above.Valid, Int32: int32(condition.Above.Value)}
		}

		_, err = hdb.db.Exec(add_condition_query, id, identity, condition.Type, condition.Sensor.Id, sql.NullString{}, sql.NullString{}, below, above, sql.NullString{}, parentID)
		if err != nil {
			log.Printf("Error during adding new condition to table: %s\n", err.Error())
		}
		return err
	} else if condition.Type == models.STATE {
		_, err := hdb.db.Exec(add_condition_query, id, identity, condition.Type, condition.Sensor.Id, sql.NullString{}, sql.NullString{}, sql.NullInt32{}, sql.NullInt32{}, condition.ComparisonState, parentID)
		if err != nil {
			log.Printf("Error during adding new condition to table: %s\n", err.Error())
		}
		return err
	} else if condition.Type == models.TIME {
		var before sql.NullString = sql.NullString{
			Valid:  len(condition.Before) > 0,
			String: condition.Before,
		}

		var after sql.NullString = sql.NullString{
			Valid:  len(condition.After) > 0,
			String: condition.After,
		}
		_, err := hdb.db.Exec(add_condition_query, id, identity, condition.Type, sql.NullInt32{}, before, after, sql.NullInt32{}, sql.NullInt32{}, sql.NullString{}, parentID)
		if err != nil {
			log.Printf("Error during adding new condition to table: %s\n", err.Error())
		}
		return err
	}

	log.Printf("Error during adding new condition to table: ConditionType %d not supported.\n", condition.Type)
	return fmt.Errorf("error during adding new condition to table: ConditionType %d not supported", condition.Type)
}

func (hdb *HonuaDatabase) DeleteCondition(conditionID int, identity string) error {
	const query = "DELETE FROM conditions WHERE id=$1 AND identity=$2;"

	_, err := hdb.db.Exec(query, conditionID, identity)
	if err != nil {
		log.Printf("An error occured during deleting the condition with id = %d of identity %s: %s\n", conditionID, identity, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) EditCondition(identity string, condition *models.Condition) error {
	exist, err := hdb.ExistCondition(condition.Id, identity)

	if err != nil {
		log.Printf("An error occured during editity condition: %s\n", err.Error())
		return err
	}

	if !exist {
		log.Println("This Condition does not exist")
		return fmt.Errorf("the Condition with a id %d does not exist in %s", condition.Id, identity)
	}

	if condition.Type < models.NUMERICSTATE {
		hasParent, err := hdb.has_no_parent(condition.Id, identity)
		if err != nil {
			log.Printf("An error occured during editity condition: %s\n", err.Error())
			return err
		}
		if !hasParent {
			log.Printf("This Condition (%d, %s) has a parent, the condition type of %d is not valid.\n", condition.Id, identity, condition.Type)
			return fmt.Errorf("this Condition (%d, %s) has a parent, the condition type of %d is not valid", condition.Id, identity, condition.Type)
		}
		query := "UPDATE conditions SET type=$1 WHERE id=$2 AND identity=$3"

		_, err = hdb.db.Exec(query, condition.Type, condition.Id, identity)
		if err != nil {
			log.Printf("An error occured during editity condition: %s\n", err.Error())
			return err
		}

		for _, c := range condition.SubConditions {
			err = hdb.EditCondition(identity, c)
			if err != nil {
				log.Printf("An error occured during editity condition: %s\n", err.Error())
				return err
			}
		}

		return nil
		
	} else {
		hasParent, err := hdb.has_no_parent(condition.Id, identity)
		if err != nil {
			log.Printf("An error occured during editity condition: %s\n", err.Error())
			return err
		}
		if hasParent {
			log.Printf("This Condition (%d, %s) hasn't a parent, the condition type of %d is not valid.\n", condition.Id, identity, condition.Type)
			return fmt.Errorf("this Condition (%d, %s) hasn't a parent, the condition type of %d is not valid", condition.Id, identity, condition.Type)
		}

		if condition.Type == models.NUMERICSTATE {
			query := "UPDATE conditions SET type=$1, sensor_id=$2, below=$3, above=$4 WHERE id=$5 AND identity=$6"
			var below sql.NullInt32 = sql.NullInt32{}
			var above sql.NullInt32 = sql.NullInt32{}

			if condition.Below != nil {
				below = sql.NullInt32{Valid: condition.Below.Valid, Int32: int32(condition.Below.Value)}
			}

			if condition.Above != nil {
				above = sql.NullInt32{Valid: condition.Above.Valid, Int32: int32(condition.Above.Value)}
			}
			_, err = hdb.db.Exec(query, condition.Type, condition.Sensor.Id, below, above, condition.Id, identity)
			if err != nil {
				log.Printf("An error occured during editity condition: %s\n", err.Error())
			}
			return err
		} else if condition.Type == models.STATE {
			query := "UPDATE conditions SET type=$1, sensor_id=$2, comparison_state=$3 WHERE id=$4 AND identity=$5"
			_, err = hdb.db.Exec(query, condition.Type, condition.Sensor.Id, condition.ComparisonState, condition.Id, identity)
			if err != nil {
				log.Printf("An error occured during editity condition: %s\n", err.Error())
			}
			return err
		} else if condition.Type == models.TIME {
			var before sql.NullString = sql.NullString{
				Valid:  len(condition.Before) > 0,
				String: condition.Before,
			}

			var after sql.NullString = sql.NullString{
				Valid:  len(condition.After) > 0,
				String: condition.After,
			}

			query := "UPDATE conditions SET type=$1, after=$2, before=$3 WHERE id=$4 AND identity=$5"

			_, err := hdb.db.Exec(query, condition.Type, after, before, condition.Id, identity)
			if err != nil {
				log.Printf("An error occured during editity condition: %s\n", err.Error())
			}
			return err
		}
		return fmt.Errorf("the condition type %d is not supported", condition.Type)
	}
}

func (hdb *HonuaDatabase) ExistCondition(conditionID int, identity string) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM conditions WHERE identity=$1 AND id = $2) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identity, conditionID)
	if err != nil {
		log.Printf("An error occured during checking if the condition with id %d exists: %s\n", conditionID, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the condition with id %d exists: %s\n", conditionID, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) GetCondition(conditionID int, identity string) (*models.Condition, error) {

	exist, err := hdb.ExistCondition(conditionID, identity)
	if err != nil {
		log.Printf("An error occured during getting the condition with id %d: %s\n", conditionID, err.Error())
		return nil, err
	}

	if !exist {
		log.Printf("the condition with id = %d does not exist!\n", conditionID)
		return nil, fmt.Errorf("the condition with id = %d does not exist", conditionID)
	}

	const query = "SELECT * FROM conditions WHERE identity =$1 AND id=$2;"

	rows, err := hdb.db.Query(query, identity, conditionID)
	if err != nil {
		log.Printf("An error occured during getting the condition with id %d: %s\n", conditionID, err.Error())
		return nil, err
	}

	var result *models.Condition

	for rows.Next() {
		condition, err := hdb.make_condition(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting the condition with id %d: %s\n", conditionID, err.Error())
			return nil, err
		}

		result = condition
	}

	rows.Close()

	return result, nil
}

func (hdb *HonuaDatabase) get_subconditions(parentID int) ([]*models.Condition, error) {
	const query = "SELECT * FROM conditions WHERE parent_id=$1;"

	rows, err := hdb.db.Query(query, parentID)
	if err != nil {
		log.Printf("An error occured during getting all subconditions of condition with id %d: %s\n", parentID, err.Error())
		return nil, err
	}

	var result []*models.Condition = []*models.Condition{}

	for rows.Next() {
		condition, err := hdb.make_condition(rows)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all subconditions of condition with id %d: %s\n", parentID, err.Error())
			return nil, err
		}

		result = append(result, condition)
	}

	rows.Close()

	return result, nil
}

func (hdb HonuaDatabase) has_no_parent(conditionID int, identity string) (bool, error) {
	const query = "SELECT parent_id FROM conditions WHERE identity = $1 AND id = $2"
	rows, err := hdb.db.Query(query, identity, conditionID)
	if err != nil {
		log.Printf("An error occured during checking if the codntion %d has a parent in %s: %s\n", conditionID, identity, err.Error())
		return false, err
	}

	var state bool = false
	for rows.Next() {
		var pid sql.NullInt32
		err = rows.Scan(&pid)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the codntion %d has a parent in %s: %s\n", conditionID, identity, err.Error())
			return false, err
		}
		state = !pid.Valid
	}

	return state, nil
}

func (hdb *HonuaDatabase) make_condition(rows *sql.Rows) (*models.Condition, error) {
	var identity string
	var id int
	var conditionType models.ConditionType
	var sensorID sql.NullInt32
	var before sql.NullString
	var after sql.NullString
	var below sql.NullInt32
	var above sql.NullInt32
	var comparisonState sql.NullString
	var parentID sql.NullInt32

	err := rows.Scan(&id, &identity, &conditionType, &sensorID, &before, &after, &below, &above, &comparisonState, &parentID)
	if err != nil {
		return nil, err
	}

	if conditionType < models.NUMERICSTATE {
		sub, err := hdb.get_subconditions(id)
		if err != nil {
			return nil, err
		}
		return &models.Condition{
			Id:            id,
			Type:          conditionType,
			SubConditions: sub,
		}, nil
	}

	if conditionType == models.NUMERICSTATE {
		if !sensorID.Valid || !(above.Valid || below.Valid) {
			return nil, errors.New("numeric_state condition is not valid")
		}

		sensor, err := hdb.GetEntity(identity, int(sensorID.Int32))
		if err != nil {
			return nil, err
		}

		// Assertion: Numeric State is Valid
		return &models.Condition{
			Id:     id,
			Type:   conditionType,
			Sensor: sensor,
			Above:  &models.ConditionValue{Valid: above.Valid, Value: int(above.Int32)},
			Below:  &models.ConditionValue{Valid: below.Valid, Value: int(below.Int32)},
		}, nil
	} else if conditionType == models.STATE {
		if !sensorID.Valid || !comparisonState.Valid {
			return nil, errors.New("state condition is not valid")
		}

		sensor, err := hdb.GetEntity(identity, int(sensorID.Int32))
		if err != nil {
			return nil, err
		}

		return &models.Condition{
			Id:              id,
			Type:            conditionType,
			Sensor:          sensor,
			ComparisonState: comparisonState.String,
		}, nil
	} else if conditionType == models.TIME {
		if !(after.Valid || before.Valid) {
			return nil, errors.New("time condition is not valid")
		}
		// Assertion: time condition is valid
		return &models.Condition{
			Id:     id,
			Type:   conditionType,
			After:  after.String,
			Before: before.String,
		}, nil
	}

	return nil, fmt.Errorf("condition type %d not supported", conditionType)
}
