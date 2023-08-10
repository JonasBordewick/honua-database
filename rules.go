package honuadatabase

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/JonasBordewick/honua-database/models"
)

func (hdb *HonuaDatabase) GetAllRulesOfIdentity(identity string) ([]*models.Rule, error) {
	const query = "SELECT * FROM rules WHERE identity=$1;"

	rows, err := hdb.db.Query(query, identity)
	if err != nil {
		log.Printf("An error occured during getting all rules of identity %s: %s\n", identity, err.Error())
		return nil, err
	}

	var result []*models.Rule = []*models.Rule{}

	for rows.Next() {
		var id int
		var identity string
		var entity_id int
		var ebe bool
		var periodic sql.NullInt32
		var description string
		var cId int
		var enabled bool

		err = rows.Scan(&id, &identity, &entity_id, &ebe, &periodic, &description, &cId, &enabled)

		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all rules of identity %s: %s\n", identity, err.Error())
			return nil, err
		}

		rule := &models.Rule{
			Id: id,
			Enabled: enabled,
			EventBasedEvaluation: ebe,
		}

		if !ebe {
			rule.PeriodicTrigger = models.PeriodicTriggerType(periodic.Int32)
		}


		entity, err := hdb.GetEntity(identity, entity_id)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all rules of identity %s: %s\n", identity, err.Error())
			return nil, err
		}

		rule.Name = fmt.Sprintf("%s -- Regel", entity.Name)
		rule.Target = entity

		tAction, eActions, err := hdb.GetActionsOfRule(identity, id)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting all rules of identity %s: %s\n", identity, err.Error())
			return nil, err
		}

		rule.ThenActions = tAction
		rule.ElseActions = eActions

		result = append(result, rule)
	}

	rows.Close()

	return result, nil
}

func (hdb *HonuaDatabase) AddRule(identity string, rule *models.Rule) error {
	cID, err := hdb.AddCondition(identity, rule.Condition)
	if err != nil {
		log.Printf("An error occured during add rule: %s\n", err.Error())
		return err
	}

	id, err := hdb.get_rule_id(identity)
	if err != nil {
		log.Printf("An error occured during add rule: %s\n", err.Error())
		return err
	}

	var periodic sql.NullInt32 = sql.NullInt32{
		Valid: !rule.EventBasedEvaluation,
		Int32: int32(rule.PeriodicTrigger),
	}

	const query = `INSERT INTO rules(
		id, identity, entity_id, event_based_evaluation,
		 periodic_trigger_type, description, condition_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7);`


	_, err = hdb.db.Exec(query, id, identity, rule.Target.Id, rule.EventBasedEvaluation, periodic, "", cID)
	if err != nil {
		log.Printf("An error occured during add rule: %s\n", err.Error())
		return err
	}

	for _, a := range rule.ThenActions {
		err = hdb.AddAction(identity, id, true, a)
		if err != nil {
			log.Printf("An error occured during add rule: %s\n", err.Error())
			return err
		}
	}

	for _, a := range rule.ElseActions {
		err = hdb.AddAction(identity, id, false, a)
		if err != nil {
			log.Printf("An error occured during add rule: %s\n", err.Error())
			return err
		}
	}


	return err
}

func (hdb *HonuaDatabase) EditRule(identity string, rule *models.Rule) error {
	err := hdb.DeleteRule(identity, rule.Id)
	if err != nil {
		log.Printf("An error occured during deleting rule %d of %s: %s\n", rule.Id, identity, err.Error())
		return err
	}
	err = hdb.AddRule(identity, rule)
	if err != nil {
		log.Printf("An error occured during deleting rule %d of %s: %s\n", rule.Id, identity, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) DeleteRule(identity string, id int) error {
	// * GET ID of Condition & DELETE Condition
	cID, err := hdb.get_condition_id_of_rule(identity, id)
	if err != nil {
		log.Printf("An error occured during deleting the rule %d in %s: %s\n", id, identity, err.Error())
		return err
	}
	// * CONSTRAINT WILL DELETE SUB CONDITIONS + RULE + ACTIONS
	err = hdb.DeleteCondition(cID, identity)
	if err != nil {
		log.Printf("An error occured during deleting the rule %d in %s: %s\n", id, identity, err.Error())
	}
	return err
}

func (hdb *HonuaDatabase) ExistRule(identity string, id int) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM rules WHERE identity=$1 AND id = $2) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identity, id)
	if err != nil {
		log.Printf("An error occured during checking if the rule with id %d exists: %s\n", id, err.Error())
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during checking if the rule with id %d exists: %s\n", id, err.Error())
			return false, err
		}
	}

	rows.Close()

	return state, nil
}

func (hdb *HonuaDatabase) ExistRules(identity string) (bool, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM rules WHERE identity = $1) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identity)
	if err != nil {
		log.Printf("An error occured during getting id of rule in %s: %s\n", identity, err.Error())
		return false, err
	}

	var exist_identity bool = false

	for rows.Next() {
		err = rows.Scan(&exist_identity)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of rule in %s: %s\n", identity, err.Error())
			return false, err
		}
	}

	rows.Close()

	return exist_identity, nil
}

func (hdb *HonuaDatabase) get_condition_id_of_rule(identifier string, id int) (int, error) {
	const query = "SELECT condition_id FROM rules WHERE id=$1 AND identity=$2;"
	rows, err := hdb.db.Query(query, id, identifier)
	if err != nil {
		log.Printf("An error occured during getting the condition_id of rule %d in %s: %s\n", id, identifier, err.Error())
		return -1, nil
	}

	var result sql.NullInt32
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting the condition_id of rule %d in %s: %s\n", id, identifier, err.Error())
			return -1, nil
		}
	}
	rows.Close()

	if !result.Valid {
		log.Printf("An error occured during getting the delay_id of action %d in %s.\n", id, identifier)
		return -1, fmt.Errorf("an error occured during getting the condition_id of rule %d in %s", id, identifier)
	}

	return int(result.Int32), nil
}

func (hdb *HonuaDatabase) get_rule_id(identity string) (int, error) {
	query := "SELECT CASE WHEN EXISTS ( SELECT * FROM rules WHERE identity = $1) THEN true ELSE false END"

	rows, err := hdb.db.Query(query, identity)
	if err != nil {
		log.Printf("An error occured during getting id of rule in %s: %s\n", identity, err.Error())
		return -1, err
	}

	var exist_identity bool = false

	for rows.Next() {
		err = rows.Scan(&exist_identity)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of rule in %s: %s\n", identity, err.Error())
			return -1, err
		}
	}

	rows.Close()

	if !exist_identity {
		return 0, nil
	}

	query = "SELECT MAX(id) FROM rules WHERE identity = $1;"

	rows, err = hdb.db.Query(query, identity)
	if err != nil {
		log.Printf("An error occured during getting id of rule in %s: %s\n", identity, err.Error())
		return -1, err
	}

	var id int = -1

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			rows.Close()
			log.Printf("An error occured during getting id of rule in %s: %s\n", identity, err.Error())
			return -1, err
		}
	}
	rows.Close()

	if id == -1 {
		return -1, errors.New("something went wrong during getting id of rule")
	}

	id = id + 1

	return id, nil
}
