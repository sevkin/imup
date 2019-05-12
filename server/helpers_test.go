package server

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testFile(fname string) (*os.File, error) {
	_, testfilename, _, _ := runtime.Caller(0)
	fname = filepath.Join(filepath.Dir(testfilename), "..", fname)
	return os.Open(fname)
}

// /////////////////////////////////////////////////////////////////////////////

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

// /////////////////////////////////////////////////////////////////////////////

func jsonReq(filename, fieldname string) *http.Request {
	body, ct := jsonBody(filename, fieldname)
	r := httptest.NewRequest("POST", "/upload/json", body)
	r.Header.Add("Content-Type", ct)

	return r
}

func jsonBody(filename, fieldname string) (io.Reader, string) {
	pr, pw := io.Pipe()

	go func() {
		file, _ := testFile(filename)
		defer file.Close()

		_, err := io.Copy(pw, jsonPipe(base64Pipe(file), fieldname))

		if err != nil {
			pw.CloseWithError(err)
		} else {
			pw.Close()
		}
	}()

	return pr, "application/json"
}

// /////////////////////////////////////////////////////////////////////////////

func urlReq(url, key string) *http.Request {
	r := httptest.NewRequest("GET", "/upload/url", nil)
	q := r.URL.Query()
	q.Add(key, url)
	r.URL.RawQuery = q.Encode()
	return r
}

// /////////////////////////////////////////////////////////////////////////////

func base64Pipe(src io.Reader) io.Reader {
	pr, pw := io.Pipe()
	encoder := base64.NewEncoder(base64.StdEncoding, pw)

	go func() {
		_, err := io.Copy(encoder, src)
		encoder.Close()

		if err != nil {
			pw.CloseWithError(err)
		} else {
			pw.Close()
		}
	}()

	return pr
}

func TestBase64Pipe(t *testing.T) {
	src := strings.NewReader("Hello, World!")
	dst, err := ioutil.ReadAll(base64Pipe(src))
	assert.Nil(t, err)
	assert.Equal(t, []byte("SGVsbG8sIFdvcmxkIQ=="), dst)
}

// /////////////////////////////////////////////////////////////////////////////

func jsonPipe(src io.Reader, fieldname string) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		fmt.Fprintf(pw, `{"%s":"`, fieldname)
		_, err := io.Copy(pw, src)
		fmt.Fprint(pw, `"}`)

		if err != nil {
			pw.CloseWithError(err)
		} else {
			pw.Close()
		}
	}()

	return pr
}

func TestJsonPipe(t *testing.T) {
	src := strings.NewReader("Hello, World!")
	dst, err := ioutil.ReadAll(jsonPipe(src, "image"))
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"image":"Hello, World!"}`), dst)

	src = strings.NewReader("Hello, World!")
	dst, err = ioutil.ReadAll(jsonPipe(base64Pipe(src), "image"))
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"image":"SGVsbG8sIFdvcmxkIQ=="}`), dst)
}

// /////////////////////////////////////////////////////////////////////////////

func md5Reader(src io.Reader) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, src)
	return fmt.Sprintf("%x", hash.Sum(nil)), err
}

func TestMd5Reader(t *testing.T) {
	src := strings.NewReader("Hello, World!")
	hash, err := md5Reader(src)
	assert.Nil(t, err)
	assert.Equal(t, "65a8e27d8879283831b664bd8b7f0ad4", hash)
}

// /////////////////////////////////////////////////////////////////////////////
