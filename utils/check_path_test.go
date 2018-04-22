package utils_test

import (
	"./"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckPath", func() {
	It("should return an error when the program passed is not available on $PATH", func() {
		result := utils.CheckPath("not-an-installed-program")
		Expect(result).To(Equal("utils.FatalError called"))
	})

	It("should return an empty string (success) when the program passed is available on $PATH", func() {
		result := utils.CheckPath("date")
		Expect(result).To(Equal(""))
	})
})
