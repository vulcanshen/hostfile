package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulcanshen/hostfile/manager"
	"github.com/vulcanshen/hostfile/parser"
)

var addCmd = &cobra.Command{
	Use:   "add <ip> <domain1> [domain2...]",
	Short: "Add domains to an IP address",
	Args: cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return completeExistingIPs()
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		domains := args[1:]

		before, block, after, err := readBlock()
		if err != nil {
			exitWithError(err)
		}

		// check which domains already exist on this IP
		var newDomains, existingDomains []string
		for _, d := range domains {
			found := false
			for _, entry := range block.Entries {
				if entry.IP == ip && entry.DisableType == parser.DisableNone && manager.ContainsDomain(entry.Domains, d) {
					found = true
					break
				}
			}
			if found {
				existingDomains = append(existingDomains, d)
			} else {
				newDomains = append(newDomains, d)
			}
		}

		if len(newDomains) == 0 {
			for _, d := range existingDomains {
				fmt.Printf("%s -> %s already exists\n", d, ip)
			}
			return
		}

		conflicts, err := manager.Add(block, ip, newDomains)
		if err != nil {
			exitWithError(err)
		}

		if len(conflicts) > 0 {
			for _, c := range conflicts {
				fmt.Printf("domain %q is already mapped to %s\n", c.Domain, c.CurrentIP)
			}
			if confirm("Move these domains to " + ip + "?") {
				if err := manager.AddForce(block, ip, newDomains); err != nil {
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
		for _, d := range newDomains {
			fmt.Printf("added %s -> %s\n", d, ip)
		}
		for _, d := range existingDomains {
			fmt.Printf("%s -> %s already exists\n", d, ip)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
