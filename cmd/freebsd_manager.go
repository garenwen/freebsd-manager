package cmd

import (
	"github.com/garenwen/freebsd-manager/cmd/iscsi"
	"github.com/garenwen/freebsd-manager/cmd/zfs"
	"github.com/spf13/cobra"
)

var FreebsdManagerCmd = &cobra.Command{
	Use: "fm",
}

func init() {
	cobra.EnableCommandSorting = false

	FreebsdManagerCmd.AddCommand(zfs.ZfsCmd)
	FreebsdManagerCmd.AddCommand(iscsi.IscsiCmd)
}
