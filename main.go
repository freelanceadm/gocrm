package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	gdb.AutoMigrate(&t_model.Album{})
	gdb.Debug().AutoMigrate(&t_model.User{})

}

func main() {
	// check DB connection and schema
	//checkDB()

	user := t_model.User{
		ID:        0,
		Nickname:  "U1",
		Email:     "mail1@mail.ru",
		Password:  "$1$dkjfngkjrnktjngkrjntgkjrntkj",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	users := []t_model.User{
		{
			ID:        0,
			Nickname:  "U2",
			Email:     "mail2@mail.ru",
			Password:  "$1$dkjfngkjrnktjngkrjntgkjrntkj",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		{
			ID:        0,
			Nickname:  "U3",
			Email:     "mail3@mail.ru",
			Password:  "$1$dkjfngkjrnktjngkrjntgkjrntkj",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		user,
	}

	// insert multiple records
	t_sql.GormCreateBatch(gdb, &users)
	// insert one record
	t_sql.GormCreateOne(gdb, &user)
	// get one record
	record, _ := t_sql.GormGetByID(gdb, &t_model.User{}, "1")
	log.Println(record)

	// Get all records
	res := []t_model.User{}
	err := t_sql.GormGetAll(gdb, t_model.User{}, &res)
	log.Println(err)
	log.Println("allrecords: ", res)

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
