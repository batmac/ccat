// Code generated by 'yaegi extract os/user'. DO NOT EDIT.

//go:build go1.21 && !go1.22
// +build go1.21,!go1.22

package stdlib

import (
	"os/user"
	"reflect"
)

func init() {
	Symbols["os/user/user"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Current":       reflect.ValueOf(user.Current),
		"Lookup":        reflect.ValueOf(user.Lookup),
		"LookupGroup":   reflect.ValueOf(user.LookupGroup),
		"LookupGroupId": reflect.ValueOf(user.LookupGroupId),
		"LookupId":      reflect.ValueOf(user.LookupId),

		// type definitions
		"Group":               reflect.ValueOf((*user.Group)(nil)),
		"UnknownGroupError":   reflect.ValueOf((*user.UnknownGroupError)(nil)),
		"UnknownGroupIdError": reflect.ValueOf((*user.UnknownGroupIdError)(nil)),
		"UnknownUserError":    reflect.ValueOf((*user.UnknownUserError)(nil)),
		"UnknownUserIdError":  reflect.ValueOf((*user.UnknownUserIdError)(nil)),
		"User":                reflect.ValueOf((*user.User)(nil)),
	}
}
