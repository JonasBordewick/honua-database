package honuadatabase

import (
	"database/sql"

	"github.com/JonasBordewick/honua-database.git/models"
)

func (hdb *HonuaDatabase) AddIdentity(identity *models.Identity) error {

}

func (hdb *HonuaDatabase) DeleteIdentity(identifier string) error {

}

func (hdb *HonuaDatabase) ExistIdentity(identifier string) (bool, error) {

}

func (hdb *HonuaDatabase) GetIdentity(identifier string) (*models.Identity, error) {

}

func (hdb *HonuaDatabase) GetIdentities() ([]*models.Identity, error) {

}

func (hdb *HonuaDatabase) make_identity(rows *sql.Rows) (*models.Identity, error) {

}
