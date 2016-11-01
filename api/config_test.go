package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/starkandwayne/covalence/api"
)

var _ = Describe("API Configuration", func() {
	Describe("Configuration", func() {
		var a *Api

		BeforeEach(func() {
			a = NewApi()
			Ω(a).ShouldNot(BeNil())
		})

		It("handles missing files", func() {
			Ω(a.ReadConfig("/path/to/nowhere")).ShouldNot(Succeed())
		})

		It("handles malformed YAML files", func() {
			Ω(a.ReadConfig("test/etc/config.xml")).ShouldNot(Succeed())
		})

		It("handles YAML files with missing directives", func() {
			Ω(a.ReadConfig("test/etc/empty.yml")).Should(Succeed())
			Ω(a.Database.Driver).Should(Equal(""))
			Ω(a.Database.DSN).Should(Equal(""))
			Ω(a.Web.Addr).Should(Equal(":8888"))
		})

		It("handles YAML files with all the directives", func() {
			Ω(a.ReadConfig("test/etc/valid.yml")).Should(Succeed())
			Ω(a.Database.Driver).Should(Equal("my-driver"))
			Ω(a.Database.DSN).Should(Equal("my:dsn=database"))
			Ω(a.Web.Addr).Should(Equal(":8988"))
		})

		It("autovivifies the api database object", func() {
			a.Database = nil
			Ω(a.ReadConfig("test/etc/valid.yml")).Should(Succeed())
			Ω(a.Database).ShouldNot(BeNil())
		})
	})
})
