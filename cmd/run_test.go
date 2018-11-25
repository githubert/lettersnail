/* run_test.go: unit tests for 'run' command
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
package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildTime(t *testing.T) {
	h := 8
	m := 30

	t1 := buildTime(time.Now(), h, m, true)

	assert.Equal(t, t1.Hour(), 8)
	assert.Equal(t, t1.Minute(), 30)
	assert.Equal(t, t1.Second(), 0)

	t2 := buildTime(time.Now(), h, m, false)

	assert.Equal(t, t2.Second(), 59)
}
