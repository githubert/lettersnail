/* run.go: send pending messages
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
	"crypto/tls"
	"fmt"
	"github.com/docopt/docopt.go"
	. "github.com/githubert/clockrotz/common"
	"github.com/jordan-wright/email"
	"net/mail"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var usageRun =
// tag::run[]
`
Usage:
  clockrotz run [options]

Options:
  --help             Show this help.
  --to=ADDR          Destination address.
  --from=ADDR        Sender address.
  --subject=ADDR     Short subject.
  --cc=ADDR          Set "Cc".
  --bcc=ADDR         Set "Bcc".
  --reply-to=ADDR    Set "Reply-To".
  --not-before=TIME  Not before TIME. (default: 00:00)
  --not-after=TIME   Not after TIME. (default: 23:59)
  --server=HOST      SMTP hostname. (default: localhost)
  --port=PORT        SMTP port. (default: 587)
  --verbose          Report on successfully sent messages.
  --dry-run          Do not send the message.
  --insecure         Accept any TLS certificate.
` // end::run[]

func Run(argv []string, conf *Configuration) {
	args, _ := docopt.Parse(usageRun, argv, true, "", false)
	conf.MergeWithDocOptArgs(CMD_RUN, &args)

	verbose := args["--verbose"].(bool)
	dryRun := args["--dry-run"].(bool)
	insecure := args["--insecure"].(bool)

	now := time.Now()

	notBeforeTime, err := time.Parse(TIME_FORMAT, conf.Get(CONF_NOT_BEFORE))

	if err != nil {
		fmt.Printf("Failed parsing not-before time: %s\n", err.Error())
		return
	}

	notAfterTime, err := time.Parse(TIME_FORMAT, conf.Get(CONF_NOT_AFTER))

	if err != nil {
		fmt.Printf("Failed parsing not-after time: %s\n", err.Error())
		return
	}

	notBefore := buildTime(now, notBeforeTime.Hour(), notBeforeTime.Minute(), true)
	notAfter := buildTime(now, notAfterTime.Hour(), notAfterTime.Minute(), false)

	// Return if we are in some quiet period.
	if now.After(notAfter) || now.Before(notBefore) {
		return
	}

	messages := NewMessagesFromDirectory(filepath.Join(conf.Get(CONF_WORKDIR), DIR_TODO))

	verificationError := false

	for _, message := range messages {
		message.Conf.MergeWith(conf)
		err := processMessage(message, now, dryRun, insecure, verbose)

		if err != nil {
			verificationError = true
			continue
		}
	}

	if verificationError {
		// FIXME: This suggests that other messages were not sent, but they were....
		fmt.Println("There were errors when verifying one or more messages.")
		fmt.Println("Please run 'clockrotz check'")
		os.Exit(1)
	}
}

func processMessage(message Message, now time.Time, dryRun, insecure, verbose bool) error {
	if errs := message.Verify(); errs != nil {
		return fmt.Errorf("Message %s failed verification.", message.Name)
	}

	date, err := ParseTime(message.Get("date"))

	if err != nil {
		return err
	}

	if !now.After(date) {
		return nil
	}

	sendErr := sendMessage(message, dryRun, insecure)

	if sendErr != nil {
		if !dryRun {
			err := moveMessage(message, DIR_ERRORS)

			if err != nil {
				fmt.Printf("Error when moving message %s: %s\n", message.Name, err.Error())
			}

			logMessage(message, DIR_ERRORS, sendErr.Error())
		}

		fmt.Printf("Error when sending message %s: %s\n", message.Name, sendErr.Error())
	} else {
		if !dryRun {
			err := moveMessage(message, DIR_DONE)

			if err != nil {
				fmt.Printf("Error when moving message %s: %s\n", message.Name, err.Error())
			}

			logMessage(message, DIR_DONE, "Successfully delivered.")
		}

		if verbose {
			fmt.Printf("Message %s delivered.\n", message.Name)
		}
	}

	return nil // TODO: what about errors that are not verification errors?
}

func logMessage(message Message, dir string, logMessage string) {
	dstDir := filepath.Join(message.Get(CONF_WORKDIR), dir)
	filename := filepath.Join(dstDir, message.Name[:len(message.Name)-len(".msg")]+".log")

	f, err := os.Create(filename)

	if err != nil {
		fmt.Printf("Error when creating log file: %s\n", err.Error())
		return
	}

	defer f.Close()

	f.WriteString("Log message:\n")
	f.WriteString(fmt.Sprintf("  %s\n", logMessage))
	f.WriteString("\n")

	f.WriteString("\nConfiguration:\n")
	for _, s := range message.Conf.DumpConfig() {
		f.WriteString(fmt.Sprintf("  %s\n", s))
	}

	f.WriteString("\nBody:\n")
	for _, s := range message.Body {
		f.WriteString(fmt.Sprintf("  %s\n", s))
	}
}

// Turn the given plain address list string into an array.
func getAddresses(addressList string) ([]string, error) {
	addresses, err := mail.ParseAddressList(addressList)

	if err != nil {
		return nil, err
	}

	result := []string{}

	for _, address := range addresses {
		result = append(result, address.String())
	}

	return result, nil
}

// Prepare a ready-to-send Email message.
func prepareEmail(message *Message) (*email.Email, error) {
	e := email.NewEmail()
	e.From = message.Get(CONF_FROM)
	e.Subject = message.Get(CONF_SUBJECT)

	// Build list of To addresses.
	to, err := getAddresses(message.Get(CONF_TO))

	if err != nil {
		return nil, err
	}

	e.To = to

	// Set optional Reply-To header.
	if r := message.Get(CONF_REPLY_TO); r != "" {
		replyTo, err := getAddresses(r)

		if err != nil {
			return nil, err
		}

		e.Headers = textproto.MIMEHeader{"Reply-To": replyTo}
	}

	// Build list of Cc addresses.
	if r := message.Get(CONF_CC); r != "" {
		cc, err := getAddresses(r)

		if err != nil {
			return nil, err
		}

		e.Cc = cc
	}

	// Build list of Bcc addresses.
	if r := message.Get(CONF_BCC); r != "" {
		bcc, err := getAddresses(r)

		if err != nil {
			return nil, err
		}

		e.Bcc = bcc
	}

	e.Text = []byte(strings.Join(message.Body, "\n"))

	return e, nil
}

// Send the given message, unless `dryRun` is true. Use `insecure` to work
// around things like self-signed certificates.
func sendMessage(message Message, dryRun bool, insecure bool) error {
	e, err := prepareEmail(&message)

	if err != nil {
		return err
	}

	smtpServer := message.Get(CONF_SMTP_SERVER) + ":" + message.Get(CONF_SMTP_PORT)

	if dryRun {
		fmt.Printf("Skip sending message %s through %s.\n", message.Name, smtpServer)
		return nil
	}

	if insecure {
		return e.SendWithTLS(smtpServer, nil, &tls.Config{InsecureSkipVerify: true})
	} else {
		return e.Send(smtpServer, nil)
	}
}

// Move the given message to a folder relative to the working directory.
func moveMessage(message Message, relative string) error {
	todo := filepath.Join(message.Get(CONF_WORKDIR), DIR_TODO)
	to := filepath.Join(message.Get(CONF_WORKDIR), relative)

	// FIXME: BUG: This will overwrite existing messages. Look at create.go:nextFreeFilename()
	//             for ideas on how to resolve this. We could try to use a similar approach.
	//             `foo.msg` to `foo.1.msg` and `foo.1.log`.

	return os.Rename(filepath.Join(todo, message.Name), filepath.Join(to, message.Name))
}

// Build a time.Time from some given base time. If floor is true, seconds will
// be set to 0, if false, 59.
// TODO: Maybe there is a better way?…
func buildTime(base time.Time, hour int, minute int, floor bool) time.Time {
	seconds := 0

	if !floor {
		// We ignore the possibility of leap seconds here, this is
		// just done so that we can get 23:59:59 instead of 23:59:00.
		// 23:59:00 would make us miss a whole minute. It would be
		// nice if there were a way to indicate to time.Date() to
		// build a date with start of the day / end of the day...
		seconds = 59
	}

	return time.Date(
		base.Year(),
		base.Month(),
		base.Day(),
		hour,
		minute,
		seconds,
		0,
		base.Location())
}
