package ec2Query

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEc2Query(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ec2Query Suite")
}
