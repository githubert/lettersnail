/* debug.go: show debug information for a provided message
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
	"github.com/docopt/docopt.go"
	. "github.com/githubert/lettersnail/common"
	"os"
)

var usageDebug =
// tag::debug[]
`
Usage:
  lettersnail debug FILENAME

Options:
  --help    Show this help.
  FILENAME  File name of the message to inspect.
` // end::debug[]

func Debug(argv []string, conf *Configuration) {
	args, _ := docopt.Parse(usageDebug, argv, true, "", false)

	message, err := NewMessageFromFile(args["FILENAME"].(string))

	if err != nil {
		fmt.Printf("Error while reading file: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Configuration\n-------------")

	message.Conf.MergeWith(conf)

	for _, line := range message.Conf.DumpConfig() {
		fmt.Println(line)
	}

	// TODO: Verify Message?
	e, err := prepareEmail(&message)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	m, err := e.Bytes()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("\nEmail message\n-------------")
	fmt.Println(string(m))
}
