package honuadatabase

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/JonasBordewick/honua-database/models"
)

var test_instance = GetHonuaDatabaseInstance("postgres", "loadscheduler", "192.168.0.138", "5432", "test-honua", "./files")

var test_vs_entity = &models.Entity{
	IdentityId:      "testidentifier",
	EntityId:        "test.test",
	Name:            "test",
	IsDevice:        false,
	AllowRules:      false,
	HasAttribute:    false,
	Attribute:       "",
	IsVictronSensor: true,
	HasNumericState: true,
}

func TestExistIdentityBeforeAdding(t *testing.T) {
	exists, err := test_instance.ExistIdentity("testidentifier")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if exists {
		t.Error("Identity already exists, should not exists")
	}
}

func TestAddIdentityAndExistIdentity(t *testing.T) {
	err := test_instance.AddIdentity(&models.Identity{Id: "testidentifier", Name: "Test"})
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	exists, err := test_instance.ExistIdentity("testidentifier")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if !exists {
		t.Error("Identity does not exist, should exist")
	}
}

func TestGetIdentity(t *testing.T) {
	id, err := test_instance.GetIdentity("testidentifier")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if id.Id != "testidentifier" || id.Name != "Test" {
		if err != nil {
			t.Errorf("FAILED: got identity (%s, %s), want (testidentifier, Test)", id.Id, id.Name)
		}
	}
}

func TestGetIdentities(t *testing.T) {
	ids, err := test_instance.GetIdentities()
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if len(ids) != 1 {
		t.Errorf("FAILED: got identities length %d want 1", len(ids))
	}
}

func TestDeleteIdentity(t *testing.T) {
	err := test_instance.DeleteIdentity("testidentifier")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
}

func TestAddEntitiyBeforeAddingIdentity(t *testing.T) {
	err := test_instance.AddEntity(test_vs_entity)
	if err == nil {
		t.Error("FAILED: got no error, but expected one because Identity should not exist")
	}
}

func TestExistEntityBeforeAdding(t *testing.T) {
	exist, err := test_instance.ExistEntity("testidentifier", "test.test")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if exist {
		t.Error("FAILED: Entity exists, it should not be exist")
	}
}

func TestAddEntity(t *testing.T) {
	err := test_instance.AddIdentity(&models.Identity{Id: "testidentifier", Name: "Test"})
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	err = test_instance.AddEntity(test_vs_entity)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	exist, err := test_instance.ExistEntity("testidentifier", "test.test")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if !exist {
		t.Error("FAILED: Entity does not exist but it should be there")
	}
}

func TestGetEntities(t *testing.T) {
	entities, err := test_instance.GetEntities("testidentifier")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if len(entities) != 1 {
		t.Errorf("FAILED: got entities length %d want 1", len(entities))
	}
	entity := entities[0]
	if !entity.Equals(test_vs_entity) {
		t.Errorf("Failed: The Entity in the list is not equal to the inserted, got %+v want %+v", entity, test_vs_entity)
	}
}

func TestEditEntity(t *testing.T) {
	tmp := test_vs_entity
	tmp.Name = "Edited"
	err := test_instance.EditEntity("testidentifier", test_vs_entity)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
}

func TestGetEntitiesWhereRulesAllowed(t *testing.T) {
	entities, err := test_instance.GetEntitiesWhereRulesAreAllowed("testidentifier")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if len(entities) != 0 {
		t.Errorf("FAILED: got entities length %d want 0", len(entities))
	}
}

func TestGetVictronEntities(t *testing.T) {
	entities, err := test_instance.GetVictronEntities("testidentifier")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if len(entities) != 1 {
		t.Errorf("FAILED: got entities length %d want 1", len(entities))
	}
}

func TestDeleteEntity(t *testing.T) {
	id, err := test_instance.GetIdOfEntity("testidentifier", "test.test")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	err = test_instance.DeleteEntity(id)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
}

func TestAddState(t *testing.T) {
	err := test_instance.AddEntity(test_vs_entity)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	id, err := test_instance.GetIdOfEntity("testidentifier", "test.test")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	err = test_instance.AddState("testidentifier", &models.State{EntityId: id, State: "69.69"})
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
}

func TestGetState(t *testing.T) {
	id, err := test_instance.GetIdOfEntity("testidentifier", "test.test")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	state, err := test_instance.GetState(id)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if state.State != "69.69" {
		t.Errorf("FAILED: got %s wanted 69.69", state.State)
	}
}

func TestGetNumbersOfStates(t *testing.T) {
	id, err := test_instance.GetIdOfEntity("testidentifier", "test.test")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	number, err := test_instance.GetNumberOfStatesOfEntity(id)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if number != 1 {
		t.Errorf("FAILED: got %d wanted 1", number)
	}
}

