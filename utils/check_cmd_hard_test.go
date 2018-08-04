package utils_test

import (
	"github.com/ironstar-io/tokaido/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckCmdHard", func() {
	It("should return an error when the program passed is not available on $PATH", func() {
		result := utils.CheckCmdHard("not-an-installed-program")
		Expect(result).To(Equal("log.Fatal called"))
	})

	It("should return an empty string (success) when the program passed is available on $PATH", func() {
		result := utils.CheckCmdHard("date")
		Expect(result).To(Equal(""))
	})
})
