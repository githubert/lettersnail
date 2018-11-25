/* check.go: check a given message for problems
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
	"path/filepath"
	"sort"
)

var usageCheck =
// tag::check[]
`
Usage:
  lettersnail check [--silent] [FILE]

Options:
  --silent  Suppress output, useful for silent checks.
  FILE      Check only the given file.

If no FILE is provided, check will inspect all messages in the todo/ and the
drafts/ folder. The command exit with a code of 0, if there were no problems.
` // end::check[]

func Check(argv []string, conf *Configuration) {
	args, _ := docopt.Parse(usageCheck, argv, true, "", false)

	silent := args["--silent"].(bool)

	ok := true

	if args["FILE"] != nil {
		message, err := NewMessageFromFile(args["FILE"].(string))

		if err != nil {
			fmt.Printf("Error while reading message: %s\n", err.Error())
			os.Exit(1)
		}

		message.Conf.MergeWith(conf)
		ok = checkMessage(message, silent)
	} else {
		draftOk := checkFolder(DIR_DRAFTS, conf, silent)

		if !silent {
			// A bit of space between the drafts/ and todo/ listing
			fmt.Println()
		}

		todoOk := checkFolder(DIR_TODO, conf, silent)

		ok = draftOk && todoOk
	}

	if ok {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// Inspect all messages in a folder.
func checkFolder(folder string, conf *Configuration, silent bool) bool {
	messages := NewMessagesFromDirectory(filepath.Join(conf.Get(CONF_WORKDIR), folder))
	sort.Sort(Messages(messages))

	if !silent {
		fmt.Printf("in %s:\n", folder)
	}

	ok := true
	count := 0

	for _, message := range messages {
		count++
		message.Conf.MergeWith(conf)

		if !checkMessage(message, silent) {
			ok = false
		}
	}

	if ok && !silent {
		fmt.Printf(" All (%d) messages are valid.\n", count)
	}

	return ok
}

// Check the given message. If `silent` is false, all problems will be printed
// to stdout.
func checkMessage(message Message, silent bool) bool {
	ok := true
	errs := message.Verify()

	if errs != nil {
		ok = false
	}

	if silent {
		return ok
	}

	if !ok {
		fmt.Printf(" %s:\n", message.Name)

		for _, err := range errs {
			fmt.Printf("  %s\n", err.Error())
		}
	}

	return ok
}
