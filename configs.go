package honuadatabase

import "github.com/JonasBordewick/honua-database/models"

func (hdb *HonuaDatabase) AddConfig(identity string, config *models.Config) error {
	const query = "INSERT INTO configs(identity) VALUES($1) RETURNING id;"
	var id int
	err := hdb.db.QueryRow(query, identity).Scan(&id)
	if err != nil {
		return err
	}
	for _, widget := range config.Widgets {
		err = hdb.AddWidget(id, widget)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hdb *HonuaDatabase) DeleteConfig(identity string) error {
	const query = "DELETE FROM configs WHERE identity=$1;"
	_, err := hdb.db.Exec(query, identity)
	return err
}

func (hdb *HonuaDatabase) EditConfig(identity string, config *models.Config) error {
	err := hdb.DeleteConfig(identity)
	if err != nil {
		return err
	}
	err = hdb.AddConfig(identity, config)
	return err
}

func (hdb *HonuaDatabase) GetConfig(identity string) (*models.Config, error) {
	const query = "SELECT id FROM configs WHERE identity=$1;"
	var id int
	err := hdb.db.QueryRow(query, identity).Scan(&id)
	if err != nil {
		return nil, err
	}
	widgets, err := hdb.GetWidgets(id)
	if err != nil {
		return nil, err
	}
	return &models.Config{
		Widgets: widgets,
	}, nil
}

func (hdb *HonuaDatabase) ExistConfig(identity string) (bool, error) {
	const query = "SELECT CASE WHEN EXISTS ( SELECT * FROM configs WHERE identity = $1) THEN true ELSE false END;"
	rows, err := hdb.db.Query(query, identity)
	if err != nil {
		return false, err
	}

	var state bool = false

	for rows.Next() {
		err = rows.Scan(&state)
		if err != nil {
			rows.Close()
			return false, err
		}
	}

	rows.Close()

	return state, nil
}
