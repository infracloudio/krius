package helm

import "os/exec"

func HelmRepoAdd(name, url string) {
	exec.Command("helm", "repo", "add", name, url)
}
