package db

import (
	"log"
	"os"
	"strconv"
)

type DataBaseConnectionParams struct {
	Host     string
	Port     int32
	Username string
	Password string
	DBName   string
}

func PrepareDBParams() (*DataBaseConnectionParams, error) {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	username := os.Getenv("USERNAME")
	pass := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &DataBaseConnectionParams{
		Host:     host,
		Port:     int32(port),
		Username: username,
		Password: pass,
		DBName:   dbname,
	}, nil
}
