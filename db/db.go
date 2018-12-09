package db

import(
    "fmt"
    "log"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func init() {
    source := fmt.Sprintf("%s:%s@/%s", "root", "", "quicksilver")

    var err error
	db, err = sql.Open("mysql", source)
	if err != nil {
		log.Fatal(err)
	}
}

func InitTables() error {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS dh(
        id INT PRIMARY KEY AUTO_INCREMENT,
        owner VARCHAR(16),
        pub BLOB,
        priv BLOB
        );`)

    return err
}

func ResetTables() error {
    _, err := db.Exec("DROP TABLE IF EXISTS dh")
    if err != nil {
        return err
    }

    return InitTables()
}

func LookupPubKey(owner string) (pub []byte, err error) {
    row := db.QueryRow(
        `SELECT pub
         FROM dh
         WHERE owner = ?
         ORDER BY id DESC`,
        owner)
    err = row.Scan(&pub)
    return
}

func LookupPrivKey(owner string, pub []byte) (priv []byte, err error) {
    row := db.QueryRow(
        `SELECT priv
         FROM dh
         WHERE owner = ? AND pub = ? AND priv IS NOT NULL`,
        owner, pub)
    err = row.Scan(&priv)
    return
}

func UploadKey(owner string, pub, priv []byte) error {
    if priv != nil {
        _, err := db.Exec(`INSERT INTO dh (owner, pub, priv) VALUES (?, ?, ?)`, owner, pub, priv)
        return err
    }
    _, err := db.Exec(`INSERT INTO dh (owner, pub) VALUES (?, ?)`, owner, pub)
    return err
}
