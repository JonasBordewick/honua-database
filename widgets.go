package honuadatabase

import "github.com/JonasBordewick/honua-database/models"

func (hdb *HonuaDatabase) AddWidget(configID int, widget *models.Widget) error {
	const query = "INSERT INTO widgets(config_id) VALUES($1) RETURNING id;"
	var id int
	err := hdb.db.QueryRow(query, configID).Scan(&id)
	if err != nil {
		return err
	}
	for contentKey, contentValue := range widget.Contents {
		err = hdb.AddContent(id, contentKey, contentValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hdb *HonuaDatabase) GetWidgets(configID int) ([]*models.Widget, error) {
	const query = "SELECT id FROM widgets WHERE config_id=$1 ORDER BY id;"

	rows, err := hdb.db.Query(query, configID)
	if err != nil {
		return nil, err
	}

	var result []*models.Widget = []*models.Widget{}

	for rows.Next() {
		var widgetID int
		err := rows.Scan(&widgetID)
		if err != nil {
			return nil, err
		}

		contents, err := hdb.GetContents(widgetID)
		if err != nil {
			return nil, err
		}

		result = append(result, &models.Widget{Contents: contents})
	}

	rows.Close()

	return result, nil
}
