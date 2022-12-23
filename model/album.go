package model

import (
	"database/sql"
	"fmt"
	"log"
)

// Link https://go.dev/doc/tutorial/database-access
// DROP TABLE IF EXISTS album;
// CREATE TABLE album (
//   id         INT AUTO_INCREMENT NOT NULL,
//   title      VARCHAR(128) NOT NULL,
//   artist     VARCHAR(255) NOT NULL,
//   price      DECIMAL(5,2) NOT NULL,
//   PRIMARY KEY (`id`)
// );

// INSERT INTO album
//
//	(title, artist, price)
//
// VALUES
//
//	('Blue Train', 'John Coltrane', 56.99),
//	('Giant Steps', 'John Coltrane', 63.99),
//	('Jeru', 'Gerry Mulligan', 17.99),
//	('Sarah Vaughan', 'Sarah Vaughan', 34.98);
//
// album represents data about a record album.
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albumsByArtist queries for albums that have the specified artist name.
// !!! mysql use ? as parameter placeholder
// !!! postgresql use $1 $2 ...
func AlbumsByArtist(name string, db *sql.DB) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func AlbumByID(id int64, db *sql.DB) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func AddAlbum(alb Album, db *sql.DB) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", alb.Title, alb.Artist, alb.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
