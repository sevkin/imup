package uploader

import (
	"errors"
	"io"

	"github.com/gofrs/uuid"
)

type (
	// Uploader just upload content from src, then return it`s uuid
	Uploader interface {
		Store(src io.Reader) (uuid.UUID, error)
	}

	dirUploader struct {
	}
)

// Store content
func (u *dirUploader) Store(src io.Reader) (uuid.UUID, error) {
	return uuid.UUID{}, errors.New("not implemented")
}

// NewDirUploader returns new Uploader instance
func NewDirUploader() Uploader {
	return &dirUploader{}
}
