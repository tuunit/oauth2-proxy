package persistence

import (
	"testing"

	logger "github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger/legacy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPersistenceSuite(t *testing.T) {
	logger.SetOutput(GinkgoWriter)
	logger.SetErrOutput(GinkgoWriter)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Persistence")
}
