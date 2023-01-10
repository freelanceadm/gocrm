package main

import (
	"database/sql"
	"fmt"
	"log"

	t_model "wcrm2/model"
	t_sql "wcrm2/pkg/client/postgresclient"
	"wcrm2/pkg/config"
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
	gdb.AutoMigrate(&t_model.User{})
	gdb.AutoMigrate(&t_model.Album{})

}

func main() {
	// check DB connection and schema
	//checkDB()

	user := t_model.User{
		Email:        "mail1@mail.ru",
		PasswordHash: "$1$dkjfngkjrnktjngkrjntgkjrntkj",
	}
	users := []t_model.User{
		{
			Email:        "mail2@mail.ru",
			PasswordHash: "$1$dkjfngkjrnktjngkrjntgkjrntkj",
		},
		{
			Email:        "mail3@mail.ru",
			PasswordHash: "$1$dkjfngkjrnktjngkrjntgkjrntkj",
		},
		user,
	}

	// insert multiple records
	t_sql.GormCreateBatch(gdb, &users)
	// insert one record
	t_sql.GormCreateOne(gdb, &user)
	// update record
	// delete

	// Start web server
	s := gmax.S
	s.StartServer()
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
