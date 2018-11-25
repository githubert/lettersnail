/* util.go: various utility functions
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
	"strings"
	"time"
)

// Split a text in two parts on the first blank line.
func SplitMessage(text []string) ([]string, []string) {
	conf := []string{}
	body := []string{}

	inBody := false

	for _, line := range text {
		if !inBody && strings.TrimSpace(line) == "" {
			inBody = true
			continue // skip the separating line
		}

		if inBody {
			body = append(body, line)
		} else {
			conf = append(conf, line)
		}
	}

	return conf, body
}

// Parse a time which may either be just a date or a date with time.
func ParseTime(datetime string) (time.Time, error) {
	result, err := time.ParseInLocation(DATE_FORMAT, datetime, time.Now().Location())

	if err == nil {
		return result, nil
	}

	result, err = time.ParseInLocation(DATETIME_FORMAT, datetime, time.Now().Location())

	if err == nil {
		return result, nil
	}

	return result, err
}
