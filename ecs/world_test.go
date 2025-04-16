package ecs_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("World", func() {
	It("should say hello", func() {
		Expect("Hello, World!").To(Equal("Hello, World!"))
	})
})