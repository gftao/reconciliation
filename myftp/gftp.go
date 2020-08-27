package myftp

import (
	"github.com/jlaffaye/ftp"
	"time"
	"fmt"
	"bytes"
	"golib/modules/logr"
	"gopkg.in/dutchcoders/goftp.v1"
	"crypto/tls"
	"runtime/debug"
)

func Myftp(user, password, host, port, fileName, rmtDir string, fileDta []byte) error {
	addr := ""
	if port != "" {
		addr = fmt.Sprintf("%s:%s", host, port)
	} else {
		addr = fmt.Sprintf("%s:%s", host, "21")
	}
	logr.Info("remote addr: ", addr)

	s, err := ftp.Dial(addr, ftp.DialWithTimeout(50*time.Second))
	if err != nil {
 		return err
	}
	defer s.Quit()
	err = s.Login(user, password)
	if err != nil {
		fmt.Println("login")
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
			debug.PrintStack()
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


func Myftp2(user, password, host, port, fileName, rmtDir string, fileDta []byte) error {
	var err error
	var ftp *goftp.FTP
	addr := ""
	if port != "" {
		addr = fmt.Sprintf("%s:%s", host, port)
	} else {
		addr = fmt.Sprintf("%s:%s", host, "21")
	}
	logr.Info("remote addr: ", addr)
	// For debug messages: goftp.ConnectDbg("ftp.server.com:21")
	if ftp, err = goftp.Connect(addr); err != nil {
		logr.Error(err)
		return err
	}

	defer ftp.Close()
	//fmt.Println("Successfully connected to", addr)

	// TLS client authentication
	//config := &tls.Config{
	//	InsecureSkipVerify: true,
	//	ClientAuth:         tls.RequestClientCert,
	//}
 	//if err = ftp.AuthTLS(config); err != nil {
	//	panic(err)
	//}
 	// Username / password authentication
  	if err = ftp.Login(user, password); err != nil {
		logr.Error(err)
		return err
	}

	if err = ftp.Cwd(rmtDir); err != nil {
		logr.Error(err)
		return err
	}
	b := bytes.NewBuffer(fileDta)

	if err := ftp.Stor(fileName, b); err != nil {
		logr.Error(err)
		return err
	}
	logr.Info("ftp post success")
	return nil
}