package honuadatabase

import (
	"testing"

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
	err = test_instance.AddState(&models.State{EntityId: id, State: "69.69"})
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
}

func TestGetState (t *testing.T) {
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

func TestGetNumbersOfStates (t *testing.T) {
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


func TestDeleteOldestState (t *testing.T) {
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

func TestClean(t *testing.T) {
	err := test_instance.DeleteIdentity("testidentifier")
	if err != nil {
		t.Errorf("FAILED: got error %s", err.Error())
	}
}
