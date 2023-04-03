package entity

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

