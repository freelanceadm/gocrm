package postgresclient

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ConnectDB() *sql.DB {
	// Capture connection properties.

	cfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("postgresql.host"),
		viper.GetString("postgresql.port"),
		viper.GetString("postgresql.user"),
		viper.GetString("postgresql.password"),
		viper.GetString("postgresql.db"))
	log.Println("Trying to connect to DB using config:", cfg)
	// Get a database handle.
	DB, err := sql.Open("postgres", cfg)
	CheckError(err)
	// defer DB.Close()

	err = DB.Ping()
	CheckError(err)

	fmt.Println("Connected!")
	return DB
}

func CheckError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// This function connect to database.
// Check if database exist. In case not it ll create required DB.
// This means that the recommendations for sql.Open apply to gorm.Open:
// The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
// Thus, the Open function should be called just once. It is rarely necessary to close a DB.
func CreateDB() error {
	// GORM configuration
	gormConfig := gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy:         nil,
		FullSaveAssociations:   false,
		//Logger:                 logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now()
		},
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false, // TODO: read doc about nested transactions
		AllowGlobalUpdate:                        false, // TODO: read what is global update
		QueryFields:                              false,
		CreateBatchSize:                          0,
		ClauseBuilders:                           map[string]clause.ClauseBuilder{},
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  map[string]gorm.Plugin{},
	}

	// make connection string
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		viper.GetString("postgresql.user"),
		viper.GetString("postgresql.password"),
		viper.GetString("postgresql.host"),
		viper.GetString("postgresql.port"),
		"postgres", // viper.GetString("postgresql.db")
		viper.GetString("postgresql.sslmode"),
	)

	// connect to the postgres db just to be able to run the create db statement
	Db, err := gorm.Open(postgres.Open(connStr), &gormConfig)
	if err != nil {
		log.Fatal("Error during connecting to database.. ", err)
	}

	// check if db exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", viper.GetString("postgresql.db"))
	rs := Db.Raw(stmt)
	if rs.Error != nil {
		log.Println("Initial database does not exist. I ll create it... Error: ", rs.Error)
	}

	// if not create it
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		stmt := fmt.Sprintf("CREATE DATABASE %s;", viper.GetString("postgresql.db"))
		if rs := Db.Exec(stmt); rs.Error != nil {
			log.Fatal("Could not create database. Error: ", rs.Error)
		}

		// close db connection
		// sql, err := Db.DB()
		// defer func() { // TODO: read how to close db connections. Use db connection pools.
		// 	_ = sql.Close()
		// }()
		// if err != nil {
		// 	log.Fatalf("Could not connect to %s database. Error:%s", viper.GetString("postgresql.db"), err)
		// }
	}

	// switch db connection to required DB
	// first close current connection
	sqlDB, err := Db.DB()
	if err != nil {
		log.Fatal("Error: Could not get sqlDB. ", err)
	}
	// Close
	err = sqlDB.Close()
	if err != nil {
		log.Fatal("Error: Could not close connection to DB. ", err)
	}

	return nil
}
