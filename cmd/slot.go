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
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pglogrepl"
	"github.com/spf13/cobra"
)

// baseSlotCmd represents the slot command
var baseSlotCmd = &cobra.Command{
	Use:   "slot <slot-name> <postgres_url>",
	Short: "",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("slot called with args %s", args)
		slotName := args[0]
		url := args[1]

		switch cmd.Parent() {
		case createCmd:
			createSlot(slotName, url)
			// case dropCmd:
			// 	dropSlot(slotName, url)
		}
	},
}

func createSlot(slotName string, url string) {
	conn, err := pgconn.Connect(context.Background(), url)
	if err != nil {
		log.Fatalln("failed to connect to Postgres server:", err)
	}
	defer conn.Close(context.Background())

	// TODO: Allow temporary replication slot to be specified
	_, err = pglogrepl.CreateReplicationSlot(context.Background(), conn, slotName, outputPlugin, pglogrepl.CreateReplicationSlotOptions{Temporary: false})

	if err != nil {
		log.Fatalln(err)
	}
}
