package main

import (
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

const (
	forceDisconnectAfter = time.Millisecond * 50
)

// ValidatEmaileHost validate email host by quering mx records followed by resolving and connecting smtp server
func ValidatEmaileHost(email string) error {
	i := strings.LastIndexByte(email, '@')
	host := email[i+1:]
	mx, err := net.LookupMX(host)
	if err != nil {
		return errors.New("error while fetching MX records")
	}
	if len(mx) == 0 {
		return errors.New("MX records not found")
	}
	addr := fmt.Sprintf("%s:%d", mx[0].Host, 25)
	conn, err := net.DialTimeout("tcp", addr, forceDisconnectAfter)
	t := time.AfterFunc(forceDisconnectAfter, func() { conn.Close() })
	defer t.Stop()

	if err != nil {
		return err
	}
	mxhost, _, _ := net.SplitHostPort(addr)
	client, err := smtp.NewClient(conn, mxhost)
	if err != nil {
		return err
	}
	client.Close()
	return nil
}
