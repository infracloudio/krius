package krius_cli_tests

import (
	"os/exec"
	"strings"
	"testing"
)

var PATH_TO_KRIUS = "./../../bin/krius"

func TestInstall(t *testing.T) {
	cmds := []string{"install", "bitnami/kube-prometheus"}

	t.Logf("Running '%v'", "krius "+strings.Join(cmds, " "))
	install := exec.Command("krius", cmds...)

	out, err := install.CombinedOutput()
	if err != nil {
		t.Logf(string(out))
		t.Fatal(err)
	}
}
