package postgresclient

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Gorm
// Gorm use existing *sql.DB to create SQL connection.
func GormConnectDB(sqlDB *sql.DB) *gorm.DB {
	if sqlDB == nil {
		sqlDB = ConnectDB()
	}

	gormDB, err := gorm.Open(postgres.New(
		postgres.Config{
			Conn: sqlDB,
		}),
		&gorm.Config{})
	if err != nil {
		log.Fatal("Error during connecting to database using GORM.", err)
	}

	return gormDB
}

// CRUD GORM functions
// Create record
func GormCreate(gdb *gorm.DB, records ...interface{}) {
	// user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	// result := db.Create(&user) // pass pointer of data to Create

	// user.ID             // returns inserted data's primary key
	// result.Error        // returns error
	// result.RowsAffected // returns inserted records count
	for rec := range records {
		gdb.Create(&rec)
	}
}
