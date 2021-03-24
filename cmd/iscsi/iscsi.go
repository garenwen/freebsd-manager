package iscsi

import (
	"github.com/garenwen/freebsd-manager/server/iscsiserver"
	"github.com/spf13/cobra"
)

var IscsiCmd = &cobra.Command{
	Use:   "iscsi",
	Short: "iscsi manager",
	Long:  `iscsi manager.`,
	// Args:  cobra.MinimumNArgs(1),
	Example: "go run main.go iscsi",
	Run:     run,
}

var (
	tlsCertificate string
)

func init() {
	flags := IscsiCmd.Flags()
	flags.StringVar(&tlsCertificate, "tls-certificate", "", "the certificate to use for secure connections")

}

func run(*cobra.Command, []string) {
	iscsiserver.NewIscsiServer().Start()
}
