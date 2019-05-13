package server

import (
	"encoding/json"
	"errors"
	"imup/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadJsonSuccess(t *testing.T) {
	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := jsonReq("testdata/image.jpg", "image")
	w := httptest.NewRecorder()

	uploader.On("Store", mock.Anything, mock.Anything).Return(uuid.Must(uuid.NewV4()), nil)

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	s := Success{}
	err := json.Unmarshal(w.Body.Bytes(), &s)
	assert.Nil(t, err)
	assert.NotEqual(t, uuid.UUID{}, s.UUID)

	f := Failed{}
	err = json.Unmarshal(w.Body.Bytes(), &f)
	assert.Nil(t, err)
	assert.Equal(t, "", f.Error)
}

func TestUploadJsonContent(t *testing.T) {
	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := jsonReq("testdata/image.jpg", "image")
	w := httptest.NewRecorder()

	uploader.On("Store", mock.Anything, mock.Anything).Return(uuid.Must(uuid.NewV4()), nil).
		Run(func(args mock.Arguments) {
			src := args.Get(1).(io.Reader)
			actual, _ := md5Reader(src)

			file, _ := testFile("testdata/image.jpg")
			defer file.Close()
			expected, _ := md5Reader(file)

			assert.Equal(t, actual, expected)
		})

	handler.ServeHTTP(w, r)
}

func TestUploadJsonFailedWrongKey(t *testing.T) {
	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := jsonReq("testdata/image.jpg", "file")
	w := httptest.NewRecorder()

	uploader.On("Store", mock.Anything, mock.Anything).Return(uuid.Must(uuid.NewV4()), nil)

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	s := Success{}
	err := json.Unmarshal(w.Body.Bytes(), &s)
	assert.Nil(t, err)
	assert.Equal(t, uuid.UUID{}, s.UUID)

	f := Failed{}
	err = json.Unmarshal(w.Body.Bytes(), &f)
	assert.Nil(t, err)
	assert.NotEqual(t, "", f.Error)
}

func TestUploadJsonFailedUploader(t *testing.T) {
	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := jsonReq("testdata/image.jpg", "image")
	w := httptest.NewRecorder()

	// something wrong inside Uploader
	uploader.On("Store", mock.Anything, mock.Anything).Return(uuid.UUID{}, errors.New("uploader failed"))

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	s := Success{}
	err := json.Unmarshal(w.Body.Bytes(), &s)
	assert.Nil(t, err)
	assert.Equal(t, uuid.UUID{}, s.UUID)

	f := Failed{}
	err = json.Unmarshal(w.Body.Bytes(), &f)
	assert.Nil(t, err)
	assert.Equal(t, "uploader failed", f.Error)
}
