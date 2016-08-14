package rex_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRux(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rux Suite")
}
