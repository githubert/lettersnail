/* clockrotz.go: main program
 *
 * Copyright (C) 2016 Clemens Fries <github-clockrotz@xenoworld.de>
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

package main

import (
	"fmt"
	"github.com/docopt/docopt.go"
	"os"
	"os/user"
	"path/filepath"
	"github.com/githubert/clockrotz/cmd"
	. "github.com/githubert/clockrotz/common"
)


var usageMain =
// tag::main[]
`
Usage:
  clockrotz [options] <command> [<args>...]

Options:
  --workdir=WORKDIR  Working directory, defaults to $HOME/clockrotz
  --config=CONFIG    Use the given INI file instead of .config/clockrotz.ini
  --help             Show this help.

The following commands are available:
  next    Show messages due in the next days.
  check   Check if a message is complete.
  debug   Print the effective configuration for a message, and the resulting
          email message.
  create  Open a new message in an editor.
  run     Send out pending messages.
` // end::main[]

var defaultConf = Configuration{
	Data: map[string]string{
		CONF_DAYS:        "7",
		CONF_SMTP_PORT:   "587",
		CONF_SMTP_SERVER: "localhost",
		CONF_NOT_BEFORE:  "00:00",
		CONF_NOT_AFTER:   "23:59",
	},
}


// Retrieve the user's home directory.
func userHome() string {
	u, err := user.Current()

	// I don't think we can go on without finding the user's home directory
	if err != nil {
		panic(err)
	}

	return u.HomeDir
}

// Retrieve the user's config home, which is either defined through the environment variable
// `XDG_CONFIG_HOME` or is assumed to be a folder called `.config` in the user's home
// directory (as retrieved by `userHome()`).
func configHome() string {
	value, present := os.LookupEnv("XDG_CONFIG_HOME")

	configHome := value

	if !present || value == "" || !filepath.IsAbs(value) {
		configHome = filepath.Join(userHome(), ".config")
	}

	return configHome
}

// Expand `~/`, if the working directory begins with it.
func expandTilde(conf *Configuration) {
	workdir := conf.Get(CONF_WORKDIR)

	if len(workdir) > 1 && workdir[0:2] == "~/" {
		conf.Set(CONF_WORKDIR, filepath.Join(userHome(), workdir[2:]))
	} else if len(workdir) == 1 && workdir == "~" {
		conf.Set(CONF_WORKDIR, userHome())
	}
}

// This will create all necessary folders in the working directory.
func createFolders(workdir string) error {
	for _, dir := range []string{DIR_TODO, DIR_DONE, DIR_ERRORS, DIR_DRAFTS} {
		err := os.MkdirAll(filepath.Join(workdir, dir), 0777)

		if err != nil {
			return err
		}
	}

	return nil
}

// Complain about files in DIR_ERRORS.
func alertIfErrors(workdir string) {
	errorsDir := filepath.Join(workdir, DIR_ERRORS)

	count := 0

	filepath.Walk(errorsDir, func(_ string, info os.FileInfo, _ error) error {
		if info.Mode().IsRegular() {
			count++
		}

		return nil
	})

	if count > 0 {
		fmt.Println("There were failed messages.")
		fmt.Printf("Please inspect the contents of the %s/ directory.\n", DIR_ERRORS)
	}
}

func main() {
	userHome := userHome()

	// This is the configuration that will be cobbled together from the
	// default configuration ('defaultConf'), the INI and the command line
	// arguments.
	sessionConf := NewConfiguration()

	sessionConf.MergeWith(&defaultConf)

	args, _ := docopt.Parse(usageMain, nil, true, "", true)

	// Determine configuration file name, defaults to `clockrotz.ini` in
	// the user's config home (usually `.config`).
	if args["--config"] != nil {
		sessionConf.Set(CONF_CONFIG_FILENAME, args["--config"].(string))
	} else {
		sessionConf.Set(CONF_CONFIG_FILENAME, filepath.Join(configHome(), "clockrotz.ini"))
	}

	// Determine the desired command (run, next, ...)
	command := args["<command>"].(string)

	// Load INI configuration. This will merge the '[default]' section with
	// the optional section named after the desired command ('cmd').
	sessionConf.MergeWithIni(command)

	// Determine working directory, defaults to the folder `clockrotz` in
	// the user's home directory. '--workdir' on the command line takes
	// precedence over the 'workdir' in the INI.
	if args["--workdir"] != nil {
		sessionConf.Set(CONF_WORKDIR, args["--workdir"].(string))
	} else {
		// If no working directory is set in the INI...
		if sessionConf.Get(CONF_WORKDIR) == "" {
			sessionConf.Set(CONF_WORKDIR, filepath.Join(userHome, "clockrotz"))
		}
	}

	// Expand the `~/` in workdir.
	expandTilde(sessionConf)

	// Make sure that all necessary folders exist.
	createFolders(sessionConf.Get(CONF_WORKDIR))

	// See if there are files in DIR_ERRORS and alert the user.
	alertIfErrors(sessionConf.Get(CONF_WORKDIR))

	commandArgs  := []string{command}
	commandArgs = append(commandArgs, args["<args>"].([]string)...)

	switch command {
	case "next":
		cmd.Next(commandArgs, sessionConf)
	case "create":
		cmd.Create(commandArgs, sessionConf)
	case "check":
		cmd.Check(commandArgs, sessionConf)
	case "run":
		cmd.Run(commandArgs, sessionConf)
	case "debug":
		cmd.Debug(commandArgs, sessionConf)
	}
}