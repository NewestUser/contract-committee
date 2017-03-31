package config

import (
	"testing"

	. "launchpad.net/gocheck"
)

type ConfigSuite struct{}

var _ = Suite(new(ConfigSuite))

func TestAll(t *testing.T) {
	TestingT(t)
}

func (suite *ConfigSuite) TestLoadingConfiguration(c *C) {
	actualConfig, _ := New([]string{"--db-host", "dev.telcong.com", "--db-name", "test"})

	expectedConfig := &Config{"dev.telcong.com", "test"}

	c.Assert(actualConfig, DeepEquals, expectedConfig)
}

func (suite *ConfigSuite) TestLoadingUndefinedDatabaseName(c *C) {
	_, err := New([]string{"--db-host", "dev.telcong.com"})

	c.Assert(err, NotNil)
}

func (suite *ConfigSuite) TestLoadingUndefinedDatabaseHost(c *C) {
	_, err := New([]string{"--db-name", "test"})

	c.Assert(err, NotNil)
}
