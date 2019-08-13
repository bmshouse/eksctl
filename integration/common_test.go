// +build integration

package integration_test

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type tHelper struct{ GinkgoTInterface }

func (t *tHelper) Helper()      { return }
func (t *tHelper) Name() string { return "eksctl-test" }

func newKubeTest() (*harness.Test, error) {
	t := &tHelper{GinkgoT()}
	l := &logger.TestLogger{}
	h := harness.New(harness.Options{Logger: l.ForTest(t)})
	if err := h.Setup(); err != nil {
		return nil, err
	}
	if err := h.SetKubeconfig(kubeconfigPath); err != nil {
		return nil, err
	}
	test := h.NewTest(t)
	test.Setup()
	return test, nil
}

type params struct {
	Args []string
	Env  []string
}

func eksctl(params params) *gexec.Session {
	command := exec.Command(eksctlPath, params.Args...)
	command.Env = append(os.Environ(), params.Env...)
	fmt.Fprintf(GinkgoWriter, "calling %q with %v\n", eksctlPath, params.Args)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	if err != nil {
		Fail(fmt.Sprintf("error starting process: %v\n", err), 1)
	}

	t := time.Minute
	switch params.Args[0] {
	case "create":
		t *= 25
	case "delete":
		t *= 15
	case "get":
		t *= 1
	case "scale":
		t *= 5
	default:
		t *= 30
	}
	session.Wait(t)
	return session
}

func eksctlSuccess(args ...string) *gexec.Session {
	return eksctlSuccessWith(params{Args: args})
}

func eksctlSuccessWith(params params) *gexec.Session {
	session := eksctl(params)
	Expect(session.ExitCode()).To(Equal(0))
	return session
}

func eksctlFail(args ...string) *gexec.Session {
	session := eksctl(params{Args: args})
	Expect(session.ExitCode()).To(Not(Equal(0)))
	return session
}
