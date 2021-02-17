package token

//Token interface to read a token
type Token interface {
	Read() (string, error)
	Write(string) error
	Delete(bool) error
}
