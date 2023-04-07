package entity

import "time"

type RegisterResponseSuccess struct {
	Message string `json:"message"`
}

type UserResponseFailure struct {
	Error string `json:"error"`
}

type LoginResponseSuccess struct {
	Token string `json:"token"`
}

type AddCountryRequest struct {
	Name string `json:"name"`
}

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CountryName struct {
	Name string `json:"name"`
}

type Percentage struct {
	Value string `json:"value"`
}

type Statistics struct {
	CountryId   int        `json:"country_id"`
	Confirmed   int        `json:"confirmed"`
	Deaths      int        `json:"death"`
	Recovered   int        `json:"recovered"`
	LastUpdated *time.Time `json:"last_updated"`
}

type CovidData struct {
	Confirmed int `json:"Confirmed"`
	Deaths    int `json:"Deaths"`
	Recovered int `json:"Recovered"`
}
