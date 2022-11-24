package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vczyh/dbinsert/mysql"
	"github.com/vczyh/dbinsert/relational"
	"time"
)

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return StartMysql()
	},
}

var (
	definitionFile string
	databaseRepeat int
	mysqlCnf       = new(mysql.Config)
)

func init() {
	rootCmd.AddCommand(mysqlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mysqlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//mysqlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	mysqlCmd.Flags().StringVar(&mysqlCnf.Host, "host", "127.0.0.1", "mysql host")
	mysqlCmd.Flags().IntVar(&mysqlCnf.Port, "port", 3306, "mysql port")
	mysqlCmd.Flags().StringVar(&mysqlCnf.Username, "username", "", "mysql username")
	mysqlCmd.Flags().StringVar(&mysqlCnf.Password, "password", "", "mysql password")
	//mysqlCmd.Flags().StringVar(&mysqlCnf.Database, "database", "", "mysql database")
	mysqlCmd.Flags().BoolVar(&mysqlCnf.CreateDatabase, "create-database", false, "auto create database if not exist")
	mysqlCmd.Flags().BoolVar(&mysqlCnf.CreateTable, "create-table", false, "auto create table if not exist")
	//mysqlCmd.Flags().IntVarP(&mysqlCnf.TableSize, "table-size", "ts", 100, "table size")
	mysqlCmd.Flags().DurationVar(&mysqlCnf.Timeout, "timeout", 10*time.Hour, "timeout")

	mysqlCmd.Flags().StringVar(&definitionFile, "definition", "", "definition file path")
	mysqlCmd.Flags().IntVar(&databaseRepeat, "db-repeat", 0, "number of times the database is repeatedly created")
}

func StartMysql() error {
	schema, err := relational.ParseSchemaFromFile(relational.DialectMysql, definitionFile, databaseRepeat)
	if err != nil {
		return err
	}
	tables, err := relational.NewTables(schema)
	if err != nil {
		return err
	}
	mysqlCnf.Tables = tables

	manager, err := mysql.NewManager(mysqlCnf)
	if err != nil {
		return err
	}
	return manager.Start()
}
