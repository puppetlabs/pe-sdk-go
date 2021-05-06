package filetoken

import (
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"

	"errors"
	"os"

	"github.com/spf13/afero"

	match "github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/testing"
	"github.com/puppetlabs/pe-sdk-go/token/testdata"
	"github.com/stretchr/testify/assert"
)

func TestReadOK(t *testing.T) {
	//arrange
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")
	fileToken := NewFileToken(path)

	//act
	token, err := fileToken.Read()

	//assert
	test.Equal(nil, err)
	test.Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5cmytoken", token)
}

func TestIsValidJWTOK(t *testing.T) {
	//arrange
	test := assert.New(t)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dnZWRJbkFzIjoiYWRtaW4iLCJpYXQiOjE0MjI3Nzk2Mzh9.gzSraSYS8EXBxLN_oWnFSRgCzcmJmMjLiuyu5CSpyHI"

	//act
	valid := isValid(token)

	//assert
	test.True(valid)
}

func TestIsValidToken(t *testing.T) {
	//arrange
	test := assert.New(t)
	token := "thisisavalidtoken"

	//act
	valid := isValid(token)

	//assert
	test.True(valid)
}

func TestIsInvalid(t *testing.T) {
	//arrange
	test := assert.New(t)
	invalidToken := ":[]"

	//act
	valid := isValid(invalidToken)

	//assert
	test.False(valid)
}

func TestAllowEmptyPath(t *testing.T) {
	//arrange
	test := assert.New(t)
	path := ""

	//act
	tokenPath := getPath(path)

	//asset
	test.NotEqual("", tokenPath)
}

func TestInvalidPath(t *testing.T) {
	//arrange
	test := assert.New(t)
	path := "invalidpath"
	fileToken := NewFileToken(path)

	//act
	_, err := fileToken.Read()

	//asset
	test.Error(err)
}

func TestDefaultPath(t *testing.T) {
	//arrange
	test := assert.New(t)

	//act
	path := defaultPath()

	//assert

	test.NotEqual("", path)

}

func TestWrite_MkdirAllFails(t *testing.T) {
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")
	errorMessage := "couldn't create directory"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fsmock := match.NewMockFs(ctrl)
	fs = fsmock

	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(errors.New(errorMessage))
	fileToken := NewFileToken(path)

	err := fileToken.Write("ok")

	test.EqualError(err, errorMessage)

}

func TestWrite_SetDirPermissionsFails(t *testing.T) {
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")
	errorMessage := "couldn't set permissions for the specified path"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fsmock := match.NewMockFs(ctrl)
	fs = fsmock

	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().Chmod(testdata.FixturePath(), os.FileMode(0700)).Return(errors.New(errorMessage))
	fileToken := NewFileToken(path)

	err := fileToken.Write("ok")

	test.EqualError(err, errorMessage)

}

func TestWrite_WriteFileSecurelyFails(t *testing.T) {
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")
	errorMessage := "couldn't set permissions for the specified path"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fsmock := match.NewMockFs(ctrl)
	fs = fsmock
	afs = &afero.Afero{Fs: fs}

	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().Chmod(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New(errorMessage))

	fileToken := NewFileToken(path)

	err := fileToken.Write("ok")

	test.EqualError(err, errorMessage)

}

func TestWrite_RenameFails(t *testing.T) {
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")
	errorMessage := "couldn't set permissions for the specified path"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fsmock := match.NewMockFs(ctrl)
	fileMock := match.NewMockFile(ctrl)
	fs = fsmock
	afs = &afero.Afero{Fs: fs}

	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().Chmod(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileMock, nil)
	fsmock.EXPECT().Rename("tmpFileName", path).Return(errors.New(errorMessage))
	fileMock.EXPECT().WriteString("ok").Return(2, nil)
	fileMock.EXPECT().Name().Return("tmpFileName")
	fileMock.EXPECT().Close()

	fileToken := NewFileToken(path)

	err := fileToken.Write("ok")

	test.EqualError(err, errorMessage)

}

func TestWrite_ChmodFails(t *testing.T) {
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")
	errorMessage := "couldn't set permissions for the specified path"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fsmock := match.NewMockFs(ctrl)
	fileMock := match.NewMockFile(ctrl)
	fs = fsmock
	afs = &afero.Afero{Fs: fs}

	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().Chmod(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileMock, nil)
	fsmock.EXPECT().Rename("tmpFileName", path).Return(nil)
	fsmock.EXPECT().Chmod(path, os.FileMode(0600)).Return(errors.New(errorMessage))
	fileMock.EXPECT().WriteString("ok").Return(2, nil)
	fileMock.EXPECT().Name().Return("tmpFileName")
	fileMock.EXPECT().Close()

	fileToken := NewFileToken(path)

	err := fileToken.Write("ok")

	test.EqualError(err, errorMessage)

}

func TestWrite_Success(t *testing.T) {
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fsmock := match.NewMockFs(ctrl)
	fileMock := match.NewMockFile(ctrl)
	fs = fsmock
	afs = &afero.Afero{Fs: fs}

	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().Chmod(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileMock, nil)
	fsmock.EXPECT().Rename("tmpFileName", path).Return(nil)
	fsmock.EXPECT().Chmod(path, os.FileMode(0600)).Return(nil)
	fileMock.EXPECT().WriteString("ok").Return(2, nil)
	fileMock.EXPECT().Name().Return("tmpFileName")
	fileMock.EXPECT().Close()

	fileToken := NewFileToken(path)

	err := fileToken.Write("ok")

	test.NoError(err)

}
