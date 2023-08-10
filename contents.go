package honuadatabase

func (hdb *HonuaDatabase) AddContent(widgetID int, contentKey, contentValue string) error {
	const query = "INSERT INTO contents(widget_id, content_key, content_value) VALUES($1, $2, $3)"
	_, err := hdb.db.Exec(query, widgetID, contentKey, contentValue)
	return err
}

func (hdb *HonuaDatabase) GetContents(widgetID int) (map[string] string, error) {
	const query = "SELECT content_key, content_value FROM contents WHERE widget_id=$1 ORDER BY id;"

	rows, err := hdb.db.Query(query, widgetID)
	if err != nil {
		return nil, err
	}

	var result map[string] string = map[string] string{}

	for rows.Next() {
		var key string
		var value string

		err := rows.Scan(&key, &value)
		if err != nil {
			rows.Close()
			return nil, err
		}

		result[key] = value
	}

	rows.Close()

	return result, nil
}
