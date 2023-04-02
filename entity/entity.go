package entity

type RegisterResponseSuccess struct {
	Message string `json:"message"`
}

type RegisterResponseFailure struct {
	Error string `json:"error"`
}
