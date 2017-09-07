package myftp

import (
	"github.com/jlaffaye/ftp"
	"time"
	"fmt"
	"bytes"
	"golib/modules/logr"
	"gopkg.in/dutchcoders/goftp.v1"
	"crypto/tls"
)

func Myftp(user, password, host, port, fileName, rmtDir string, fileDta []byte) error {
	addr := ""
	if port != "" {
		addr = fmt.Sprintf("%s:%s", host, port)
	} else {
		addr = fmt.Sprintf("%s:%s", host, "21")
	}
	logr.Info("remote addr: ", addr)
	s, err := ftp.DialTimeout(addr, 20*time.Second)
	if err != nil {
		return err
	}
	defer s.Quit()
	err = s.Login(user, password)
	if err != nil {
		return err
	}
	defer s.Logout()
	if rmtDir != "" {
		err = s.ChangeDir(rmtDir)
		if err != nil {
			err = s.MakeDir(rmtDir)
			if err != nil {
				return err
			}
			s.ChangeDir(rmtDir)
		}
	}
	b := bytes.NewBuffer(fileDta)
	err = s.Stor(fileName, b)
	if err != nil {
		return err
	}
	logr.Info("ftp post success")
	return nil
}

func MyftpTSL(user, password, host, port, fileName, rmtDir string, fileDta []byte) error {
	defer func() {
		if r := recover(); r != nil {
			logr.Info("recover:", r)
		}
	}()
	var ftp *goftp.FTP
	var err error
	addr := ""
	if port != "" {
		addr = fmt.Sprintf("%s:%s", host, port)
	} else {
		addr = fmt.Sprintf("%s:%s", host, "21")
	}
	logr.Info("remote addr: ", addr)
	ftp, err = goftp.Connect(addr)
	if err != nil {
		logr.Info("Connect err:", err)
		return err
	}
	//defer ftp.Close()

	config := &tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequestClientCert,
	}
	err = ftp.AuthTLS(config)
	if err != nil {
		logr.Info("AuthTLS err:", err)
		return err
	}
	logr.Info("user: ", user, "---", "password:", password)

	err = ftp.Login(user, password)
	if err != nil {
		logr.Info("Login err:", err)
		return err
	}
	defer ftp.Quit()

	if rmtDir != "" {
		err = ftp.Cwd(rmtDir)
		if err != nil {
			err = ftp.Mkd(rmtDir)
			if err != nil {
				return err
			}
			ftp.Cwd(rmtDir)
		}
	}
	p, _ := ftp.Pwd()
	logr.Info("pwd:", p)
	b := bytes.NewReader(fileDta)
	err = ftp.Stor(fileName, b)
	if err != nil {
		return err
	}
	logr.Info("ftp post success")

	return nil
}
