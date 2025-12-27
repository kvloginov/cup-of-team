package domain

type User struct {
	ID                string      `json:"id"`
	FirstName         string      `json:"first_name"`
	Initials          string      `json:"initials"`
	ParentNames       []string    `json:"parent_names"`       // list of parent names, no more than 2
	GrandParentsNames []string    `json:"grandparents_names"` // list of grandparent names, no more than 4
	Country           CountryCode `json:"country"`
}

type Team struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Users []User `json:"users"`
}
