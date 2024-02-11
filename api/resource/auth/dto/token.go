package authdto

var (
	Active  = "active"
	Refresh = "refresh"
	Access  = "access"
)

type PrivateClaimsForToken struct {
	UserId string
	Role   string
}

type TokenDetailForRedis struct {
	UserId    string
	TokenType string
	Token     string
}
