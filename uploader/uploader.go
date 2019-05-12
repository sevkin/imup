package uploader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

type (
	// Uploader just upload content from src, then return it`s uuid
	Uploader interface {
		Store(src io.Reader) (uuid.UUID, error)
	}

	dirUploader struct {
		storage string
	}
)

// Store content
func (u *dirUploader) Store(src io.Reader) (uuid.UUID, error) {
	if src == nil {
		return uuid.UUID{}, errors.New("empty source")
	}

	buf := bufio.NewReader(src)
	hdr, err := buf.Peek(512)
	if err == nil {
		contentType := http.DetectContentType(hdr)
		var ext string
		switch {
		case contentType == "image/jpeg":
			ext = ".jpg"
		case contentType == "image/png":
			ext = ".png"
		case contentType == "image/gif":
			ext = ".gif"
		default:
			return uuid.UUID{}, fmt.Errorf("unsupported content: %s", contentType)
		}

		UUID, err := uuid.NewV4()
		if err == nil {
			fname := filepath.Join(u.storage, UUID.String()+ext)

			file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0644)
			if err == nil {
				defer file.Close()
				_, err = io.Copy(file, buf)
				if err == nil {
					// TODO make thumbnail
					return UUID, nil
				}
			}
		}
	}
	return uuid.UUID{}, err
}

// NewDirUploader returns new Uploader instance
func NewDirUploader(storage string) Uploader {
	return &dirUploader{
		storage: storage,
	}
}
