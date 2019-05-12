package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"imup/mocks"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testFile(fname string) (*os.File, error) {
	_, testfilename, _, _ := runtime.Caller(0)
	fname = filepath.Join(filepath.Dir(testfilename), "..", fname)
	return os.Open(fname)
}

func multipartReq(filename, fieldname string) *http.Request {
	body, ct := multipartBody(filename, fieldname)
	r := httptest.NewRequest("POST", "/upload/form", body)
	r.Header.Add("Content-Type", ct)

	return r
}

func multipartBody(filename, fieldname string) (io.Reader, string) {
	file, _ := testFile(filename)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, _ := writer.CreateFormFile(fieldname, filepath.Base(file.Name()))
	io.Copy(part, file)
	return body, writer.FormDataContentType()
}

func TestUploadFormSuccess(t *testing.T) {
	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := multipartReq("testdata/image.jpg", "image")
	w := httptest.NewRecorder()

	uploader.On("Store", mock.Anything).Return(uuid.Must(uuid.NewV4()), nil)

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	s := Success{}
	err := json.Unmarshal(w.Body.Bytes(), &s)
	assert.Nil(t, err)
	assert.NotEqual(t, uuid.UUID{}, s.UUID)

	f := Failed{}
	err = json.Unmarshal(w.Body.Bytes(), &f)
	assert.Nil(t, err)
	assert.Equal(t, "", f.Message)
}

func TestUploadFormFailedHeader(t *testing.T) {
	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := multipartReq("testdata/image.jpg", "image")
	// multipart/form-data; boundary=<something wrong>
	r.Header["Content-Type"][0] = r.Header["Content-Type"][0] + "wrong"
	w := httptest.NewRecorder()

	uploader.On("Store", mock.Anything).Return(uuid.Must(uuid.NewV4()), nil)

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	s := Success{}
	err := json.Unmarshal(w.Body.Bytes(), &s)
	assert.Nil(t, err)
	assert.Equal(t, uuid.UUID{}, s.UUID)

	f := Failed{}
	err = json.Unmarshal(w.Body.Bytes(), &f)
	assert.Nil(t, err)
	assert.NotEqual(t, "", f.Message)
}

func TestUploadFormFailedField(t *testing.T) {
	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	// expected "image" but actual "file"
	r := multipartReq("testdata/image.jpg", "file")
	w := httptest.NewRecorder()

	uploader.On("Store", mock.Anything).Return(uuid.Must(uuid.NewV4()), nil)

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	s := Success{}
	err := json.Unmarshal(w.Body.Bytes(), &s)
	assert.Nil(t, err)
	assert.Equal(t, uuid.UUID{}, s.UUID)

	f := Failed{}
	err = json.Unmarshal(w.Body.Bytes(), &f)
	assert.Nil(t, err)
	assert.NotEqual(t, "", f.Message)
}

func TestUploadFormFailedUploader(t *testing.T) {
	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := multipartReq("testdata/image.jpg", "image")
	w := httptest.NewRecorder()

	// something wrong inside Uploader
	uploader.On("Store", mock.Anything).Return(uuid.UUID{}, errors.New("uploader failed"))

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	s := Success{}
	err := json.Unmarshal(w.Body.Bytes(), &s)
	assert.Nil(t, err)
	assert.Equal(t, uuid.UUID{}, s.UUID)

	f := Failed{}
	err = json.Unmarshal(w.Body.Bytes(), &f)
	assert.Nil(t, err)
	assert.Equal(t, "uploader failed", f.Message)
}
