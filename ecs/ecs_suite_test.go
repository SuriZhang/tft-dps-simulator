package ecs_test

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestEcs(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Ecs Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	// Initialization logic for the test suite
})

var _ = ginkgo.AfterSuite(func() {
	// Teardown logic for the test suite
})
