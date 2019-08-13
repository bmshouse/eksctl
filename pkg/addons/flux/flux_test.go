package flux_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/eksctl/pkg/utils/random"
)

const (
	namespace         = "flux"
	repository        = "git@github.com:eksctl-bot/my-gitops-repo.git"
	email             = "eksctl-bot@weave.works"
	privateSSHKeyPath = "/Users/marc/.ssh/eksctl-bot_id_rsa"
	name              = "flux"
	region            = "ap-northeast-1"
)

func TestFlux(t *testing.T) {
	// Use a random branch to ensure test runs don't step on each others.
	branch := random.String(8)
	cloneDir, err := createBranch(branch)
	assert.NoError(t, err)
	defer deleteBranch(branch, cloneDir)

	assertFluxManifestsAbsentInGit(t, branch)
	assertFluxPodsAbsentInKubernetes(t)

	eksctl := exec.Command("../../../eksctl", "install", "flux",
		"--git-url", repository,
		"--git-email", email,
		"--git-private-ssh-key-path", privateSSHKeyPath,
		"--git-branch", branch,
		"--name", name,
		"--region", region)
	eksctl.Env = append(os.Environ(), "EKSCTL_EXPERIMENTAL=true", "AWS_PROFILE=default-mfa")
	eksctl.Stdout = os.Stdout
	eksctl.Stderr = os.Stderr
	err = eksctl.Run()
	assert.NoError(t, err)

	assertFluxManifestsPresentInGit(t, branch)
	assertFluxPodsPresentInKubernetes(t)
}
