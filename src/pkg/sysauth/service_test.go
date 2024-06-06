package sysauth

import (
	"encoding/json"

	"github.com/IceBearAI/aigc/tests"
	"github.com/go-kit/log"
)

var logger log.Logger

func initSvc() Service {
	apiSvc, err := tests.Init()
	if err != nil {
		panic(err)
	}

	return New(logger, "", tests.Store, apiSvc)
}

func getCurrentWeather(location string, unit string) (string, error) {
	weatherInfo := map[string]interface{}{
		"location":    location,
		"temperature": "72",
		"unit":        unit,
		"forecast":    []string{"sunny", "windy"},
	}
	b, err := json.Marshal(weatherInfo)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
