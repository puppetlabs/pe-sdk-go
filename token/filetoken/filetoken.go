package filetoken

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/puppetlabs/pe-sdk-go/log"
	"github.com/puppetlabs/pe-sdk-go/token"
)

// fileToken struct
type fileToken struct {
	path string
}

// NewFileToken constructs a new filetoken
func NewFileToken(path string) token.Token {
	fileToken := fileToken{path: path}
	return &fileToken
}

// ReadTokenFile reads the token from a file and returns a valid token.
// If the token is not valid an error will be returned
func (ft *fileToken) ReadTokenFile() (string, error) {
	var err error

	tokenPath, err := getPath(&ft.path)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(*tokenPath)
	if err != nil {
		return "", err
	}

	token := strings.TrimRight(string(data), "\r\n")

	valid := isValid(token)
	if !valid {
		err = fmt.Errorf("Token %s is invalid", token)
		return "", err
	}

	return token, nil
}

//DeleteTokenFile deletes the token file
//If the token is not valid an error will be returned
func (ft *fileToken) DeleteTokenFile() (string, error) {
	var err error

	tokenPath, err := getPath(&ft.path)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(*tokenPath)
	if err != nil {
		return "", err
	}

	token := strings.TrimRight(string(data), "\r\n")
	valid := isValid(token)
	if !valid {
		err = fmt.Errorf("Token %s is invalid", token)
		return "", err
	}
	err = os.Remove(*tokenPath)
	if err != nil {
		return "", err
	}

	return token, nil
}

func getPath(path *string) (*string, error) {
	var err error
	tokenPath := path
	if *path == "" {
		tokenPath, err = defaultPath()
		if err != nil {
			return nil, err
		}
	}
	return tokenPath, nil
}

func isValid(content string) bool {

	jwtExpresion := `([A-Za-z0-9_-]{4,})\.([A-Za-z0-9_-]{4,})\.([A-Za-z0-9_-]{4,})`
	tokenExpr := `([A-Za-z0-9_-]+)`

	jwtMatched, _ := regexp.MatchString(jwtExpresion, content)
	if jwtMatched {
		log.Debug("Token is in JWT format")
		return true
	}

	exprMatch, _ := regexp.MatchString(tokenExpr, content)
	if exprMatch {
		log.Debug("Token format is valid")
		return true
	}

	return false
}

//defaultPath returns token default path
func defaultPath() (*string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	result := filepath.Join(usr.HomeDir, ".puppetlabs", "token")
	return &result, nil
}
