/* message.go: module for managing message information
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
	"bufio"
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
)

type Message struct {
	Conf Configuration
	Body []string
	Name string
}

// Supporting sort.Interface.
type Messages []Message

// Load all messages in the given directory.
func NewMessagesFromDirectory(dir string) Messages {
	messages := []Message{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only interested in regular files
		if !info.Mode().IsRegular() {
			return nil
		}

		// Parse only .msg files
		if filepath.Ext(info.Name()) != ".msg" {
			return nil
		}

		message, err := NewMessageFromFile(path)

		if err != nil {
			return err
		}

		messages = append(messages, message)

		return nil
	})

	// TODO: Maybe return error?
	if err != nil {
		fmt.Printf("Error while reading messages from %s:\n\t%s\n", dir, err.Error())
	}

	return messages
}

func (m Messages) Len() int      { return len(m) }
func (m Messages) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m Messages) Less(i, j int) bool {
	// TODO: We are being very confident here w.r.t. errors
	t1, _ := ParseTime(m[i].Get(CONF_DATE))
	t2, _ := ParseTime(m[j].Get(CONF_DATE))

	return t1.Before(t2)
}

// Create a new, empty Message.
func NewMessage() *Message {
	return &Message{
		Conf: *NewConfiguration(),
		Body: []string{},
		Name: "",
	}
}

// Construct new Message from the given file.
func NewMessageFromFile(path string) (Message, error) {
	f, err := os.Open(path)

	if err != nil {
		return Message{}, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	msgConf, msgBody := SplitMessage(lines)

	configuration := Configuration{}
	configuration.Load(msgConf)

	return Message{
		Body: msgBody,
		Conf: configuration,
		Name: filepath.Base(path),
	}, nil
}

// Write a message to a file, such that it could be loaded again.
func (m *Message) WriteToFile(file string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	defer f.Close()

	for _, s := range m.Conf.DumpConfig() {
		f.WriteString(s)
		f.WriteString("\n")
	}

	f.WriteString("\n")

	for _, s := range m.Body {
		f.WriteString(s)
		f.WriteString("\n")
	}

	return nil
}

// Return the specified configuration key's value.
func (m *Message) Get(key string) string {
	return m.Conf.Data[key]
}

// Checks if set and if ParseAddressList succeeds
func verifyAddressList(addresses string) error {
	if addresses == "" {
		// TODO: Why is this part different than the one in verifyAddress()? Maybe not intentional?
		return nil
	}

	_, err := mail.ParseAddressList(addresses)

	if err != nil {
		return err
	}

	return nil
}

// Checks if no nil and if ParseAddress succeeds.
func verifyAddress(paramName string, address string) error {
	if address == "" {
		return fmt.Errorf("'%s' parameter is missing", paramName)
	}

	_, err := mail.ParseAddress(address)

	if err != nil {
		return err
	}

	return nil
}

// Verify if a message has all necessary parameters. We need at least to, from,
// subject, and date. This will also verify optional address lists, etc.
func (m *Message) Verify() []error {
	errors := []error{}

	if err := verifyAddress("from", m.Get("from")); err != nil {
		errors = append(errors, err)
	}

	if err := verifyAddress("to", m.Get("to")); err != nil {
		errors = append(errors, err)
	}

	if m.Get("subject") == "" {
		errors = append(errors, fmt.Errorf("'subject' parameter is missing"))
	}

	if m.Get("date") == "" {
		errors = append(errors, fmt.Errorf("'date' parameter is missing"))
	} else {
		_, err := ParseTime(m.Get("date"))

		if err != nil {
			errors = append(errors, fmt.Errorf("'date' format error: %s", err.Error()))
		}
	}

	if m.Get("reply-to") != "" {
		if err := verifyAddress("nil", m.Get("to")); err != nil {
			errors = append(errors, err)
		}
	}

	if err := verifyAddressList(m.Get("cc")); err != nil {
		errors = append(errors, err)
	}

	if err := verifyAddressList(m.Get("bcc")); err != nil {
		errors = append(errors, err)
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
