package bootstrap

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func Bootstrap(target string) error {
	// target format: user@host[:port]
	parts := strings.Split(target, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid target format, expected user@host")
	}
	user := parts[0]
	host := parts[1]
	if !strings.Contains(host, ":") {
		host = host + ":22"
	}

	fmt.Printf("Bootstrapping %s...\n", target)

	// 1. SSH Client Setup (assume ssh-agent or id_rsa)
	client, err := dialSSH(user, host)
	if err != nil {
		return fmt.Errorf("failed to connect via SSH: %w", err)
	}
	defer client.Close()

	// 2. Cross-compile op-server for Linux/AMD64 (most common)
	fmt.Println("Building op-server for linux/amd64...")
	tmpDir, _ := os.MkdirTemp("", "op-simple-*")
	defer os.RemoveAll(tmpDir)
	serverBin := filepath.Join(tmpDir, "op-server")

	cmd := exec.Command("go", "build", "-o", serverBin, "cmd/op-server/main.go")
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to compile op-server: %w", err)
	}

	// 3. Upload binary
	fmt.Println("Uploading op-server...")
	if err := uploadFile(client, serverBin, "/tmp/op-server"); err != nil {
		return fmt.Errorf("failed to upload binary: %w", err)
	}

	// 4. Run server and capture config
	fmt.Println("Initializing server on remote...")
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Ensure remote dir exists and move binary
	remoteCmd := "chmod +x /tmp/op-server && /tmp/op-server > /tmp/op-server.log 2>&1 & sleep 2 && cat ~/.op-simple/server.yaml"
	output, err := session.CombinedOutput(remoteCmd)
	if err != nil {
		return fmt.Errorf("failed to start remote server: %w. Output: %s", err, string(output))
	}

	fmt.Println("Bootstrap complete.")
	fmt.Printf("Remote config:\n%s\n", string(output))

	return nil
}

func dialSSH(user, host string) (*ssh.Client, error) {
	var auths []ssh.AuthMethod

	// 1. Try SSH Agent
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}

	// 2. Try default keys
	home, _ := os.UserHomeDir()
	for _, keyName := range []string{"id_rsa", "id_ed25519"} {
		keyPath := filepath.Join(home, ".ssh", keyName)
		if key, err := os.ReadFile(keyPath); err == nil {
			if signer, err := ssh.ParsePrivateKey(key); err == nil {
				auths = append(auths, ssh.PublicKeys(signer))
			}
		}
	}

	if len(auths) == 0 {
		return nil, fmt.Errorf("no ssh auth methods found")
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return ssh.Dial("tcp", host, config)
}

func uploadFile(client *ssh.Client, localPath, remotePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, _ := f.Stat()

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C%04o %d %s\n", 0755, stat.Size(), filepath.Base(remotePath))
		io.Copy(w, f)
		fmt.Fprint(w, "\x00")
	}()

	return session.Run("/usr/bin/scp -t " + filepath.Dir(remotePath))
}
