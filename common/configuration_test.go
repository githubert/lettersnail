/* configuration_test.go: unit tests for the configuration module
 *
 * Copyright (C) 2016-2018 Clemens Fries <github-lettersnail@xenoworld.de>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package common

import (
	"github.com/go-ini/ini"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfiguration_MergeWith(t *testing.T) {
	src := Configuration{
		Data: map[string]string{
			"port":   "1",
			"server": "example.com",
		},
	}

	dst := Configuration{
		Data: map[string]string{
			"server":  "example.net",
			"workdir": "/home/foo",
		},
	}

	expected := map[string]string{
		"port":    "1",
		"server":  "example.com",
		"workdir": "/home/foo",
	}

	dst.MergeWith(&src)

	assert.Equal(t, expected, dst.Data)
}

func TestConfiguration_MergeWithIni(t *testing.T) {
	dst := map[string]string{
		"port":   "1",
		"server": "example.com",
	}

	iniContents := []byte(`
	[foo]
	port = 1234
	workdir = /home/foo
	`)

	expected := map[string]string{
		"port":    "1234",
		"server":  "example.com",
		"workdir": "/home/foo",
	}

	cfg, _ := ini.Load(iniContents)

	section := cfg.Section("foo")

	mergeIniSection(section, &dst)

	assert.Equal(t, expected, dst)
}

func TestConfiguration_Load(t *testing.T) {
	text := []string{
		"to: me@example.com",
		"from: ",
		"subject",
		"",
	}

	expected := map[string]string{
		"to":      "me@example.com",
		"from":    "",
		"subject": "",
	}

	conf := Configuration{}
	conf.Load(text)

	assert.Equal(t, expected, conf.Data)
}

func TestConfiguration_Get(t *testing.T) {
	conf := NewConfiguration()

	conf.Set("foo", "bar")

	assert.Equal(t, "bar", conf.Get("foo"))
}
