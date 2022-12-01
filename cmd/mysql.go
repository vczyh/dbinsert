package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/vczyh/dbinsert/relational"
	"time"
)

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Quick insert tool for mysql.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return StartMysql()
	},
}

var (
	mysqlCnf = new(relational.MysqlConfig)
)

func init() {
	rootCmd.AddCommand(mysqlCmd)

	mysqlCmd.Flags().StringVar(&mysqlCnf.Host, "host", "127.0.0.1", "mysql host")
	mysqlCmd.Flags().IntVar(&mysqlCnf.Port, "port", 3306, "mysql port")
	mysqlCmd.Flags().StringVar(&mysqlCnf.Username, "username", "", "mysql username")
	mysqlCmd.Flags().StringVar(&mysqlCnf.Password, "password", "", "mysql password")
	mysqlCmd.Flags().BoolVar(&mysqlCnf.CreateDatabase, "create-database", false, "auto create database if not exist")
	mysqlCmd.Flags().BoolVar(&mysqlCnf.CreateTable, "create-table", false, "auto create table if not exist")
	mysqlCmd.Flags().DurationVar(&mysqlCnf.Timeout, "timeout", 10*time.Hour, "timeout")
	mysqlCmd.Flags().IntVar(&mysqlCnf.TableSize, "table-size", 0, "table size")
	mysqlCmd.Flags().IntVar(&mysqlCnf.DatabaseRepeat, "db-repeat", 0, "number of times the database is repeatedly created")
}

func StartMysql() error {
	mysqlCnf.SchemaFile = definitionFile
	manager, err := relational.CreateManagerForMysql(mysqlCnf)
	if err != nil {
		return err
	}
	return manager.Start(context.Background())
}
