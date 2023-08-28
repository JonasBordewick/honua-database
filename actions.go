package honuadatabase

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/JonasBordewick/honua-database/models"
)


func (hdb *HonuaDatabase) GetActionsOfRule(identifier string, ruleID int) ([]*models.Action, []*models.Action, error) {
	const query = "SELECT * FROM actions WHERE identity=$1 AND rule_id=$2;"


	rows, err := hdb.db.Query(query, identifier, ruleID)
	if err != nil {
		log.Printf("An error occured during getting all actions of rule %d in %s: %s\n", ruleID, identifier, err.Error())
		return nil, nil, err
	}
	
	thenActions := []*models.Action{}
	elseActions := []*models.Action{}

	for rows.Next() {
		var id int
		var identity string
		var aType models.ActionType
		var ruleid int
		var isThenAction bool
		var serviceID sql.NullInt32
		var delayID sql.NullInt32

		err := rows.Scan(&id, &identity, &aType, &ruleid, &isThenAction, &serviceID, &delayID)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all actions of rule %d in %s: %s\n", ruleID, identifier, err.Error())
			return nil, nil, err
		}

		if aType == models.SERVICE {
			service, err := hdb.GetHassService(identifier, int(serviceID.Int32))
			if err != nil {
				rows.Close()
				log.Printf("An error occured during getting all actions of rule %d in %s: %s\n", ruleID, identifier, err.Error())
				return nil, nil, err
			}

			if isThenAction {
				thenActions = append(thenActions, &models.Action{
					Id: id,
					Type: aType,
					Service: service.Domain,
				})
			} else {
				elseActions = append(elseActions, &models.Action{
					Id: id,
					Type: aType,
					Service: service.Domain,
				})
			}
		} else if aType == models.DELAY {
			delay, err := hdb.GetDelay(identifier, int(delayID.Int32))
			if err != nil {
				rows.Close()
				log.Printf("An error occured during getting all actions of rule %d in %s: %s\n", ruleID, identifier, err.Error())
				return nil, nil, err
			}

			if isThenAction {
				thenActions = append(thenActions, &models.Action{
					Id: id,
					Type: aType,
					Delay: delay,
				})
			} else {
				elseActions = append(elseActions, &models.Action{
					Id: id,
					Type: aType,
					Delay: delay,
				})
			}
		}
	}

	rows.Close()

	return thenActions, elseActions, nil
}

func (hdb *HonuaDatabase) AddAction(identifier string, ruleID int, isThenAction bool, action *models.Action) error {
	id, err := hdb.get_action_id(identifier)
	if err != nil {
		log.Printf("An error occured during adding a new action: %s\n", err.Error())
		return err
	}

	if action.Type == models.DELAY {
		delayID, err := hdb.AddDelay(identifier, action.Delay)
		if err != nil {
			log.Printf("An error occured during adding a new action: %s\n", err.Error())
			return err
		}
		query := "INSERT INTO actions(id, identity, type, rule_id, is_then_action, delay_id) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err  = hdb.db.Exec(query, id, identifier, action.Type, ruleID, isThenAction, delayID)
		if err != nil {
			log.Printf("An error occured during adding a new action: %s\n", err.Error())
			return err
		}
		return nil
	} else if action.Type == models.SERVICE {
		serviceID, err := hdb.GetIDofHassService(identifier, action.Service)
		if err != nil {
			log.Printf("An error occured during adding a new action: %s\n", err.Error())
			return err
		}
		query := "INSERT INTO actions(id, identity, type, rule_id, is_then_action, service_id) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err  = hdb.db.Exec(query, id, identifier, action.Type, ruleID, isThenAction, serviceID)
		if err != nil {
			log.Printf("An error occured during adding a new action: %s\n", err.Error())
			return err
		}
		return nil
	}

	return fmt.Errorf("actiontype %d not supported", action.Type)
}

func (hdb *HonuaDatabase) DeleteAction(identifier string, id int) error {

	aType, err := hdb.get_action_type(identifier, id)
	if err != nil {
		log.Printf("An error occured during deleting the action %d in %s: %s\n", id, identifier, err.Error())
		return err
	}

	if aType == models.DELAY {
		dId, err := hdb.get_delay_id_of_action(identifier, id)
		if err != nil {
			log.Printf("An error occured during deleting the action %d in %s: %s\n", id, identifier, err.Error())
			return err
		}

		err = hdb.DeleteDelay(identifier, dId)
		if err != nil {
			log.Printf("An error occured during deleting the action %d in %s: %s\n", id, identifier, err.Error())
			return err
		}
	}

	const query = "DELETE FROM actions WHERE id=$1 AND identity=$2;"
	_, err = hdb.db.Exec(query, id, identifier)
	if err != nil {
		log.Printf("An error occured during deleting the action %d in %s: %s\n", id, identifier, err.Error())
		return err
	}


	return nil
}

func (hdb *HonuaDatabase) ExistAction(identifier string, id int) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM actions WHERE identity=$1 AND id = $2) THEN true ELSE false END"
	rows, err := hdb.db.Query(query, identifier, id)
	if err != nil {
		log.Printf("An error occured during checking if the action with id %d exists: %s\n", id, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the action with id %d exists: %s\n", id, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) get_action_type(identifier string, id int) (models.ActionType, error) {
	const query = "SELECT type FROM actions WHERE id=$1 AND identity=$2;"
	rows, err := hdb.db.Query(query, id, identifier)
	if err != nil {
		log.Printf("An error occured during getting the action type of action %d in %s: %s\n", id, identifier, err.Error())
		return -1, nil
	}

	var result models.ActionType
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting the action type of action %d in %s: %s\n", id, identifier, err.Error())
			return -1, nil
		}
	}
	rows.Close()

	return result, nil
}

func (hdb *HonuaDatabase) get_delay_id_of_action(identifier string, id int) (int, error) {
	const query = "SELECT delay_id FROM actions WHERE id=$1 AND identity=$2;"
	rows, err := hdb.db.Query(query, id, identifier)
	if err != nil {
		log.Printf("An error occured during getting the delay_id of action %d in %s: %s\n", id, identifier, err.Error())
		return -1, nil
	}

	var result sql.NullInt32
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting the delay_id of action %d in %s: %s\n", id, identifier, err.Error())
			return -1, nil
		}
	}
	rows.Close()

	if !result.Valid {
		log.Printf("An error occured during getting the delay_id of action %d in %s.\n", id, identifier)
		return -1, fmt.Errorf("an error occured during getting the delay_id of action %d in %s", id, identifier)
	}

	return int(result.Int32), nil
}

func (hdb *HonuaDatabase) get_action_id(identifier string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM actions WHERE identity = $1) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting id of action in %s: %s\n", identifier, err.Error())
		return -1, err
	}

	var exist_identity bool = false

	for rows.Next() {
		err = rows.Scan(&exist_identity)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of action in %s: %s\n", identifier, err.Error())
			return -1, err
		}
	}

	rows.Close()

	if !exist_identity {
		return 0, nil
	}

	query = "SELECT MAX(id) FROM actions WHERE identity = $1;"

	rows, err = hdb.db.Query(query, identifier)
	if err != nil {
		log.Printf("An error occured during getting id of action in %s: %s\n", identifier, err.Error())
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of action in %s: %s\n", identifier, err.Error())
			return -1, err
		}
	}
	rows.Close()

	if id == -1 {
		return -1, errors.New("something went wrong during getting id of action")
	}

	id = id + 1

	return id, nil
}
