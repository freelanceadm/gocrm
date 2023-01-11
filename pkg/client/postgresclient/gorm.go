package postgresclient

import (
	"database/sql"
	"errors"
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
// Create one record
func GormCreateOne(gdb *gorm.DB, r interface{}) error {
	// user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	// result := db.Create(&user) // pass pointer of data to Create

	// user.ID             // returns inserted data's primary key
	// result.Error        // returns error
	// result.RowsAffected // returns inserted records count

	log.Println(r)
	result := gdb.Create(r)
	if result.Error != nil {
		log.Printf("Error: %v", result.Error)
		return result.Error
	}
	log.Printf("gorm: data: %v rows affected %v", r, result.RowsAffected)
	return nil
}

// Create multiple/batch records
func GormCreateBatch(gdb *gorm.DB, r interface{}) error {
	log.Printf("gormbatch: %v", r)

	result := gdb.CreateInBatches(r, 100)
	if result.Error != nil {
		log.Printf("Error: %v", result.Error)
		return result.Error
	}
	log.Printf("gorm: data: %v rows affected %v", r, result.RowsAffected)
	return nil
}

// Get one record
func GormGetByID(gdb *gorm.DB, r interface{}, id string) (interface{}, error) {
	var err error

	result := map[string]interface{}{}

	err = gdb.Debug().Model(r).Where("id = ?", id).Take(&result).Error
	if err != nil {
		return r, err
	}
	if gorm.ErrRecordNotFound == err {
		return r, errors.New("Record Not Found")
	}
	return result, nil
}

// Get multiple records
func GormGetAll(gdb *gorm.DB, model interface{}, result interface{}) error {
	var err error

	//gdb.Debug().Model(r).Limit(100).Find(&results) //.FindInBatches(&rows, 100, ) //	Find(&rows)
	err = gdb.Debug().Model(model).Find(result).Error
	if err != nil {
		return err
	}
	return nil
}
