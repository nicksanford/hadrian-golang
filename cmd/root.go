/*
Copyright © 2021 Nick Sanford <nicholascsanford@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var cfgFile string
var slotName string
var publicationName string

const outputPlugin = "pgoutput"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hadrian",
	Short: "WIP: A single binary debezium for postgres",
	Long: `
	THIS IS A WORK IN PROGRESS & NOT YET PRODUCTION READY

	Hadrian is a CLI that replicates postgres database changes for application
	consumption.

	It does this by connecting to a target posgres server as a logical streaming
	replication client, and outputs all database changes in a format applications
	can understand, such as json.

	Hadrian is intended to solve the dual write problem, described in Kleppmann's
	Bottled Water blog post, which frequently occur when multiple systems need to
	be notified of data changes in a given database.

	Debezium is currently the industry standard solution for this problem, however
	it requires both JVM and/or kafka expertise to use effectively.

	Hadrian is intended to provide a subset of debezium's features, so that small
	teams or single person projects without high availability or high throughput
	requirements can replicate data from their postgres databases to downstream
	systems without needing to change application code.

	Credits:
		Hadrian is a wrapper around @jackc's fantastic pglogrepl and was inspired by
		both cainophile: and supabase's realtime.

	References:
		Postgres Streaming Replication:
		https://www.postgresql.org/docs/current/protocol-replication.html.

		Debezium:
		https://debezium.io/

		Bottled Water Blog Post:
		https://www.confluent.io/blog/bottled-water-real-time-integration-of-postgresql-and-kafka/.

		@jackc's pglogrepl:
		https://github.com/jackc/pglogrepl

		Cainophile:
		https://github.com/cainophile/cainophile

		Supabase Realtime:
		https://github.com/supabase/realtime.


		`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	var dropPublicationCmd = *basePublicationCmd
	var dropSlotCmd = *baseSlotCmd

	dropCmd.AddCommand(&dropPublicationCmd)
	dropCmd.AddCommand(&dropSlotCmd)
	rootCmd.AddCommand(dropCmd)

	var createPublicationCmd = *basePublicationCmd
	var createSlotCmd = *baseSlotCmd

	createCmd.AddCommand(&createPublicationCmd)
	createCmd.AddCommand(&createSlotCmd)
	rootCmd.AddCommand(createCmd)

	replicateCmd.PersistentFlags().StringVarP(&slotName, "slot", "s", "", "replication slot (required)")
	replicateCmd.PersistentFlags().StringVarP(&publicationName, "publication", "p", "", "publication (required)")
	replicateCmd.MarkPersistentFlagRequired("slot")
	replicateCmd.MarkPersistentFlagRequired("publication")
	rootCmd.AddCommand(replicateCmd)
}
