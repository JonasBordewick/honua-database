package honuadatabase

import (
	"database/sql"

	"github.com/JonasBordewick/honua-database.git/models"
)

func (hdb *HonuaDatabase) AddState(state *models.State) error {

}

func (hdb *HonuaDatabase) GetState(entityID int) (*models.State, error) {

}

func (hdb *HonuaDatabase) DeleteOldestState(entityID int) error {

}

func (hdb *HonuaDatabase) make_state(rows *sql.Rows) (*models.State, error) {

}
