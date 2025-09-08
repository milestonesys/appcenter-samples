package vms

import (
	"apigateway-webserver/src/pkg/constants/enums"
)

type User struct {
	username            string
	password            string
	credentialsFlowType enums.CredentialsFlowType
}

func NewUser(username, password string, credentialsFlowType enums.CredentialsFlowType) *User {
	return &User{
		username:            username,
		password:            password,
		credentialsFlowType: credentialsFlowType,
	}
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Password() string {
	return u.password
}

func (u *User) CredentialsFlowType() enums.CredentialsFlowType {
	return u.credentialsFlowType
}
