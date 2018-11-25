/* util_test.go: unit tests for utility functions
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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitMessage(t *testing.T) {
	message := []string{
		"to: me@example.com",
		"subject: foo",
		"",
		"Dear Me,",
		"",
		"have a nice day.",
	}

	conf, body := SplitMessage(message)

	expectedConf := []string{
		"to: me@example.com",
		"subject: foo",
	}

	expectedBody := []string{
		"Dear Me,",
		"",
		"have a nice day.",
	}

	assert.Equal(t, expectedConf, conf)
	assert.Equal(t, expectedBody, body)
}

func TestParseTime(t *testing.T) {
	d1 := "2061-07-28"
	d2 := "2061-07-28 12:23"
	invalid1 := "2061/07/28"
	invalid2 := "2061-07-28 12:23:00"

	r, err := ParseTime(d1)
	assert.Nil(t, err)

	if r.Hour() != 0 || r.Minute() != 0 {
		t.Errorf("expected Hour and Minute to be 0, but was %s", r)
	}

	r, err = ParseTime(d2)
	assert.Nil(t, err)

	if r.Hour() != 12 || r.Minute() != 23 {
		t.Errorf("expected Hour and Minute to be 12:23, but was %s", r)
	}

	// Expected to fail
	_, err = ParseTime(invalid1)
	assert.NotNil(t, err)

	// Expected to fail
	_, err = ParseTime(invalid2)
	assert.NotNil(t, err)
}
