package user

type RegisterRequestData struct {
	Email    string
	Name     string
	Surname  string
	Password string
}

type RegisterResponseData struct {
	AuthToken string
	User      UserModel
}

type UpdateRequestData struct {
	Name    string
	Surname string
}

type LoginRequestData struct {
	Email    string
	Password string
}

type LoginResponseData struct {
	AuthToken string
	User      UserModel
}
