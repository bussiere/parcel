package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Parcel", func() {
	var (
		cmd      *exec.Cmd
		dir      string
		args     []string
		resource string
	)

	BeforeEach(func() {
		args = []string{}
	})

	JustBeforeEach(func() {
		var err error

		dir, err = ioutil.TempDir("", "gom")
		Expect(err).To(BeNil())

		cmd = exec.Command(embedoPath, args...)
		cmd.Dir = dir

		path := filepath.Join(cmd.Dir, "/database")
		Expect(os.MkdirAll(path, 0700)).To(Succeed())

		path = filepath.Join(path, "main.sql")
		Expect(ioutil.WriteFile(path, []byte("main"), 0700)).To(Succeed())

		path = filepath.Join(cmd.Dir, "/database/command")
		Expect(os.MkdirAll(path, 0700)).To(Succeed())

		path = filepath.Join(path, "commands.sql")
		Expect(ioutil.WriteFile(path, []byte("command"), 0700)).To(Succeed())

		resource = filepath.Join(cmd.Dir, "/resource.go")
	})

	It("generates resource on root level", func() {
		cmd.Args = append(cmd.Args, "-r")

		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		Expect(session.Out).To(gbytes.Say("Compressing 'database/main.sql'"))
		Expect(session.Out).NotTo(gbytes.Say("Compressing 'database/command/commands.sql'"))
		Expect(resource).To(BeARegularFile())

		data, err := ioutil.ReadFile(resource)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(ContainSubstring("// Auto-generated"))
		Expect(string(data)).To(ContainSubstring(fmt.Sprintf("package %s", filepath.Base(cmd.Dir))))
		Expect(string(data)).To(ContainSubstring("parcel.AddResource"))
	})

	Context("when the commands.sql is ignored", func() {
		BeforeEach(func() {
			args = append(args, "-r", "-i", "commands.sql")
		})

		It("does not generate embedded resource for it", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out).To(gbytes.Say("Compressing 'database/main.sql'"))
			Expect(session.Out).NotTo(gbytes.Say("Compressing 'database/command/commands.sql'"))
			Expect(resource).To(BeARegularFile())
		})
	})

	Context("when the documentation is disabled", func() {
		BeforeEach(func() {
			args = append(args, "-r", "-include-docs=false")
		})

		It("does not include API documentation", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out).To(gbytes.Say("Compressing 'database/main.sql'"))
			Expect(resource).To(BeARegularFile())

			data, err := ioutil.ReadFile(resource)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(ContainSubstring("parcel.AddResource"))
			Expect(string(data)).NotTo(ContainSubstring("// Auto-generated"))
		})
	})

	Context("when quite model is enabled", func() {
		BeforeEach(func() {
			args = append(args, "-r", "-q")
		})

		It("does not print anything on stdout", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out).NotTo(gbytes.Say("Compressing 'database/main.sql'"))
			Expect(session.Out).NotTo(gbytes.Say("Compressing 'database/command/commands.sql'"))
			Expect(resource).To(BeARegularFile())
		})
	})

	Context("when the recursion is disabled", func() {
		It("generates resource for all directories", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out).NotTo(gbytes.Say("Compressing"))
			Expect(resource).NotTo(BeARegularFile())
		})
	})

	Context("when the directory is provided", func() {
		BeforeEach(func() {
			args = []string{"-r", "-d", "./database"}
		})

		It("compresses the directory successfully", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out).To(gbytes.Say("Compressing 'main.sql'"))
			Expect(resource).To(BeARegularFile())
		})
	})

	Context("when the bundle-dir is provided", func() {
		BeforeEach(func() {
			args = []string{"-r", "-b", "./database"}
		})

		It("returns an error", func() {
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			resource = filepath.Join(cmd.Dir, "database", "resource.go")
			Expect(resource).To(BeARegularFile())

			data, err := ioutil.ReadFile(resource)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(ContainSubstring("package database"))
		})
	})
})
