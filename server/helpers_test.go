package server

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
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
