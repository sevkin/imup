package uploader

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gofrs/uuid"
	"golang.org/x/sync/semaphore"
)

type (
	// Uploader just upload content from src, then return it`s uuid
	Uploader interface {
		Store(ctx context.Context, src io.Reader) (uuid.UUID, error)
	}

	dirUploader struct {
		storage  string
		thumbCMD string
		semThumb *semaphore.Weighted
	}
)

// Store content
func (u *dirUploader) Store(ctx context.Context, src io.Reader) (uuid.UUID, error) {
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

		var UUID uuid.UUID
		UUID, err = uuid.NewV4()
		if err == nil {
			fname := filepath.Join(u.storage, UUID.String()+ext)

			var file *os.File
			file, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0644)
			if err == nil {
				func() {
					defer file.Close()
					_, err = io.Copy(file, buf)
				}()
				if err == nil {
					tname := filepath.Join(u.storage, UUID.String()+".thumb.100x100"+ext)

					func() {
						u.semThumb.Acquire(ctx, 1)
						defer u.semThumb.Release(1)
						err = exec.Command(u.thumbCMD, fname, tname).Run()
					}()
					if err == nil {
						return UUID, nil
					}
				}
			}
		}
	}
	return uuid.UUID{}, err
}

// NewDirUploader returns new Uploader instance
func NewDirUploader(storage, thumbCMD string) Uploader {
	return &dirUploader{
		storage:  storage,
		thumbCMD: thumbCMD,
		semThumb: semaphore.NewWeighted(int64(runtime.NumCPU())),
	}
}
