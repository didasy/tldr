package tldr_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTldr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tldr Suite")
}
