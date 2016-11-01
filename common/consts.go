/* consts.go: global constants
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
package common

// Project wide constants.
const (
	TIME_FORMAT     = "15:04"
	DATE_FORMAT     = "2006-01-02"
	DATETIME_FORMAT = DATE_FORMAT + " " + TIME_FORMAT

	DIR_TODO   = "todo"
	DIR_DRAFTS = "drafts"
	DIR_ERRORS = "errors"
	DIR_DONE   = "done"

	CMD_USAGE = "usage"
	CMD_RUN   = "run"

	CONF_WORKDIR         = "workdir"
	CONF_DATE            = "date"
	CONF_DAYS            = "days"
	CONF_SUBJECT         = "subject"
	CONF_TO              = "to"
	CONF_FROM            = "from"
	CONF_REPLY_TO        = "reply-to"
	CONF_CC              = "cc"
	CONF_BCC             = "bcc"
	CONF_CONFIG_FILENAME = "config"
	CONF_SMTP_SERVER     = "server"
	CONF_SMTP_PORT       = "port"
	CONF_SMTP_INSECURE   = "insecure"
	CONF_NOT_BEFORE      = "not-before"
	CONF_NOT_AFTER       = "not-after"
)
