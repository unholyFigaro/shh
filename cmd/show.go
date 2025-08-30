/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unholyFigaro/shh/internal/usecases/hosts"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show host details by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) >= 1 {
			names := make([]string, 0, len(args))
			for _, arg := range args {
				names = append(names, arg)
			}
			hosts.ShowHostsByName(cmd.OutOrStdout(), names)
			return nil
		}

		return fmt.Errorf("expected at least one host name")
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
