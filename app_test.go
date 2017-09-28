package config

import (
	_ "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AppTestSuite ...
type AppTestSuite struct {
	suite.Suite
}

// SetupTest ...
func (suite *AppTestSuite) SetupTest() {
}

// TestLoad ...
func (suite *AppTestSuite) TestLoad() {
	assert.NotNil(suite.T(), App)
}

// TestPrintable ...
func (suite *AppTestSuite) TestPrintable() {
	assert.NotEqual(suite.T(), "", App.String())
}

// TestName ...
func (suite *AppTestSuite) TestName() {
	assert.Equal(suite.T(), "rai", App.Name)
}

// TestLicense ...
func (suite *AppTestSuite) TestLicense() {
	assert.Equal(suite.T(), "NCSA or Apache-2.0", App.License)
}

// TestAppConfig ...
func TestAppConfig(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
