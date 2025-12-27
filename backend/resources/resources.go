package resources

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/kvloginov/cup-of-team/backend/internal/domain"
)

//go:embed holidays/holidays_ru_12.json
var holidaysRu12Json []byte

//go:embed namedays/namedays_ru.json
var namedaysRuJson []byte

func LoadHolidaysRu12() (domain.Holidays, error) {
	var dayHolidaysList []domain.DayHolidays
	err := json.Unmarshal(holidaysRu12Json, &dayHolidaysList)
	if err != nil {
		return nil, err
	}

	holidays := make(domain.Holidays)
	for _, dayHolidays := range dayHolidaysList {
		holidays[dayHolidays.Date] = dayHolidays
	}

	return holidays, nil
}

func LoadNamedaysRu() (domain.Namedays, error) {
	var namedayList []domain.Nameday
	err := json.Unmarshal(namedaysRuJson, &namedayList)
	if err != nil {
		return nil, err
	}

	namedays := make(domain.Namedays)
	for _, nameday := range namedayList {
		namedays[nameday.Date] = strings.Join(nameday.Names, ", ")
	}

	return namedays, nil
}
