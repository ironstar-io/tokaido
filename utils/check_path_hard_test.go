package utils_test

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckPathHard", func() {
	It("should return an error when the program passed is not available on $PATH", func() {
		result := utils.CheckPathHard("not-an-installed-program")
		Expect(result).To(Equal("utils.FatalError called"))
	})

	It("should return an empty string (success) when the program passed is available on $PATH", func() {
		result := utils.CheckPathHard("date")
		Expect(result).To(Equal(""))
	})
})
