package main

import (
	"fmt"

	honuadatabase "github.com/JonasBordewick/honua-database"
)

func main() {
	hdb := honuadatabase.GetHonuaDatabaseInstance("postgres", "loadscheduler", "192.168.0.138", "5432", "test-honua", "./files")
	exists, err := hdb.ExistIdentity("bordewickbgd")
	if err != nil {
		return
	}
	fmt.Printf("%t", exists)
}
