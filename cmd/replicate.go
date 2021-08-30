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
	"fmt"

	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgproto3/v2"
	"github.com/spf13/cobra"
)

const newline = '\n'

// replicateCmd represents the replicate command
var replicateCmd = &cobra.Command{
	Use:   "replicate <postgres_url>",
	Short: "begins replication from postgres",
	Long:  `begins replication from postgres`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		replicate(url)
	},
}

func replicate(url string) {
	log.Printf("beginning replication on url: %s, slotName: %s, publicationName: %s", url, slotName, publicationName)
	conn, err := pgconn.Connect(context.Background(), url)
	if err != nil {
		log.Fatalln("failed to connect to PostgresQL server:", err)
	}
	defer conn.Close(context.Background())

	pluginArguments := []string{"proto_version '1'", fmt.Sprintf("publication_names '%s'", publicationName)}
	var clientXLogPos pglogrepl.LSN = 0
	err = pglogrepl.StartReplication(context.Background(), conn, slotName, clientXLogPos, pglogrepl.StartReplicationOptions{PluginArgs: pluginArguments})
	if err != nil {
		log.Fatalln("StartReplication failed:", err)
	}
	log.Println("Logical replication started on slot", slotName)
	standbyMessageTimeout := time.Second * 10
	nextStandbyMessageDeadline := time.Now().Add(standbyMessageTimeout)
	writer := bufio.NewWriter(os.Stdout)

	for {
		if time.Now().After(nextStandbyMessageDeadline) {
			err = pglogrepl.SendStandbyStatusUpdate(context.Background(), conn, pglogrepl.StandbyStatusUpdate{WALWritePosition: clientXLogPos})
			if err != nil {
				log.Fatalln("SendStandbyStatusUpdate failed:", err)
			}
			nextStandbyMessageDeadline = time.Now().Add(standbyMessageTimeout)
		}
		ctx, cancel := context.WithDeadline(context.Background(), nextStandbyMessageDeadline)
		msg, err := conn.ReceiveMessage(ctx)
		cancel()
		if err != nil {
			if pgconn.Timeout(err) {
				continue
			}
			log.Fatalln("ReceiveMessage failed:", err)
		}
		switch msg := msg.(type) {
		case *pgproto3.CopyData:
			switch msg.Data[0] {
			case pglogrepl.PrimaryKeepaliveMessageByteID:
				pkm, err := pglogrepl.ParsePrimaryKeepaliveMessage(msg.Data[1:])
				if err != nil {
					log.Fatalln("ParsePrimaryKeepaliveMessage failed:", err)
				}
				if pkm.ReplyRequested {
					nextStandbyMessageDeadline = time.Time{}
				}
			case pglogrepl.XLogDataByteID:
				xld, err := pglogrepl.ParseXLogData(msg.Data[1:])
				if err != nil {
					log.Fatalln("ParseXLogData failed:", err)
				}
				logicalMsg, err := pglogrepl.Parse(xld.WALData)
				if err != nil {
					log.Fatalf("Parse logical replication message: %s", err)
				}
				jsonBytes, err := json.Marshal(logicalMsg)
				if err != nil {
					log.Fatalf("json.Marshal(logicalMsg) failed: ", err)
				}
				count, err := writer.WriteString(string(append(jsonBytes, newline)))
				if err != nil {
					log.Fatalf("Error writing json marshaled logical message. Wrote %d bytes to stdout buffer: %v", count, err)
				}
				writer.Flush()
				clientXLogPos = xld.WALStart + pglogrepl.LSN(len(xld.WALData))
			}
		default:
			log.Fatalf("Received unexpected message: %#v\n", msg)
		}
	}
}
