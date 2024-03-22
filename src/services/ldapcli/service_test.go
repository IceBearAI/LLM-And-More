package ldapcli

import (
	"context"
	"fmt"
	"testing"
)

func InitClient() Service {
	return New(Config{
		Host:         "",
		Port:         389,
		UseSSL:       false,
		BindUser:     "",
		BindPassword: "",
		BindDN:       "OU=HABROOT,DC=corp",
		Attributes:   []string{"userPrincipalName", "sAMAccountName", "displayName", "mail", "telephoneNumber", "name"},
		Filter:       "",
	})
}

func TestService_Authenticate(t *testing.T) {
	res, err := InitClient().Authenticate(context.Background(), "admin@admin.com", "!asdf@asdfasdfasd")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(res)
}
