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
	"github.com/spf13/cobra"
)

// basePublicationCmd represents the publication command
var basePublicationCmd = &cobra.Command{
	Use:   "publication <slot-name> <postgres_url>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		publicationName := args[0]
		url := args[1]

		switch cmd.Parent() {
		case createCmd:
			createPublication(publicationName, url)
		case dropCmd:
			dropPublication(publicationName, url)
		}
	},
}

func createPublication(publicationName string, url string) {
	conn, err := pgconn.Connect(context.Background(), url)
	if err != nil {
		log.Fatalln("failed to connect to PostgresQL server:", err)
	}
	defer conn.Close(context.Background())

	command := fmt.Sprintf("CREATE PUBLICATION %s FOR ALL TABLES;", publicationName)
	log.Printf("executing '%s'", command)
	result := conn.Exec(context.Background(), command)
	_, err = result.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
}

func dropPublication(publicationName string, url string) {
	conn, err := pgconn.Connect(context.Background(), url)
	if err != nil {
		log.Fatalln("failed to connect to PostgresQL server:", err)
	}
	defer conn.Close(context.Background())

	command := fmt.Sprintf("DROP PUBLICATION IF EXISTS %s;", publicationName)
	log.Printf("executing '%s'", command)
	result := conn.Exec(context.Background(), command)
	_, err = result.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
}
