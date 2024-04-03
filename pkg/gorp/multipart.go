package gorp

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
)

// Multipart is an interface that allows to pass over different types
// of multipart data sources
type Multipart interface {
	// Load loads multipart data
	Load() (fileName, contentType string, reader io.Reader, err error)
}

// FileMultipart is a multipart content in form of file
type FileMultipart struct {
	*os.File
}

//nolint:nonamedreturns // for readability
func (fm *FileMultipart) Load() (fileName, contentType string, reader io.Reader, err error) {
	if fm.File == nil {
		return "", "", nil, errors.New("file shouldn't be nil")
	}
	fName := fm.File.Name()
	if _, sErr := os.Stat(fName); os.IsNotExist(sErr) {
		return "", "", nil, fmt.Errorf("file %s does not exist", fName)
	}
	contentType = mime.TypeByExtension(filepath.Ext(fName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return filepath.Base(fName), contentType, fm.File, nil
}

// ReaderMultipart is a multipart content in form of io.Reader
type ReaderMultipart struct {
	FileName, ContentType string
	io.Reader
}

//nolint:nonamedreturns // for readability
func (fm *ReaderMultipart) Load() (fileName, contentType string, reader io.Reader, err error) {
	if fm.FileName == "" {
		return "", "", nil, errors.New("multipart filename shouldn't be nil")
	}
	return fm.FileName, fm.ContentType, fm, nil
}
