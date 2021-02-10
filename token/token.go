package token

//Token interface to read a token
type Token interface {
	ReadTokenFile() (string, error)
	DeleteTokenFile() (string, error)
}
