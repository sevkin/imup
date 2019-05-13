package uploader

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type (
	UploaderTestSuite struct {
		suite.Suite
		Storage  string
		Ctx      context.Context
		Uploader Uploader
	}
)

func (suite *UploaderTestSuite) SetupSuite() {
	suite.Storage = filepath.Join(testDir(), "tmp")
	os.Mkdir(suite.Storage, 0755)
	suite.Ctx = context.TODO()
}

func (suite *UploaderTestSuite) SetupTest() {
	suite.Uploader = NewDirUploader(suite.Storage, "/bin/true") // windows sucks?
}

func (suite *UploaderTestSuite) TearDownTest() {
	os.RemoveAll(filepath.Join(suite.Storage, "*"))
}

func (suite *UploaderTestSuite) TearDownSuite() {
	os.RemoveAll(suite.Storage)
}

func (suite *UploaderTestSuite) TestStore() {
	file, _ := testFile("testdata/image.jpg")
	defer file.Close()
	UUID, err := suite.Uploader.Store(suite.Ctx, file)
	suite.Nil(err)
	suite.NotEqual(uuid.UUID{}, UUID)
}

func (suite *UploaderTestSuite) TestPNG() {
	file, _ := testFile("testdata/image.png")
	defer file.Close()
	UUID, err := suite.Uploader.Store(suite.Ctx, file)
	suite.Nil(err)
	suite.NotEqual(uuid.UUID{}, UUID)
}

func (suite *UploaderTestSuite) TestGIF() {
	file, _ := testFile("testdata/image.gif")
	defer file.Close()
	UUID, err := suite.Uploader.Store(suite.Ctx, file)
	suite.Nil(err)
	suite.NotEqual(uuid.UUID{}, UUID)
}

func (suite *UploaderTestSuite) TestUnsupported() {
	file, _ := testFile("testdata/image.zip")
	defer file.Close()
	_, err := suite.Uploader.Store(suite.Ctx, file)
	suite.NotNil(err)
}

func (suite *UploaderTestSuite) TestNil() {
	_, err := suite.Uploader.Store(suite.Ctx, nil)
	suite.NotNil(err)
}

func TestUploaderTestSuite(t *testing.T) {
	suite.Run(t, new(UploaderTestSuite))
}

func TestThumbFailed(t *testing.T) {
	storage := filepath.Join(testDir(), "tmp")
	os.Mkdir(storage, 0755)
	defer os.RemoveAll(storage)

	uploader := NewDirUploader(storage, "/bin/false")

	file, _ := testFile("testdata/image.jpg")
	defer file.Close()
	_, err := uploader.Store(context.TODO(), file)
	assert.NotNil(t, err)
}

// /////////////////////////////////////////////////////////////////////////////

func testDir() string {
	_, testfilename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(testfilename), "..")
}

func testFile(fname string) (*os.File, error) {
	fname = filepath.Join(testDir(), fname)
	return os.Open(fname)
}
