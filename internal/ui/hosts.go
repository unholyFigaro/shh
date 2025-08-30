package ui

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/unholyFigaro/shh/internal/domain"
)

func PrintHosts(w io.Writer, hosts map[string]domain.Host) error {
	sortedNames := make([]string, 0, len(hosts))
	for name := range hosts {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "NAME\tHOST\tPORT\tUSER")
	for _, name := range sortedNames {
		host := hosts[name]
		port := host.Port
		if port == 0 {
			port = 22
		}
		user := host.User
		if user == "" {
			user = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n", name, host.Host, port, user)
	}
	return nil
}
