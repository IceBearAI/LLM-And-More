package ldapcli

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
)

type user struct {
	Displayname    string `json:"displayname,omitempty"`
	Password       string `json:"password,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
	LoginFrequency int    `json:"loginFrequency,omitempty"`
	AuthType       int8   `json:"authType,omitempty"`
	OnlyOss        bool   `json:"onlyOSS"`
}

type Config struct {
	Host         string
	Port         int
	UseSSL       bool
	BindUser     string
	BindPassword string
	BindDN       string
	Attributes   []string
	Filter       string
}

type Middleware func(Service) Service

type Service interface {
	Connection(ctx context.Context) (err error)
	// Authenticate 授权登陆
	Authenticate(ctx context.Context, account string, password string) (bool, error)
}

type service struct {
	conn     *ldap.Conn
	pageSize uint32
	Config
}

func (s *service) Connection(ctx context.Context) (err error) {
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		err = errors.Wrap(err, "ldap.Connection.ldap.Dial")
		return err
	}
	if s.UseSSL {
		if err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			return errors.Wrap(err, "conn.StartTLS")
		}
	}
	s.conn = conn
	return nil
}

func (s *service) Authenticate(ctx context.Context, account string, password string) (res bool, err error) {
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		err = errors.Wrap(err, "ldap.Connection.ldap.Dial")
		return
	}
	if s.UseSSL {
		if err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			return false, errors.Wrap(err, "conn.StartTLS")
		}
	}
	defer conn.Close()

	// First bind with a read only user
	err = conn.Bind(s.BindUser, s.BindPassword)
	if err != nil {
		err = errors.Wrap(err, "conn.Bind")
		return
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		s.BindDN,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf(s.Filter, ldap.EscapeFilter(account)),
		s.Attributes,
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		err = errors.Wrap(err, "conn.Search")
		return
	}

	if len(sr.Entries) != 1 {
		err = errors.New("User does not exist or too many entries returned")
		return
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = conn.Bind(userdn, password)
	if err != nil {
		err = errors.Wrap(err, "conn.Bind: Bind as the user to verify their password")
		return
	}
	//
	// Rebind as the read only user for any further queries
	err = conn.Bind(s.BindUser, s.BindPassword)
	if err != nil {
		err = errors.Wrap(err, "conn.Bind: Rebind as the read only user for any further queries")
		return
	}

	return true, nil
}

func New(cfg Config) Service {
	return &service{
		Config:   cfg,
		pageSize: 1000,
	}
}
