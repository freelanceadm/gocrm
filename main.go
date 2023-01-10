package main

import (
	"database/sql"
	"fmt"
	"log"

	"wcrm2/config"
	t_model "wcrm2/model"
	t_sql "wcrm2/pkg/client/postgresclient"
	"wcrm2/router/gmax"

	"gorm.io/gorm"
	// t_sql "wcrm2/pkg/client/mysql"
)

var (
	db  *sql.DB
	gdb *gorm.DB
	s   gmax.APIServer
)

func init() {
	// read configuration
	if err := config.ReadConfig(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Get a database handle.
	db = t_sql.ConnectDB()
	gdb = t_sql.GormConnectDB(db)

	// Migrate schema(s)
	gdb.AutoMigrate(&t_model.User{}, &t_model.Album{})

	// router and handlers
	s := gmax.S
	s.StartServer()
}

func main() {
	// check DB connection and schema
	checkDB()

	// start http server

}

// Check db and tables exist with required data
func checkDB() {
	albums, err := t_model.AlbumsByArtist("John Coltrane", db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 here to test the query.
	alb, err := t_model.AlbumByID(2, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := t_model.AddAlbum(t_model.Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}
