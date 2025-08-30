/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unholyFigaro/shh/internal/usecases/hosts"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <name> [flags]",
	Short: "Add a new host",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			return err
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}
		user, err := cmd.Flags().GetString("user")
		if err != nil {
			return err
		}
		name := args[0]

		params := map[string]any{
			"name":  name,
			"host":  host,
			"port":  port,
			"force": force,
			"user":  user,
		}
		err = hosts.AddHost(cmd.Context(), params)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Host %q added on port %d \n", host, port)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().StringP("host", "H", "", "hostname or IP address (required)")
	addCmd.Flags().StringP("user", "u", "", "user for remote host")
	addCmd.MarkFlagRequired("host")
	addCmd.Flags().IntP("port", "p", 22, "port number (default 22)")
	addCmd.Flags().BoolP("force", "f", false, "overwrite existing host with the same name")
}
