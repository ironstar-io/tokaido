package utils_test

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StdoutCmd", func() {
	It("should return an error when the program passed is not available on $PATH", func() {
		result := utils.StdoutCmd("not-an-installed-program")
		Expect(result).To(Equal("utils.FatalError called"))
	})

	It("should return an error when the program passed is available on $PATH, but args are not valid", func() {
		result := utils.StdoutCmd("date", "-notvalid")
		Expect(result).To(Equal("utils.FatalError called"))
	})

	It("should return an empty string (success) when the program passed is available on $PATH", func() {
		result := utils.StdoutCmd("date")
		Expect(result).To(Equal(""))
	})
})
