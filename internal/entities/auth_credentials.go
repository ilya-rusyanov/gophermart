package entities

type AuthCredentials struct {
	Login    Login  `json:"login"`
	Password string `json:"password"`
}
