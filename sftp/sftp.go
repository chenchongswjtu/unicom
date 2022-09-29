package sftp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const bufSize = 1 << 20 //1M

var ErrSftpClientIsNil = errors.New("sftp client is nil")

func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	//hostKey := getHostKey(host)
	//config := ssh.ClientConfig{
	//	User: user,
	//	Auth: auths,
	//	// Uncomment to ignore host key check
	//	//HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	//	HostKeyCallback: ssh.FixedHostKey(hostKey),
	//}

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: hostKeyCallback,
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func uploadFile(sc *sftp.Client, localFile, remoteFileDir string) error {
	if sc == nil {
		return ErrSftpClientIsNil
	}

	srcFile, err := os.Open(localFile)
	if err != nil {
		fmt.Printf("Unable to open local file: %v\n\n", err)
		return err
	}
	defer srcFile.Close()

	fileName := filepath.Base(localFile) // 获得文件名
	err = sc.MkdirAll(remoteFileDir)     // 远程创建目录
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Note: SFTP To Go doesn't support O_RDWR mode
	dstFile, err := sc.OpenFile(filepath.Join(remoteFileDir, fileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		fmt.Printf("Unable to open remote file: %v\n", err)
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, bufSize)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dstFile.Write(buf[:n]); err != nil {
			return err
		}
	}

	fmt.Printf("upload file [%s] to [%s] done.\n", localFile, filepath.Join(remoteFileDir, fileName))
	return err
}

func downloadFile(sc *sftp.Client, remoteFile, localFileDir string) error {
	if sc == nil {
		return ErrSftpClientIsNil
	}

	// Note: SFTP To Go doesn't support O_RDWR mode
	srcFile, err := sc.OpenFile(remoteFile, os.O_RDONLY)
	if err != nil {
		fmt.Printf("Unable to open remote file: %v\n", err)
		return err
	}
	defer srcFile.Close()

	fileName := filepath.Base(remoteFile)
	err = sc.MkdirAll(localFileDir)

	dstFile, err := os.Create(filepath.Join(localFileDir, fileName))
	if err != nil {
		fmt.Printf("Unable to open local file: %v\n", err)
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, bufSize)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dstFile.Write(buf[:n]); err != nil {
			return err
		}
	}

	fmt.Printf("download file [%s] to [%s] done.\n", remoteFile, filepath.Join(localFileDir, fileName))
	return nil
}

// Get host key from local known hosts
func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read known_hosts file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing %q: %v\n", fields[2], err)
				os.Exit(1)
			}
			break
		}
	}

	if hostKey == nil {
		fmt.Fprintf(os.Stderr, "No hostkey found for %s", host)
		os.Exit(1)
	}

	return hostKey
}
