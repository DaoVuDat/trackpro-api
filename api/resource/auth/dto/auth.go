package authdto

type AuthSignUp struct {
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthLoginTemp struct {
	Token string `json:"token"`
}

type AuthLogin struct {
	Type   *string `json:"type,omitempty"`
	Status string  `json:"status,omitempty"`
}
