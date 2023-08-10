package models

type Config struct {
	Widgets []*Widget
}

type Widget struct {
	Contents map[string]string
}