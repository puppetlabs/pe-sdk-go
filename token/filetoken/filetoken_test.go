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

func Test_Read_OK(t *testing.T) {
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

func Test_IsValid_JWTOK(t *testing.T) {
	//arrange
	test := assert.New(t)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dnZWRJbkFzIjoiYWRtaW4iLCJpYXQiOjE0MjI3Nzk2Mzh9.gzSraSYS8EXBxLN_oWnFSRgCzcmJmMjLiuyu5CSpyHI"

	//act
	valid := isValid(token)

	//assert
	test.True(valid)
}

func Test_IsValid_Token(t *testing.T) {
	//arrange
	test := assert.New(t)
	token := "thisisavalidtoken"

	//act
	valid := isValid(token)

	//assert
	test.True(valid)
}

func Test_IsValid_Invalid(t *testing.T) {
	//arrange
	test := assert.New(t)
	invalidToken := ":[]"

	//act
	valid := isValid(invalidToken)

	//assert
	test.False(valid)
}

func Test_GetPath_AllowEmptyPath(t *testing.T) {
	//arrange
	test := assert.New(t)
	path := ""

	//act
	tokenPath := getPath(path)

	//asset
	test.NotEqual("", tokenPath)
}

func Test_Read_InvalidPath(t *testing.T) {
	//arrange
	test := assert.New(t)
	path := "invalidpath"
	fileToken := NewFileToken(path)

	//act
	_, err := fileToken.Read()

	//asset
	test.Error(err)
}

func Test_DefaultPath(t *testing.T) {
	//arrange
	test := assert.New(t)

	//act
	path := defaultPath()

	//assert
	test.NotEqual("", path)
}

func Test_Write_MkdirAllFails(t *testing.T) {
	//arrange
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")
	errorMessage := "couldn't create directory"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fsmock := match.NewMockFs(ctrl)
	fs = fsmock
	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(errors.New(errorMessage))

	fileToken := NewFileToken(path)

	//act
	err := fileToken.Write("ok")

	//assert
	test.EqualError(err, errorMessage)
}

func Test_Write_SetDirPermissionsFails(t *testing.T) {
	//arrange
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

	//act
	err := fileToken.Write("ok")

	//assert
	test.EqualError(err, errorMessage)
}

func Test_Write_WriteFileSecurelyFails(t *testing.T) {
	//arrange
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

	//act
	err := fileToken.Write("ok")

	//assert
	test.EqualError(err, errorMessage)
}

func Test_Write_RenameFails(t *testing.T) {
	//arrange
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")
	errorMessage := "couldn't set permissions for the specified path"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileMock := match.NewMockFile(ctrl)
	fileMock.EXPECT().WriteString("ok").Return(2, nil)
	fileMock.EXPECT().Name().Return("tmpFileName")
	fileMock.EXPECT().Close()

	fsmock := match.NewMockFs(ctrl)
	fs = fsmock
	afs = &afero.Afero{Fs: fs}
	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().Chmod(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileMock, nil)
	fsmock.EXPECT().Rename("tmpFileName", path).Return(errors.New(errorMessage))

	fileToken := NewFileToken(path)

	//act
	err := fileToken.Write("ok")

	//assert
	test.EqualError(err, errorMessage)
}

func Test_Write_ChmodFails(t *testing.T) {
	//arrange
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")
	errorMessage := "couldn't set permissions for the specified path"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileMock := match.NewMockFile(ctrl)
	fileMock.EXPECT().WriteString("ok").Return(2, nil)
	fileMock.EXPECT().Name().Return("tmpFileName")
	fileMock.EXPECT().Close()

	fsmock := match.NewMockFs(ctrl)
	fs = fsmock
	afs = &afero.Afero{Fs: fs}
	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().Chmod(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileMock, nil)
	fsmock.EXPECT().Rename("tmpFileName", path).Return(nil)
	fsmock.EXPECT().Chmod(path, os.FileMode(0600)).Return(errors.New(errorMessage))

	fileToken := NewFileToken(path)

	//act
	err := fileToken.Write("ok")

	//assert
	test.EqualError(err, errorMessage)
}

func Test_Write_Success(t *testing.T) {
	//arrange
	test := assert.New(t)
	path := filepath.Join(testdata.FixturePath(), "token")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileMock := match.NewMockFile(ctrl)
	fileMock.EXPECT().WriteString("ok").Return(2, nil)
	fileMock.EXPECT().Name().Return("tmpFileName")
	fileMock.EXPECT().Close()

	fsmock := match.NewMockFs(ctrl)
	fs = fsmock
	afs = &afero.Afero{Fs: fs}
	fsmock.EXPECT().MkdirAll(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().Chmod(testdata.FixturePath(), os.FileMode(0700)).Return(nil)
	fsmock.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(fileMock, nil)
	fsmock.EXPECT().Rename("tmpFileName", path).Return(nil)
	fsmock.EXPECT().Chmod(path, os.FileMode(0600)).Return(nil)

	fileToken := NewFileToken(path)

	//act
	err := fileToken.Write("ok")

	//assert
	test.NoError(err)
}
