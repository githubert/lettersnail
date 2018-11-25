/* create_test.go: unit tests for 'create' command
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
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestEditor(t *testing.T) {
	os.Setenv("VISUAL", "")
	os.Setenv("EDITOR", "")

	assert.Equal(t, "vi", editor())

	os.Setenv("EDITOR", "editor")
	assert.Equal(t, "editor", editor())

	os.Setenv("VISUAL", "visual")
	assert.Equal(t, "visual", editor())
}