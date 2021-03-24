package zfs

import (
	"github.com/garenwen/freebsd-manager/server/zfsserver"
	"github.com/spf13/cobra"
)

var ZfsCmd = &cobra.Command{
	Use:   "zfs",
	Short: "zfs manager",
	Long:  `zfs manager.`,
	// Args:  cobra.MinimumNArgs(1),
	Example: "go run main.go zfs",
	Run:     run,
}

var (
	tlsCertificate string
)

func init() {
	flags := ZfsCmd.Flags()
	flags.StringVar(&tlsCertificate, "tls-certificate", "", "the certificate to use for secure connections")

}

func run(*cobra.Command, []string) {
	zfsserver.NewZfsServer().Start()
}
