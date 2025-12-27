package domain

type Namedays map[string]string

type Nameday struct {
	Date  string   `json:"date"` // in format DDMM
	Names []string `json:"names"`
}
