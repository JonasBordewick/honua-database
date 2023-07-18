package honuadatabase

import (
	"database/sql"

	"github.com/JonasBordewick/honua-database.git/models"
)

func (hdb *HonuaDatabase) AddEntity(entity *models.Entity) error {

}

func (hdb *HonuaDatabase) DeleteEntity(id int) error {

}

func (hdb *HonuaDatabase) ExistEntity(identifier, entityId string) (bool, error) {

}

func (hdb *HonuaDatabase) GetIdOfEntity(identifier, entityId string) (int, error) {

}

func (hdb *HonuaDatabase) GetEntities() ([]*models.Entity, error) {

}

func (hdb *HonuaDatabase) GetEntitiesWhereRulesAreAllowed() ([]*models.Entity, error) {

}

func (hdb *HonuaDatabase) GetEntitiesWithoutAnyRule() ([]*models.Entity, error) {

}

func (hdb *HonuaDatabase) GetVictronEntities() ([]*models.Entity, error) {

}

func (hdb *HonuaDatabase) make_entity(rows *sql.Rows) (*models.Entity, error) {

}
