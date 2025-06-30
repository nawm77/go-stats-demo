package main

import (
	"defig-stats-artifact/pkg/apacheBuildr"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
	"regexp"
	"strconv"
	"strings"
)

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func parseBuildString(s string) (*apacheBuildr.ApacheBuildr, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 5 {
		return nil, fmt.Errorf("некорректный формат строки: ожидается 5 частей, получено %d", len(parts))
	}

	disrtibApacheBuildr := &apacheBuildr.ApacheBuildr{
		GroupId:     parts[0],
		ArtifactId:  parts[1],
		PackageType: apacheBuildr.ValueOf(parts[2]),
		Classifier:  parts[3],
		Version:     parts[4],
	}

	return disrtibApacheBuildr, nil
}

func findSimilarVersionDisrtib(version string, distribs []string) (string, bool) {
	for _, distrib := range distribs {
		distribVersion, err := extractVersionFromDistribLink(distrib)
		if err == nil {
			if matchRegexp, _ := regexp.MatchString(version, distribVersion); distribVersion == version || matchRegexp {
				return distrib, true
			}
		}
	}
	return "", false
}

func extractVersionFromDistribLink(url string) (string, error) {
	parts := strings.Split(url, "/")

	if len(parts) < 2 {
		return "", fmt.Errorf("URL %s слишком короткий", url)
	}

	return parts[len(parts)-2], nil
}

func prepareTableHeader(sheet *xlsx.Sheet) {
	headerRow := sheet.AddRow()
	headerDefaultDistrib := headerRow.AddCell()
	headerDefaultDistrib.Value = "defaultDistribLink"

	headerUnpacked := headerRow.AddCell()
	headerUnpacked.Value = "unpackedDistribLink"

	headerDate := headerRow.AddCell()
	headerDate.Value = "date"
}
