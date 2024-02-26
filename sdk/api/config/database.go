/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package sdkcfg

// Database is the database configuration.
type Database struct {
	// Host is the hostname or IP address of the database server.
	Host string

	// Username is the username to use when connecting to the database.
	Username string

	// Password is the password to use when connecting to the database.
	Password string

	// Database is the name of the database to use.
	Database string
}

// DbCfg is the interface for database configuration.
type DbCfg interface {
	// Read reads the database configuration.
	Read() (Database, error)

	// Write writes the database configuration.
	Write(Database) error
}