/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/unholyFigaro/shh/internal/completion"
	"github.com/unholyFigaro/shh/internal/usecases/hosts"
)

var flagBastion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "shh <name>",
	Short:             "SSH management tool",
	SilenceUsage:      true,
	Args:              cobra.ArbitraryArgs,
	ValidArgsFunction: completeHostArg,
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.Help()
			return nil
		}
		name := args[0]
		return hosts.ConnectToHostByName(cmd.Context(), name, flagBastion)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDescriptions = true
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.RegisterFlagCompletionFunc("jump", completeBastionFlag)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shh.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&flagBastion, "jump", "j", "", "Jump host in <host> format")
}

func completeHostArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) >= 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names, _ := completion.HostNamesByPrefix(toComplete)
	return names, cobra.ShellCompDirectiveNoFileComp
}

func completeBastionFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	names, _ := completion.HostNamesByPrefix(toComplete)
	return names, cobra.ShellCompDirectiveNoFileComp
}
