/* message_test.go: unit tests for the message module
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
	"path/filepath"
	"sort"
	"testing"
)

func TestReadMessagesAndVerify(t *testing.T) {
	workdir := "testdata"

	conf := NewConfiguration()
	conf.Set("from", "me@example.com")

	messages := NewMessagesFromDirectory(filepath.Join(workdir, "todo"))

	for _, message := range messages {
		// If we don't to this, it will fail because "from" is missing
		message.Conf.MergeWith(conf)

		if err := message.Verify(); err != nil {
			t.Errorf("Verification of message '%s' failed with %s", message.Name, err)
		}
	}
}

func TestMessage_Get(t *testing.T) {
	message := Message{Conf: *NewConfiguration()}

	message.Conf.Set("foo", "bar")

	assert.Equal(t, "bar", message.Get("foo"))
}

func TestMessagesSort(t *testing.T) {
	messages := []Message{
		{
			Name: "1",
			Conf: Configuration{
				Data: map[string]string{
					"date": "2000-01-01",
				}}},
		{
			Name: "2",
			Conf: Configuration{
				Data: map[string]string{
					"date": "2010-01-01",
				}}},
		{
			Name: "3",
			Conf: Configuration{
				Data: map[string]string{
					"date": "2005-01-01",
				}}},
	}

	sort.Sort(Messages(messages))

	for i, k := range []string{"1", "3", "2"} {
		if messages[i].Name != k {
			t.Errorf("expected: %s, but was %s", k, messages[i].Name)
		}
	}
}
