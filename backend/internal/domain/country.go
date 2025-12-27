package domain

type CountryCode string

const (
	CountryCodeRU CountryCode = "RU"
	CountryCodeUS CountryCode = "US"
	CountryCodeWW CountryCode = "WW"
	CountryCodeUK CountryCode = "UK"
	CountryCodeDE CountryCode = "DE"
	CountryCodeFR CountryCode = "FR"
	CountryCodeIT CountryCode = "IT"
	CountryCodeES CountryCode = "ES"
	CountryCodePT CountryCode = "PT"
	CountryCodeNL CountryCode = "NL"
	CountryCodeBE CountryCode = "BE"
	CountryCodeCH CountryCode = "CH"
)

type Country struct {
	Code CountryCode `json:"code"`
	Name string      `json:"name"`
}
