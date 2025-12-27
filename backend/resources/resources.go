package resources

import (
	_ "embed"
	"encoding/json"

	"github.com/kvloginov/cup-of-team/backend/internal/domain"
)

//go:embed holidays/holidays_ru_12.json
var holidaysRu12Json []byte

//go:embed namedays/namedays_ru.json
var namedaysRuJson []byte

func LoadHolidaysRu12() (domain.Holidays, error) {
	var holidays domain.Holidays
	err := json.Unmarshal(holidaysRu12Json, &holidays)
	if err != nil {
		return nil, err
	}
	return holidays, nil
}

func LoadNamedaysRu() (domain.Namedays, error) {
	var namedays domain.Namedays
	err := json.Unmarshal(namedaysRuJson, &namedays)
	if err != nil {
		return nil, err
	}
	return namedays, nil
}
