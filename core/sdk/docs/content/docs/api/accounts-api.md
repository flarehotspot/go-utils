+++
title = "Accounts API"
description = "Manage admin accounts and permissions."
date = 2021-05-01T08:00:00+00:00
updated = 2021-05-01T08:00:00+00:00
draft = false
weight = 20
sort_by = "weight"
template = "docs/page.html"

[extra]
lead = "Manage admin accounts and permissions."
toc = true
top = false
+++

# AccountsApi

## Overview
The `AccountsApi` let's you create, modify, remove and manage system admin accounts and permissions.

# Methods
First, get an instance of the `AccountsApi` from the [PluginApi](../plugin-api):
```go
package main
// imports...
func Init(api sdkplugin.PluginApi) {
    accountsApi := api.Acct()
}
```

The following are the available methods in `AccountsApi`.

## Create
It creates a new admin account with the given username, password and [permissions](#permissions). It returns an [Account](#account-instance) instance and an `error` object.
```go
username := "admin"
password := "admin"
permissions := []string{"admin"}
acct, err := accountsApi.Create(username, password, permissions)
if err != nil {
    fmt.Println(err) // Error
}
fmt.Println(acct) // Account
```

## Find
It finds an admin account by the given username. It returns an [Account](#account-instance) instance and an `error` object.
```go
acct, err := accountsApi.Find("admin")
if err != nil {
    fmt.Println(err) // Error
}
fmt.Println(acct) // Account
```

## GetAll
It returns all the admin accounts, admin and non-admin. It returns a slice of [Account](#account-instance) instance and an `error` object.
```go
accts, err := accountsApi.GetAll()
if err != nil {
    fmt.Println(err) // Error
}
fmt.Println(accts) // []Account
```

## GetAdmins
It returns all the admin accounts. It returns a slice of [Account](#account-instance) instance and an `error` object.
```go
accts, err := accountsApi.GetAdmins()
if err != nil {
    fmt.Println(err) // Error
}
fmt.Println(accts) // []Account
```

## NewPerm
It creates a new permission with the given name and description. It returns an `error` object.
```go
name := "newperm"
desc := "New permission"
err := accountsApi.NewPerm(name, desc)
if err != nil {
    fmt.Println(err) // Error
}
```

## GetPerms
It returns all the available permissions, including custom ones from plugins. The return type is `map[string]string` (name and description pairs of permissions).
```go
perms := accountsApi.GetPerms()
fmt.Println(perms) // map[string]string{"admin": "The admin permission"}
```

## PermDesc
Returns the description of the given permission. It returns a `string` and an `error` object.
```go
desc, err := accountsApi.PermDesc("newperm")
if err != nil {
    fmt.Println(err) // Error
}
fmt.Println(desc) // "New permission"
```

---

# Account Instance
Account instance represents a system admin account. First, find an admin account:
```go
acct, err := accountsApi.Find("admin")
if err != nil {
    fmt.Println(err) // Error
}
fmt.Println(acct) // Account
```

Given an admin account instance, you can access the following properties and methods.

## Username
It returns the username of the admin account.
```go
acct.Username() // "admin"
```

## Permissions {#acct-perms}
It returns the [permissions](#permissions) of the admin account.
```go
acct.Permissions() // []string{"admin"}
```

## HasAllPerms
Returns `true` if the admin account has all the given permissions. It can be used to check if an admin account has all the required permissions to access a certain part of the system.
```go
acct, _ := accountsApi.Find("admin")
hasAll := acct.HasAllPerms([]string{"admin"})
fmt.Println(hasAll) // true
```

## HasAnyPerm
It returns `true` if the admin account has any of the given permissions. It can be used to check if an admin account has any of the required permissions to access a certain part of the system.
```go
acct, _ := accountsApi.Find("admin")
hasAny := acct.HasAnyPerm([]string{"admin"})
fmt.Println(hasAny) // true
```

## IsAdmin
It returns `true` if the admin account has the `admin` permission.
```go
acct.IsAdmin() // true
```

## Update
It updates the admin account with the given username, password and [permissions](#permissions). It returns an `error` object.
```go
newUsername := "newadmin"
newPassword := "********"
err := acct.Update(newUsername, newPassword, []string{"admin"})
if err != nil {
    fmt.Println(err) // Error
}
```

## Delete
It deletes the admin account. It returns an `error` object. Note: You cannot delete the last admin account since it is required for the system to function.
```go
err := acct.Delete()
if err != nil {
    fmt.Println(err) // Error
}
```

## Emit
Emit an [event](#events) to the admin account. It returns an `error` object.
```go
evt := "some_event"
data := map[string]any{"key": "value"}
acct, _ := accountsApi.Find("admin")
acct.Emit(evt, data)
```

---

# Permissions
Permissions are used to control the access to various parts of the system. Users without the appropriate permissions will not be able to access the restricted parts of the system.

These are the default permissions that you can assign to an admin account. Although you may define your custom permissions using the [Accounts API](#newperm).

| Permission | Description
| --- | --- |
| `admin` | The admin permission grants full access to the system. |

**TODO:**
- Add more permissions.

---

# Events
Events are emitted to the admin accounts. You can listen to these events and perform certain actions when they are emitted. Here are the available events:

TODO: Add events