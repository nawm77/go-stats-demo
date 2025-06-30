package main

import (
	"database/sql"
	db "defig-stats-artifact/db"
	"defig-stats-artifact/internal/config"
	"defig-stats-artifact/internal/logger"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strings"
)

const (
	filename = "output.xlsx"
)

func main() {
	config.LoadEnv()

	nexusRegistryPathTemplate := os.Getenv("NEXUS_REGISTRY_PATH_TEMPLATE")
	dbConnectionParams, err := db.PrepareDBParams()
	dbDriverName := os.Getenv("DB_DRIVER")

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		logger.ErrorLogger.Fatalf("Ошибка создания листа: %s", err)
	}

	prepareTableHeader(sheet)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConnectionParams.Host,
		dbConnectionParams.Port,
		dbConnectionParams.Username,
		dbConnectionParams.Password,
		dbConnectionParams.DBName)

	client, err := sql.Open(dbDriverName, psqlInfo)
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.ErrorLogger.Fatal(err)
		}
	}(client)

	err = client.Ping()
	if err != nil {
		panic(err)
	}
	logger.InfoLogger.Printf("Successfully connected to database %s", dbConnectionParams.DBName)

	rows, err := client.Query("SELECT composition, unpacked_distribs, build_date FROM build WHERE composition is not null AND unpacked_distribs IS NOT NULL AND build_date is not null order by build_date desc;")
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var dataMap = make(map[string]string)

	for rows.Next() {
		var composition, unpackedDistribs, buildDate string
		if err := rows.Scan(&composition, &unpackedDistribs, &buildDate); err != nil {
			log.Fatal(err)
		}

		var currentDistribs []string

		for _, artifact := range strings.Split(composition, "\n") {
			cleaned := strings.TrimSpace(artifact)
			if cleaned != "" {
				apacheBuildr, _ := parseBuildString(cleaned)
				groupId := apacheBuildr.GroupId
				groupId = strings.ReplaceAll(groupId, ".", "/")
				artifactId := apacheBuildr.ArtifactId
				artifactId = strings.ReplaceAll(artifactId, ".", "/")
				version := apacheBuildr.Version

				distribLink := fmt.Sprintf(nexusRegistryPathTemplate, groupId, artifactId, version, artifactId, version)

				currentDistribs = append(currentDistribs, distribLink)
			}
		}

		var currentUnpackedDistribs []string

		for _, unpackedDistribLink := range strings.Split(unpackedDistribs, "\n") {
			cleaned := strings.TrimSpace(unpackedDistribLink)
			if cleaned != "" && !isNumber(cleaned) {
				currentUnpackedDistribs = append(currentUnpackedDistribs, cleaned)
			}
		}

		for _, currentDistrib := range currentDistribs {
			currentDisrtibVersion, err := extractVersionFromDistribLink(currentDistrib)
			if err == nil {
				unpackedDistrb, found := findSimilarVersionDisrtib(currentDisrtibVersion, currentUnpackedDistribs)
				if found {
					dataMap[unpackedDistrb] = currentDistrib
					row := sheet.AddRow()
					cellDefault := row.AddCell()
					cellDefault.Value = currentDistrib

					cellUnpacked := row.AddCell()
					cellUnpacked.Value = unpackedDistrb

					cellDate := row.AddCell()
					cellDate.Value = strings.Split(buildDate, "T")[0]
				}
			}
		}
	}

	err = file.Save(filename)
	if err != nil {
		logger.ErrorLogger.Fatalf("Ошибка сохранения файла %s : %s", filename, err)
	}

	logger.InfoLogger.Printf("Файл Excel успешно создан: %s", filename)
}
