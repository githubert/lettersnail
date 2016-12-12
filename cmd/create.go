/* create.go: help creating a new message
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
package cmd

import (
	"bufio"
	"fmt"
	"github.com/docopt/docopt.go"
	. "github.com/githubert/clockrotz/common"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var usageCreate =
// tag::create[]
`
Usage:
  clockrotz create [--draft=FILE] [options]

Options:
  --help           Show this help.
  --to=ADDR        Destination address.
  --from=ADDR      Sender address.
  --subject=ADDR   Short subject.
  --cc=ADDR        Set "Cc".
  --bcc=ADDR       Set "Bcc".
  --reply-to=ADDR  Set "Reply-To".
  --draft=FILE     Use FILE from the drafts/ folder as template.
` // end::create[]

func Create(argv []string, conf *Configuration) {
	args, _ := docopt.Parse(usageCreate, argv, true, "", false)

	tmpFile, err := ioutil.TempFile("", "clockrotz")

	if err != nil {
		fmt.Printf("Error while creating temporary file: %s\n", err.Error())
		os.Exit(1)
	}

	// We close the file right away, because we need only its soulless
	// shell, mwhaha.
	tmpFile.Close()

	draftsDir := filepath.Join(conf.Get(CONF_WORKDIR), DIR_DRAFTS)
	todoDir := filepath.Join(conf.Get(CONF_WORKDIR), DIR_TODO)

	defer os.Remove(tmpFile.Name())

	message := NewMessage()

	if args["--draft"] != nil {
		draft := filepath.Join(draftsDir, args["--draft"].(string))
		m, err := NewMessageFromFile(draft)

		if err != nil {
			fmt.Printf("Error while reading draft: %s\n", err.Error())
			os.Exit(1)
		}

		message = &m
	}

	message.Conf.MergeWithDocOptArgs(CMD_USAGE, &args)

	// MergeWithDocOptArgs will also copy --draft and --help over, but we do not want
	// that.
	delete(message.Conf.Data, "draft")
	delete(message.Conf.Data, "help")

	if message.Get("date") == "" {
		// Add tomorrow's date.
		message.Conf.Set("date", time.Now().AddDate(0, 0, 1).Format(DATE_FORMAT))
	}

	if message.Get("subject") == "" {
		message.Conf.Set("subject", "Type subject here")
	}

	if len(message.Body) == 0 {
		message.Body = append(message.Body, "Add message to the world of tomorrow here.")
	}

	message.WriteToFile(tmpFile.Name())

	editor := editor()

	fmt.Printf("Opening %s using %s.\n", tmpFile.Name(), editor)

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	if cmd.ProcessState.Success() {
		fmt.Printf("\nSave message? ([(y)es], (r)enamed, (d)raft, (n)o): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')

		if err != nil {
			fmt.Printf("Error when reading response: %s\n", err.Error())
			os.Exit(1)
		}

		var dst string

		saved := false

		response = strings.TrimSpace(response)

		switch response {
		case "y", "yes", "":
			dst = filepath.Join(todoDir, filepath.Base(tmpFile.Name())+".msg")
			dst, err = copyFile(tmpFile.Name(), dst, false)
			saved = true
		case "r", "renamed":
			fmt.Printf("\nSpecify new name: ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')

			if err != nil {
				fmt.Printf("Error when reading response: %s\n", err.Error())
				os.Exit(1)
			}

			dst = filepath.Join(todoDir, strings.TrimSpace(response)+".msg")
			dst, err = copyFile(tmpFile.Name(), dst, false)
			saved = true
		case "d", "draft":
			dst = filepath.Join(draftsDir, filepath.Base(tmpFile.Name())+".msg")
			dst, err = copyFile(tmpFile.Name(), dst, false)
			saved = true
		}

		if err != nil {
			fmt.Printf("Error when saving %s: %s\n", dst, err)
		} else if saved {
			fmt.Printf("Saved as: %s\n", dst)
		} else {
			fmt.Println("Discarding message.")
		}
	}
}

// Copy file from `src` to `dst`. If `overwrite` is false, then an alternative
// file name will be used and returned as string.
// TODO: Portability issues (cp) / https://github.com/golang/go/issues/8868
func copyFile(src, dst string, overwrite bool) (string, error) {
	if !overwrite {
		var err error

		dst, err = nextFreeFilename(dst)

		if err != nil {
			return "", err
		}
	}

	err := exec.Command("cp", "-f", src, dst).Run()

	return dst, err
}

func nextFreeFilename(dst string) (string, error) {
	_, e := os.Stat(dst)

	// If the file does not exist, we can use the name
	if os.IsNotExist(e) {
		return dst, nil
	}

	alt := ""

	for i := 0; i < 255; i++ {
		name := fmt.Sprintf("%s.%d", dst, i)
		_, e := os.Stat(name)

		if os.IsNotExist(e) {
			alt = name
			break
		}
	}

	if alt == "" {
		return "", fmt.Errorf("No suitable file name could be found.")
	}

	return alt, nil
}

func editor() string {
	editor := os.Getenv("VISUAL")

	if editor != "" {
		return editor
	}

	editor = os.Getenv("EDITOR")

	if editor != "" {
		return editor
	}

	return "vi"
}
