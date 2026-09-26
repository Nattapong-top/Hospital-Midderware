package domain

type CustomClaims struct {
	Username string
	Hospital string
}

type TokenProvider interface {
	GenerateToken(staff *Staff) (string, error)
	ValidateToken(token string) (*CustomClaims, error)
}
