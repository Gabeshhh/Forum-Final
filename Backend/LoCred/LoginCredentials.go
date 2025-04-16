package LoCred

// Login credential
type LoginCredentials struct {
	Email    string `form:"Email"`
	Password string `form:"Password"`
}
