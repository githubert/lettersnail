/* configuration.go: module for managing message properties
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
	"fmt"
	"github.com/go-ini/ini"
	"sort"
	"strconv"
	"strings"
)

type Configuration struct {
	Data map[string]string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Data: map[string]string{},
	}
}

func (c *Configuration) Get(key string) string {
	return c.Data[key]
}

func (c *Configuration) Set(key string, value string) {
	c.Data[key] = value
}

// Load configuration from an array of strings in the form `key: value`.
func (c *Configuration) Load(text []string) {
	c.Data = map[string]string{}

	for _, line := range text {
		r := strings.SplitN(line, ":", 2)

		key := strings.TrimSpace(r[0])

		if len(r) == 2 {
			c.Data[key] = strings.TrimSpace(r[1])
		} else {
			if key != "" {
				c.Data[key] = ""
			}
		}
	}
}

// Merge the `src` configuration into this configuration.
func (c *Configuration) MergeWith(src *Configuration) {
	for k, v := range (*src).Data {
		(*c).Data[k] = v
	}
}

// Merge arguments from a INI section into the given map.
func mergeIniSection(section *ini.Section, dst *map[string]string) {
	for _, k := range section.KeyStrings() {
		(*dst)[k] = section.Key(k).String()
	}
}

// Merges the "default" and the "cmd" section from CONF_CONFIG_FILENAME.
func (c *Configuration) MergeWithIni(cmd string) {
	cfg, err := ini.Load(c.Get(CONF_CONFIG_FILENAME))

	if err == nil {
		section, err := cfg.GetSection("default")

		if err == nil {
			mergeIniSection(section, &c.Data)
		}

		section, err = cfg.GetSection(cmd)

		if err == nil {
			mergeIniSection(section, &c.Data)
		}
	}
}

// Merge arguments from the DocOpt parser into a configuration map. All
// arguments that are not `nil` and start with "--" will be merged.
// Booleans will be converted to strings.
func (c *Configuration) MergeWithDocOptArgs(cmd string, args *map[string]interface{}) {
	for k, v := range *args {

		// The args list contains the name of the command, but we are not
		// interested in it.
		if k == cmd {
			continue
		}

		// Merge all args that start with -- and where the value is not nil
		if v != nil && (len(k) > 2 && k[0:2] == "--") {
			switch v.(type) {
			case string:
				c.Data[k[2:]] = v.(string)
			case bool:
				c.Data[k[2:]] = strconv.FormatBool(v.(bool))
			}

		}
	}
}

// Dump the configuration as strings in the form `key: value`.
func (c *Configuration) DumpConfig() []string {
	keys := make([]string, len(c.Data))

	i := 0
	for k, _ := range c.Data {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	result := make([]string, len(keys))

	for i, k := range keys {
		result[i] = fmt.Sprintf("%s: %s", k, c.Get(k))
	}

	return result
}
