package cmd

import (
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/spf13/cobra"
)

// HashCmd - `tok up`
var HashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Returns a valid Drupal Hash Salt",
	Long:  "Hash simply generates and returns Drupal hash salt. It is not saved anywhere",
	Run: func(cmd *cobra.Command, args []string) {
		drupal.Hash()
	},
}
