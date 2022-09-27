package sftp

import (
	"fmt"
	"testing"

	"github.com/pkg/sftp"
)

func TestUploadFile(t *testing.T) {
	var (
		err        error
		sftpClient *sftp.Client
	)
	sftpClient, err = connect("root", "root", "127.0.0.1", 22)
	if err != nil {
		panic(err)
	}

	err = uploadFile(sftpClient, "/home/chenchong/gopath/src/awesomeProject/test/local/local.txt", "/home/chenchong/gopath/src/awesomeProject/test/remote")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestDownloadFile(t *testing.T) {
	var (
		err        error
		sftpClient *sftp.Client
	)
	sftpClient, err = connect("root", "root", "127.0.0.1", 22)
	if err != nil {
		panic(err)
	}

	err = downloadFile(sftpClient, "/home/chenchong/gopath/src/awesomeProject/test/remote/remote.txt", "/home/chenchong/gopath/src/awesomeProject/test/local")
	if err != nil {
		fmt.Println(err)
		return
	}
}
