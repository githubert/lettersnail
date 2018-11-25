/* next.go: show upcoming messages
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
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

var usageNext =
// tag::next[]
`
Usage:
  lettersnail next [--days=DAYS | --all]

Options:
  --help       Show this help.
  --days=DAYS  List messages for the next DAYS days. (default: 7)
  --all        List all pending messages.
` // end::next[]

func Next(argv []string, conf *Configuration) {
	args, _ := docopt.Parse(usageNext, argv, true, "", false)

	if args["--days"] != nil {
		conf.Set("days", args["--days"].(string))
	}

	days, err := strconv.ParseInt(conf.Get(CONF_DAYS), 10, 0)

	if err != nil {
		fmt.Println("Error while parsing value of --days")
		return
	}

	all := args["--all"].(bool)

	future := buildTime(time.Now().AddDate(0, 0, int(days)), 23, 59, false)

	messages := NewMessagesFromDirectory(filepath.Join(conf.Get(CONF_WORKDIR), DIR_TODO))
	sort.Sort(messages)

	if all {
		fmt.Printf("Showing all messages.\n\n")
	} else {
		fmt.Printf("Showing messages before %s.\n\n", future.Format(DATETIME_FORMAT))
	}

	count := 0

	for _, message := range messages {
		message.Conf.MergeWith(conf)

		if errs := message.Verify(); errs != nil {
			fmt.Printf("Error in message \"%s\". Please run 'lettersnail check'.\n", message.Name)
			continue
		}

		// Errors are caught already by Verify()
		messageDate, _ := ParseTime(message.Get(CONF_DATE))

		if all || messageDate.Before(future) {
			count++
			fmt.Printf("%s  %s (%s)\n", messageDate.Format(DATETIME_FORMAT), message.Get(CONF_SUBJECT), message.Name)
		}
	}

	if count == 0 {
		fmt.Println("No messages.")
	}
}
