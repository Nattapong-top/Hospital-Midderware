package domain

type TokenProvider interface {
	GenerateToken(staff *Staff) (string, error)
}
