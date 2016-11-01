package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	. "github.com/githubert/clockrotz/common"
)

func TestExpandTilde(t *testing.T) {
	c := NewConfiguration()

	c.Set(CONF_WORKDIR, "/foo")
	expandTilde(c)
	assert.Equal(t, "/foo", c.Get(CONF_WORKDIR))

	c.Set(CONF_WORKDIR, "~/foo")
	expandTilde(c)
	assert.Equal(t, userHome() + "/foo", c.Get(CONF_WORKDIR))

	// Shortest possible
	c.Set(CONF_WORKDIR, "~/")
	expandTilde(c)
	assert.Equal(t, userHome(), c.Get(CONF_WORKDIR))

	// Single character workdir with only a tilde
	c.Set(CONF_WORKDIR, "~")
	expandTilde(c)
	assert.Equal(t, userHome(), c.Get(CONF_WORKDIR))

	// Ignored tilde
	c.Set(CONF_WORKDIR, "~test")
	expandTilde(c)
	assert.Equal(t, "~test", c.Get(CONF_WORKDIR))
}