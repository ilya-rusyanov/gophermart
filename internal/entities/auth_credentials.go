package entities

type Login string

type AuthCredentials struct {
	Login    Login  `json:"login"`
	Password string `json:"password"`
}
