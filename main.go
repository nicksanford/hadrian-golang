package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgproto3/v2"
)

const newline = '\n'

/*
Commands
hadrian replicate --slot <hadrian> --temporary <url>
hadrian create publication <hadrian> <url>
hadrian create replication_slot <hadrian> <url>
hadrian drop publication <hadrian> <url>
hadrian drop replication_slot <hadrian> <url>

when hadrian boots, it ensures:
1.that there is a publication for all tables
	a. if there is no plublication with that name, then it creates one
	b. if there is a publication with that name, but not for all tables, terminate with an error
	c. if there is a publication with that name for all tables, no op
2. that there is a replication slot of the type specified
	a. if there isn't a slot with that name it creates one
	b. if there is a slot with the name, but it is of a different type, terminate with an error
	c. if there is a slot with the name & desired type, no op
3.
*/

func main() {
	const outputPlugin = "pgoutput"
	slotName := "hadrian"
	conn, err := pgconn.Connect(context.Background(), os.Getenv("PGLOGREPL_DEMO_CONN_STRING"))
	if err != nil {
		log.Fatalln("failed to connect to PostgresQL server:", err)
	}
	defer conn.Close(context.Background())

	// result := conn.Exec(context.Background(), "DROP PUBLICATION IF EXISTS hadrian;")
	// _, err := result.ReadAll()
	// if err != nil {
	// 	log.Fatalln("drop publication if exists error", err)
	// }
	// result = conn.Exec(context.Background(), "CREATE PUBLICATION hadrian FOR ALL TABLES;")
	// _, err = result.ReadAll()
	// if err != nil {
	// 	log.Fatalln("create publication error", err)
	// }
	// log.Println("create publication hadrian")

	pluginArguments := []string{"proto_version '1'", "publication_names 'hadrian'"}
	// sysident, err := pglogrepl.IdentifySystem(context.Background(), conn)
	// if err != nil {
	// 	log.Fatalln("IdentifySystem failed:", err)
	// }
	// log.Println("SystemID:", sysident.SystemID, "Timeline:", sysident.Timeline, "XLogPos:", sysident.XLogPos, "DBName:", sysident.DBName)

	// TODO: set Temporary: false
	// _, err = pglogrepl.CreateReplicationSlot(context.Background(), conn, slotName, outputPlugin, pglogrepl.CreateReplicationSlotOptions{Temporary: false})
	// if err != nil {
	// 	log.Fatalln("CreateReplicationSlot failed:", err)
	// }
	// log.Println("Created temporary replication slot:", slotName)
	// err = pglogrepl.StartReplication(context.Background(), conn, slotName, sysident.XLogPos, pglogrepl.StartReplicationOptions{PluginArgs: pluginArguments})
	// setting offset position to be where postgres thinks we left off
	var clientXLogPos pglogrepl.LSN = 0
	err = pglogrepl.StartReplication(context.Background(), conn, slotName, clientXLogPos, pglogrepl.StartReplicationOptions{PluginArgs: pluginArguments})
	if err != nil {
		log.Fatalln("StartReplication failed:", err)
	}
	log.Println("Logical replication started on slot", slotName)
	standbyMessageTimeout := time.Second * 10
	nextStandbyMessageDeadline := time.Now().Add(standbyMessageTimeout)
	writer := bufio.NewWriter(os.Stdout)
	go func() {
		log.Println("waiting for CleanupDone")
		<-conn.CleanupDone()
		log.Fatalln("connection died")
	}()

	for {
		log.Println("time.Now().After(nextStandbyMessageDeadline) %s", time.Now().After(nextStandbyMessageDeadline))
		log.Printf("now: %s, nextStandbyMessageDeadline: %s", time.Now(), nextStandbyMessageDeadline)
		log.Printf("conn: %v", conn)
		if time.Now().After(nextStandbyMessageDeadline) {
			err = pglogrepl.SendStandbyStatusUpdate(context.Background(), conn, pglogrepl.StandbyStatusUpdate{WALWritePosition: clientXLogPos})
			if err != nil {
				log.Fatalln("SendStandbyStatusUpdate failed:", err)
			}
			log.Println("Sent Standby status message")
			nextStandbyMessageDeadline = time.Now().Add(standbyMessageTimeout)
		}
		ctx, cancel := context.WithDeadline(context.Background(), nextStandbyMessageDeadline)
		log.Printf("calling receive message, ctx: %v, cancel: %v, nextStandbyMessageDeadline: %v", ctx, cancel, nextStandbyMessageDeadline)
		msg, err := conn.ReceiveMessage(ctx)
		log.Printf("err: %v, msg: %v, IsBusy: %v", "IsClosed: %v", err, msg, conn.IsBusy(), conn.IsClosed())
		cancel()
		if err != nil {
			log.Printf(" iv", err)
			if pgconn.Timeout(err) {
				log.Println("timeout")
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
				log.Println("Primary Keepalive Message =>", "ServerWALEnd:", pkm.ServerWALEnd, "ServerTime:", pkm.ServerTime, "ReplyRequested:", pkm.ReplyRequested)
				if pkm.ReplyRequested {
					nextStandbyMessageDeadline = time.Time{}
				}
			case pglogrepl.XLogDataByteID:
				xld, err := pglogrepl.ParseXLogData(msg.Data[1:])
				if err != nil {
					log.Fatalln("ParseXLogData failed:", err)
				}
				log.Println("XLogData =>", "WALStart", xld.WALStart, "ServerWALEnd", xld.ServerWALEnd, "ServerTime:", xld.ServerTime, "WALData", string(xld.WALData))
				logicalMsg, err := pglogrepl.Parse(xld.WALData)
				if err != nil {
					log.Fatalf("Parse logical replication message: %s", err)
				}
				log.Printf("Receive a logical replication message: %s", logicalMsg.Type())
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
			log.Printf("Received unexpected message: %#v\n", msg)
		}
	}
}
