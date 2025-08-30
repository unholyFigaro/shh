/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/unholyFigaro/shh/internal/usecases/hosts"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm <name>",
	Args:  cobra.ExactArgs(1),
	Short: "Remove host",
	Long:  "Remove host by name or by host",
	RunE: func(cmd *cobra.Command, args []string) error {
		return hosts.RemoveHost(cmd.OutOrStdout(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
