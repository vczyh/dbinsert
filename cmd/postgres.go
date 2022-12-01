package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/vczyh/dbinsert/relational"
	"time"
)

// mysqlCmd represents the mysql command
var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Quick insert tool for postgresql.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return StartPostgres()
	},
}

var (
	postgresCnf = new(relational.PostgresConfig)
)

func init() {
	rootCmd.AddCommand(postgresCmd)

	postgresCmd.Flags().StringVar(&postgresCnf.Host, "host", "127.0.0.1", "mysql host")
	postgresCmd.Flags().IntVar(&postgresCnf.Port, "port", 3306, "mysql port")
	postgresCmd.Flags().StringVar(&postgresCnf.Username, "username", "", "mysql username")
	postgresCmd.Flags().StringVar(&postgresCnf.Password, "password", "", "mysql password")
	postgresCmd.Flags().BoolVar(&postgresCnf.CreateDatabase, "create-database", false, "auto create database if not exist")
	postgresCmd.Flags().BoolVar(&postgresCnf.CreateTable, "create-table", false, "auto create table if not exist")
	postgresCmd.Flags().DurationVar(&postgresCnf.Timeout, "timeout", 10*time.Hour, "timeout")
	postgresCmd.Flags().IntVar(&postgresCnf.TableSize, "table-size", 0, "table size")
	postgresCmd.Flags().IntVar(&postgresCnf.DatabaseRepeat, "db-repeat", 0, "number of times the database is repeatedly created")
}

func StartPostgres() error {
	postgresCnf.SchemaFile = definitionFile
	manager, err := relational.CreateManagerForPostgres(postgresCnf)
	if err != nil {
		return err
	}
	return manager.Start(context.Background())
}
