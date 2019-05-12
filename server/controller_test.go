package server

import (
	"bytes"
	"encoding/json"
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
)

func testFile(fname string) (*os.File, error) {
	_, testfilename, _, _ := runtime.Caller(0)
	fname = filepath.Join(filepath.Dir(testfilename), "..", fname)
	return os.Open(fname)
}

func multipartReq(fname string) *http.Request {
	file, _ := testFile(fname)
	defer file.Close()

	body := &bytes.Buffer{}
	ct := func() string {
		writer := multipart.NewWriter(body)
		defer writer.Close()
		part, _ := writer.CreateFormFile("image", filepath.Base(file.Name()))
		io.Copy(part, file)
		return writer.FormDataContentType()
	}()

	r := httptest.NewRequest("POST", "/upload/form", body)
	r.Header.Add("Content-Type", ct)

	return r
}

func TestUploadForm(t *testing.T) {
	r := multipartReq("testdata/image.jpg")
	w := httptest.NewRecorder()
	handler := newController()

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
