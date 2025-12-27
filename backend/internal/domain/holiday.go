package domain

type HolidayType string

const (
	HolidayTypeNational      HolidayType = "national"
	HolidayTypeInternational HolidayType = "international"
	HolidayTypeHistorical    HolidayType = "historical"
	HolidayTypeInteresting   HolidayType = "interesting"
)

type Holidays map[string]DayHolidays // in format DDMM -> DayHolidays

type DayHolidays struct {
	Date  string    `json:"date"` // in format DDMM
	Dates []Holiday `json:"dates"`
}

type Holiday struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        HolidayType   `json:"type"`
	Countries   []CountryCode `json:"countries"`
	Stats       Stats         `json:"stats"`
}
