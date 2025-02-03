package oidc

import (
	"testing"

	logger "github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger/legacy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOIDCSuite(t *testing.T) {
	logger.SetOutput(GinkgoWriter)

	RegisterFailHandler(Fail)
	RunSpecs(t, "OIDC")
}
