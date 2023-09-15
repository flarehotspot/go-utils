package accounts

import (
	"os"
	"path/filepath"

	"github.com/flarehotspot/core/sdk/libs/yaml-3"
	"github.com/flarehotspot/core/sdk/utils/paths"
	"github.com/flarehotspot/core/sdk/utils/sse"
)

var (
	AcctDir = filepath.Join(paths.ConfigDir, "accounts")
)

type Account struct {
	Uname  string   `yaml:"username"`
	Passwd string   `yaml:"password"`
	Perms  []string `yaml:"permissions"`
}

// return the file path of yaml file
func (acct *Account) YamlFile() string {
	return FilepathForUser(acct.Uname)
}

// Username returns the username of the account
func (acct *Account) Username() string {
	return acct.Uname
}

// Auth returns true if the password is correct
func (acct *Account) Auth(pw string) bool {
	return acct.Passwd == pw
}

// Permissions returns the permissions of the account
func (acct *Account) Permissions() []string {
	return acct.Perms
}

// IsAdmin returns true if the account is admin
func (acct *Account) IsAdmin() bool {
	for _, p := range acct.Perms {
		if p == PermMngUsers {
			return true
		}
	}
	return false
}

// AddSocket adds a sse socket to the account
func (acct *Account) AddSocket(s *sse.SseSocket) {
	sse.AddSocket(acct.Username(), s)
}

// Emit emits an event to the account that will propage to the browser
func (acct *Account) Emit(event string, data interface{}) {
	sse.Emit(acct.Username(), event, data)
}

// Save saves the account to yaml file
func (acct *Account) Save() error {
	b, err := yaml.Marshal(acct)
	if err != nil {
		return err
	}
	return os.WriteFile(acct.YamlFile(), b, 0644)
}

// Update updates the account with new username, password and permissions
func (acct *Account) Update(uname string, pass string, perms []string) error {
	_, err := Update(acct.Uname, uname, pass, perms)
	if err != nil {
		return err
	}

	acct.Uname = uname
	acct.Passwd = pass
	acct.Perms = perms
	return nil
}

// Delete deletes the account
func (acct *Account) Delete() error {
	return Delete(acct.Uname)
}