package config

import (
	_ "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AppTestSuite struct {
	suite.Suite
}

func (suite *AppTestSuite) SetupTest() {
}

func (suite *AppTestSuite) TestLoad() {
	assert.NotNil(suite.T(), App)
}

func (suite *AppTestSuite) TestPrintable() {
	assert.NotEqual(suite.T(), "", App.String())
}

func (suite *AppTestSuite) TestName() {
	assert.Equal(suite.T(), "rai", App.Name)
}

func (suite *AppTestSuite) TestLicense() {
	assert.Equal(suite.T(), "NCSA or Apache-2.0", App.License)
}

func TestAppConfig(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
