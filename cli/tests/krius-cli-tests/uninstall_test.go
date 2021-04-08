package krius_cli_tests

import (
	"os/exec"
	"strings"
	"testing"
)

func TestUninstall(t *testing.T) {
	cmds := []string{"uninstall", "prom"}

	t.Logf("Running '%v'", "krius "+strings.Join(cmds, " "))
	install := exec.Command("krius", cmds...)

	out, err := install.CombinedOutput()
	if err != nil {
		t.Logf(string(out))
		t.Fatal(err)
	}
}
