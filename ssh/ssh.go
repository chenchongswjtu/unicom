package ssh

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"

	"golang.org/x/crypto/ssh"
)

// 这种方式第一次连接一个主机时需要输入yes/no
func runCmd() {
	var stdOut, stdErr bytes.Buffer

	cmd := exec.Command("ssh", "root@127.0.0.1", "ls")
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	if err := cmd.Run(); err != nil {
		panic(fmt.Sprintf("cmd exec failed: %s : %s", fmt.Sprint(err), stdErr.String()))
	}

	fmt.Print(stdOut.String())
	fmt.Printf(stdErr.String())
}

func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: hostKeyCallback,
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

// 使用ssh的方式远程执行命令，需要指导用户名，密码，主机ip，端口
func runSSH() {
	var stdOut, stdErr bytes.Buffer

	session, err := SSHConnect("root", "root", "127.0.0.1", 22)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	session.Stdout = &stdOut
	session.Stderr = &stdErr

	err = session.Run("ls") // ssh执行的命令
	if err != nil {
		fmt.Println("ssh run err:", err)
		return
	}

	fmt.Println(stdOut.String())
	fmt.Println(stdErr.String())
}
