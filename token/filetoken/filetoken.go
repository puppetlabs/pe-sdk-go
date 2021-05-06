package filetoken

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"

	"github.com/puppetlabs/pe-sdk-go/log"
	"github.com/puppetlabs/pe-sdk-go/token"
)

var fs = afero.NewOsFs()
var afs = &afero.Afero{Fs: fs}

// fileToken struct
type fileToken struct {
	path string
}

// NewFileToken constructs a new filetoken
func NewFileToken(path string) token.Token {
	fileToken := fileToken{path: getPath(path)}
	return &fileToken
}

// Read reads the token from a file and returns a valid token.
// If the token is not valid an error will be returned
func (ft *fileToken) Read() (string, error) {
	return ft.readToken()
}

func writeTempFile(dir string, content string) (string, error) {
	tmpFile, err := afs.TempFile(dir, "temp")
	if err != nil {
		return "", err
	}

	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), err
}

// Write ...
func (ft *fileToken) Write(token string) error {
	err := fs.MkdirAll(filepath.Dir(ft.path), 0700)
	if err != nil {
		return err
	}

	err = fs.Chmod(filepath.Dir(ft.path), 0700)
	if err != nil {
		return err
	}

	tmpFile, err := writeTempFile(filepath.Dir(ft.path), token)
	if err != nil {
		return err
	}

	err = fs.Rename(tmpFile, ft.path)
	if err != nil {
		return err
	}

	return fs.Chmod(ft.path, 0600)
}

//DeleteFile deletes the token file
//If the token is not valid, or file is a symlink an error will be returned
func (ft *fileToken) Delete(force bool) error {
	if force == true {
		return ft.deleteFile()
	}
	_, err := ft.readToken()
	if err != nil {
		return err
	}
	return ft.deleteFile()
}

func (ft *fileToken) deleteFile() error {
	file, err := os.Lstat(ft.path) // FIXME afero
	if err != nil {
		return err
	}
	if file.Mode()&os.ModeSymlink == os.ModeSymlink {
		return errors.New("Cannot delete token because it is a symbolic link instead of a token file:" + ft.path)
	}
	err = fs.Remove(ft.path)
	if err != nil {
		return err
	}
	fmt.Println("Token file deleted")
	return nil
}

func getPath(path string) string {
	if path == "" {
		return defaultPath()
	}
	return path
}

func (ft *fileToken) readToken() (string, error) {
	data, err := afs.ReadFile(ft.path)
	if err != nil {
		return "", err
	}
	token := strings.TrimRight(string(data), "\r\n")
	valid := isValid(token)
	if !valid {
		return "", fmt.Errorf("Malformed token: token does not match expected formats")
	}
	return token, nil
}

func isValid(content string) bool {
	// tokens can be of one of two flavors -- they can be jwt like, or new token like
	// jwt tokens consist of three segments, dot separated, which each section being at
	// least four characters from A-Za-z0-9_-
	// new tokens consist only of characters A-Za-z0-9_-

	jwtExpresion := `([A-Za-z0-9_-]{4,})\.([A-Za-z0-9_-]{4,})\.([A-Za-z0-9_-]{4,})`
	tokenExpr := `^([A-Za-z0-9_-]+)$` //regexp.MatchString searches for any match, so we should specify that it should match the whole string with ^$

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
func defaultPath() string {
	usr, err := user.Current()
	if err != nil {
		//FIXME log level
		log.Debug(err.Error())
		return ""
	}
	return filepath.Join(usr.HomeDir, ".puppetlabs", "token")
}
