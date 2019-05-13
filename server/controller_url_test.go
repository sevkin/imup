package server

import (
	"encoding/json"
	"errors"
	"imup/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadUrlSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	file, _ := testFile("testdata/image.jpg")
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	httpmock.RegisterNoResponder(httpmock.NewBytesResponder(200, body))

	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := urlReq("http://localhost:5000/testdata/image.jpg", "image")
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

func TestUploadUrlContent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	file, _ := testFile("testdata/image.jpg")
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	httpmock.RegisterNoResponder(httpmock.NewBytesResponder(200, body))

	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := urlReq("http://localhost:5000/testdata/image.jpg", "image")
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

func TestUploadUrlFailedNotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, ""))

	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := urlReq("http://localhost:5000/testdata/image.jpg", "image")
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

func TestUploadUrlFailedUploader(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	file, _ := testFile("testdata/image.jpg")
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	httpmock.RegisterNoResponder(httpmock.NewBytesResponder(200, body))

	uploader := new(mocks.Uploader)
	handler := newController(uploader)

	r := urlReq("http://localhost:5000/testdata/image.jpg", "image")
	w := httptest.NewRecorder()

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
