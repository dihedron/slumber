package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsAnyEditorActive(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "proc_mock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a dummy PID dir with an editor process
	pidDir := filepath.Join(tempDir, "123")
	if err := os.MkdirAll(pidDir, 0755); err != nil {
		t.Fatal(err)
	}
	cmdline := "node\x00/home/user/.vscode-server/bin/some-id/out/server-main.js"
	if err := os.WriteFile(filepath.Join(pidDir, "cmdline"), []byte(cmdline), 0644); err != nil {
		t.Fatal(err)
	}

	// Mock SSH connections (Empty by default)
	netDir := filepath.Join(tempDir, "net")
	if err := os.MkdirAll(netDir, 0755); err != nil {
		t.Fatal(err)
	}
	tcpFile := filepath.Join(netDir, "tcp")
	os.WriteFile(tcpFile, []byte("  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retr   timeout inode\n"), 0644)
	SetNetworkPaths(tcpFile, "/dev/null")

	// Test negative (no SSH)
	editors := IsAnyEditorActive(tempDir)
	if len(editors) > 0 {
		t.Error("expected no active editors when no SSH connections")
	}

	// Mock active SSH connection
	sshLine := "   0: 00000000:0016 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 14467 1 0000000000000000 100 0 0 10 -1\n"
	os.WriteFile(tcpFile, []byte("  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retr   timeout inode\n"+sshLine), 0644)

	// Test positive (with SSH)
	editors = IsAnyEditorActive(tempDir)
	if len(editors) == 0 {
		t.Error("expected active editors when SSH connection is present")
	}
	if len(editors) != 1 || editors[0] != "vscode-server" {
		t.Errorf("expected [vscode-server], got %v", editors)
	}

	// Test false positive: flag value
	cmdline = "myappl\x00--path=vscode-server"
	if err := os.WriteFile(filepath.Join(pidDir, "cmdline"), []byte(cmdline), 0644); err != nil {
		t.Fatal(err)
	}
	editors = IsAnyEditorActive(tempDir)
	if len(editors) > 0 {
		t.Error("expected no active editors for flag value")
	}

	// Test false positive: unrelated command arg
	cmdline = "ls\x00vscode-server"
	if err := os.WriteFile(filepath.Join(pidDir, "cmdline"), []byte(cmdline), 0644); err != nil {
		t.Fatal(err)
	}
	editors = IsAnyEditorActive(tempDir)
	if len(editors) > 0 {
		t.Error("expected no active editors for unrelated command arg")
	}

	// Test positive: interpreter script
	cmdline = "node\x00/usr/bin/vscode-server/server.js"
	if err := os.WriteFile(filepath.Join(pidDir, "cmdline"), []byte(cmdline), 0644); err != nil {
		t.Fatal(err)
	}
	editors = IsAnyEditorActive(tempDir)
	if len(editors) == 0 {
		t.Error("expected active editors for node script")
	}
	if len(editors) != 1 || editors[0] != "vscode-server" {
		t.Errorf("expected [vscode-server], got %v", editors)
	}
}
