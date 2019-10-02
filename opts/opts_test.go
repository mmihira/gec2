package opts

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TaggedKeyName", func() {
	It("Should be Name", func() {
		Expect(TaggedKeyName()).To(Equal("Name"))
	})
})