func TestDeleteOldestState(t *testing.T) {
	id, err := test_instance.GetIdOfEntity("testidentifier", "test.test")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	test_instance.DeleteOldestState(id)
	number, err := test_instance.GetNumberOfStatesOfEntity(id)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if number != 0 {
		t.Errorf("FAILED: got %d wanted 0", number)
	}
}

func TestHassService(t *testing.T) {

	var identifier = randSeq(10)
	var domain = fmt.Sprintf("%s.domain", identifier)

	// setup
	t.Run("Create Identity", func(t *testing.T) {
		err := test_instance.AddIdentity(&models.Identity{Id: identifier, Name: "Zufällige Identität"})
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistIdentity(identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Identity does not exist, should exist")
		}
	})

	t.Run("Exist Hass Service Before Adding a Hass Service", func(t *testing.T) {
		exist, err := test_instance.ExistsHassService(identifier, domain)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if exist {
			t.Error("FAILED: expected that HassService does not exists")
		}
	})

	t.Run("Add hass Service", func(t *testing.T) {
		err := test_instance.AddHassService(&models.HassService{Domain: domain, Name: "Zufällige Domain", Enabled: true}, identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
	})

	t.Run("Exist Hass Service After Adding", func(t *testing.T) {
		exist, err := test_instance.ExistsHassService(identifier, domain)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exist {
			t.Error("FAILED: expected that HassService exists but it does not exist")
		}
	})

	t.Run("GetID", func(t *testing.T) {
		_, err := test_instance.GetIDofHassService(identifier, domain)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
	})

	t.Run("ToggleService", func(t *testing.T) {
		err := test_instance.ToggleHassService(identifier, domain)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
	})

	t.Run("DeleteService", func(t *testing.T) {
		err := test_instance.DeleteHassService(identifier, domain)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
	})

	// Clean
	err := test_instance.DeleteIdentity(identifier)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
}

func TestAllowService(t *testing.T) {
	var identifier = randSeq(10)
	var domain = fmt.Sprintf("%s.domain", identifier)
	var entity = &models.Entity{
		IdentityId:      identifier,
		EntityId:        "test.entity",
		Name:            "Test Entity",
		IsDevice:        true,
		AllowRules:      true,
		HasAttribute:    false,
		Attribute:       "",
		IsVictronSensor: false,
		HasNumericState: false,
	}

	// setup
	t.Run("Create Identity", func(t *testing.T) {
		err := test_instance.AddIdentity(&models.Identity{Id: identifier, Name: "Zufällige Identität"})
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistIdentity(identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Identity does not exist, should exist")
		}
	})

	t.Run("Create Entity", func(t *testing.T) {
		err := test_instance.AddEntity(entity)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistEntity(identifier, entity.EntityId)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Entity does not exist, should exist")
		}
	})

	t.Run("Create HassService", func(t *testing.T) {
		err := test_instance.AddHassService(&models.HassService{
			Domain: domain,
			Name:   "Zufälliger Service",
		}, identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistsHassService(identifier, domain)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Hass Serivce does not exist, should exist")
		}
	})

	err := test_instance.AllowService(identifier, domain, entity.EntityId)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	allowed, err := test_instance.IsServiceAllowed(identifier, domain, entity.EntityId)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if !allowed {
		t.Errorf("Hass Serivce %s is not allowed for %s in %s", domain, entity.EntityId, identifier)
	}
	err = test_instance.DisallowService(identifier, domain, entity.EntityId)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}

	t.Run("Clean", func(t *testing.T) {
		err := test_instance.DeleteIdentity(identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
	})
}

