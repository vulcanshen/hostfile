package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
)

var addCmd = &cobra.Command{
	Use:   "add <ip> <domain1> [domain2...]",
	Short: "Add domains to an IP address",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		domains := args[1:]

		before, block, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		conflicts, err := manager.Add(block, ip, domains)
		if err != nil {
			exitWithError(err)
		}

		if len(conflicts) > 0 {
			for _, c := range conflicts {
				fmt.Printf("domain %q is already mapped to %s\n", c.Domain, c.CurrentIP)
			}
			if confirm("Move these domains to " + ip + "?") {
				if err := manager.AddForce(block, ip, domains); err != nil {
					exitWithError(err)
				}
			} else {
				fmt.Println("aborted")
				return
			}
		}

		if err := writeBlock(before, block, after); err != nil {
			exitWithError(err)
		}
		for _, d := range domains {
			fmt.Printf("added %s -> %s\n", d, ip)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
