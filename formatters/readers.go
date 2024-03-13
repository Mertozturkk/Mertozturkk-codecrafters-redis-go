package formatters

import (
	"github.com/codecrafters-io/redis-starter-go/models"
	"strconv"
	"strings"
	"time"
)

func ArrayReader(arr string, processCount int) (models.CliData, error) {
	var data []string
	var setTime time.Duration
	var cmd string

	parsedData := strings.Split(arr, "\r\n")

	for i := 0; i < len(parsedData); i++ {
		element := parsedData[i]

		if element == "" || strconv.Itoa(int(element[0])) == "" || element[0] == models.Bulk {
			continue
		}

		if _, ok := models.Commands[element]; ok {
			cmd = element
		}

		if element == models.Px {
			pxvalue, err := PxReader(parsedData[i+2])
			if err != nil {
				return models.CliData{}, err
			}
			setTime = pxvalue
			break
		}

		data = append(data, element)
		if i == 2*processCount {
			break
		}
	}

	return models.NewCliData(cmd, data, setTime), nil
}

func StringParser(s string) (models.CliData, error) {
	var cliData models.CliData
	cliInputType := s[0]

	switch cliInputType {
	case models.Array:
		processCount, err := strconv.Atoi(string(s[1]))
		if err != nil {
			return cliData, err
		}
		cliData, err = ArrayReader(s[2:], processCount)
		if err != nil {
			return cliData, err
		}
	}
	return cliData, nil
}

func PxReader(value string) (time.Duration, error) {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(val) * time.Millisecond, nil
}