func TestAllowSensor(t *testing.T) {
	var identifier = randSeq(10)
	var entity = &models.Entity{
		IdentityId:      identifier,
		EntityId:        "test.entity",
		Name:            "Test Entity",
		IsDevice:        true,
		AllowRules:      true,
		HasAttribute:    false,
		Attribute:       "",
		IsVictronSensor: false,
		HasNumericState: false,
	}

	var sensor = &models.Entity{
		IdentityId:      identifier,
		EntityId:        "test.sensor",
		Name:            "Test Sensor",
		IsDevice:        true,
		AllowRules:      true,
		HasAttribute:    false,
		Attribute:       "",
		IsVictronSensor: false,
		HasNumericState: false,
	}

	// setup
	t.Run("Create Identity", func(t *testing.T) {
		err := test_instance.AddIdentity(&models.Identity{Id: identifier, Name: "Zufällige Identität"})
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistIdentity(identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Identity does not exist, should exist")
		}
	})

	t.Run("Create Entity", func(t *testing.T) {
		err := test_instance.AddEntity(entity)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistEntity(identifier, entity.EntityId)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Entity does not exist, should exist")
		}
	})

	t.Run("Create Entity", func(t *testing.T) {
		err := test_instance.AddEntity(sensor)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistEntity(identifier, sensor.EntityId)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Entity does not exist, should exist")
		}
	})

	err := test_instance.AllowSensor(identifier, sensor.EntityId, entity.EntityId)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	allowed, err := test_instance.IsSensorAllowed(identifier, sensor.EntityId, entity.EntityId)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
	if !allowed {
		t.Errorf("Hass Serivce %s is not allowed for %s in %s", sensor.EntityId, entity.EntityId, identifier)
	}
	err = test_instance.DisallowSensor(identifier, sensor.EntityId, entity.EntityId)
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}

	t.Run("Clean", func(t *testing.T) {
		err := test_instance.DeleteIdentity(identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
	})
}

func TestConditions(t *testing.T) {
	var identifier = randSeq(10)

	var sensor_a = &models.Entity{
		IdentityId:      identifier,
		EntityId:        "test.sensor",
		Name:            "Test Sensor",
		IsDevice:        false,
		AllowRules:      false,
		HasAttribute:    false,
		Attribute:       "",
		IsVictronSensor: true,
		HasNumericState: true,
	}

	var sensor_b = &models.Entity{
		IdentityId:      identifier,
		EntityId:        "test.sensor2",
		Name:            "Test Sensor2",
		IsDevice:        false,
		AllowRules:      false,
		HasAttribute:    false,
		Attribute:       "",
		IsVictronSensor: false,
		HasNumericState: false,
	}

	t.Run("Create Identity", func(t *testing.T) {
		err := test_instance.AddIdentity(&models.Identity{Id: identifier, Name: "Zufällige Identität"})
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistIdentity(identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Identity does not exist, should exist")
		}
	})

	t.Run("Create Entity", func(t *testing.T) {
		err := test_instance.AddEntity(sensor_a)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistEntity(identifier, sensor_a.EntityId)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Entity does not exist, should exist")
		}
	})

	t.Run("Create Entity", func(t *testing.T) {
		err := test_instance.AddEntity(sensor_b)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		exists, err := test_instance.ExistEntity(identifier, sensor_b.EntityId)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
		if !exists {
			t.Error("Entity does not exist, should exist")
		}
	})

	t.Run("AddConditions", func(t *testing.T) {
		condition := &models.Condition{
			Type: models.AND,
			SubConditions: []*models.Condition{
				{
					Type: models.NUMERICSTATE,
					Sensor: sensor_a,
					Above: &models.ConditionValue{
						Valid: true,
						Value: 10,
					},
				},
				{
					Type: models.STATE,
					Sensor: sensor_b,
					ComparisonState: "on",
				},
				{
					Type: models.TIME,
					After: "10:00",
					Before: "11:00",
				},
			},	
		}
		id, err := test_instance.AddCondition(identifier, condition)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}

		exists, err := test_instance.ExistCondition(id, identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}

		if !exists {
			t.Error("[FAILED] Expected that Condition Exists")
		}

		c, err := test_instance.GetCondition(id, identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}

		c.Type = models.NAND
		c.SubConditions[0].Above.Value = 12
		c.SubConditions[1].ComparisonState = "off"
		c.SubConditions[2].After = "10:30"

		err = test_instance.EditCondition(identifier, c)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}

		c, err = test_instance.GetCondition(id, identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}

		log.Printf("%d | %d | %s | %s\n", condition.Type, condition.SubConditions[0].Above.Value, condition.SubConditions[1].ComparisonState, condition.SubConditions[2].After)

		log.Printf("%d | %d | %s | %s\n", c.Type, c.SubConditions[0].Above.Value, c.SubConditions[1].ComparisonState, c.SubConditions[2].After)

		if c.Type != models.NAND || c.SubConditions[0].Above.Value != 12 || c.SubConditions[1].ComparisonState != "off" || c.SubConditions[2].After != "10:30" {
			t.Error("FAILED: edit condition failed")
		}

		err = test_instance.DeleteCondition(id, identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}


	})

	t.Run("Clean", func(t *testing.T) {
		err := test_instance.DeleteIdentity(identifier)
		if err != nil {
			t.Errorf("FAILED: got error %s", err.Error())
		}
	})
}

func TestClean(t *testing.T) {
	err := test_instance.DeleteIdentity("testidentifier")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	rand.NewSource(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
