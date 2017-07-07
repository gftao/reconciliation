// sftp project doc.go

/*
sftp document
*/
package myfstp

import (
	"fmt"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path"
	"log"
	"strconv"
	"strings"
	"golib/modules/logr"
)

func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	addr = fmt.Sprintf("%s:%d", host, port)
	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, addr net.Addr, key ssh.PublicKey) error {
			return nil
		},
	} // connet to ssh

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func PosSftp(user, password, host, port, filePath, rmtDir string) error {
	var (
		err error
		sftpClient *sftp.Client
	)

	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	p, _ := strconv.Atoi(port)
	sftpClient, err = connect(user, password, host, p)
	if err != nil {
		logr.Info("SSH 连接出错", err)
	}
	defer sftpClient.Close()
	//用来测试的本地文件路径 和 远程机器上的文件夹
	var localFilePath = filePath
	var remoteDir = rmtDir
	localFilePath = strings.Replace(localFilePath, "\\", "//", -1)
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	logr.Info("--->", localFilePath)
	var remoteFileName = path.Base(localFilePath)
	logr.Info("--->", remoteFileName)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		return error(err)
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}
	logr.Info("copy file to remote server finished!")
	return nil
}

func PosByteSftp(user, password, host, port, fileName, rmtDir string, fileData []byte) error {
	var (
		err error
		sftpClient *sftp.Client
	)

	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	p, _ := strconv.Atoi(port)
	sftpClient, err = connect(user, password, host, p)
	if err != nil {
		logr.Info("SSH 连接出错", err)
	}
	defer sftpClient.Close()
	//用来测试的本地文件路径 和 远程机器上的文件夹
	//var localFilePath = filePath
	var remoteDir = rmtDir
	var remoteFileName = fileName
	logr.Info("--->", remoteFileName)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		return error(err)
	}
	defer dstFile.Close()
	l := len(fileData)
	c := 0
	for {
		n, err := dstFile.Write(fileData)

		if err != nil {
			break
		}
		c += n
		if c == l {
			break
		}
	}
	logr.Info("copy file to remote server finished!")
	return nil
}