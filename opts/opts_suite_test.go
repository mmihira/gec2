package opts

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOpts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Opts Suite")
}
