package sftp_client

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path/filepath"
)

type sftpClient struct {
	sshConn  *ssh.Client
	sftpConn *sftp.Client
	verbose  bool
}

func (c *sftpClient) Close() {
	c.sshConn.Close()
	c.sftpConn.Close()
}

func (c *sftpClient) Create(path string) (*sftp.File, error) {
	return c.sftpConn.Create(path)
}

func (c *sftpClient) PutFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := c.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		log.Fatalf("put file failed: %v", err)
		return err
	}

	return err
}

func (c *sftpClient) Mkdir(p string) error {
	return c.sftpConn.Mkdir(p)
}

func (c *sftpClient) PutDir(srcDir, dstDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		relFileName, _ := filepath.Rel(srcDir, path)
		dstPath := filepath.Join(dstDir, relFileName)
		dstPathDir, _ := filepath.Split(dstPath)
		dstFileInfo, err := c.sftpConn.Lstat(dstPathDir)
		if dstFileInfo == nil {
			c.Mkdir(dstPathDir)
		}

		err1 := c.PutFile(path, dstPath)
		if err == nil && c.verbose {
			log.Printf("put %v to %v success", path, dstPath)
		}
		return err1
	})

	return err
}

func NewSftpClient(host, user, passwd string, verbose bool) (*sftpClient, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
	}

	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatalf("connect to host failed: %v", err)
	}

	sftp, err1 := sftp.NewClient(conn)
	if err1 != nil {
		log.Fatalf("sftp failed: %v", err1)
	}

	s := &sftpClient{
		conn,
		sftp,
		verbose,
	}

	return s, err1
}
