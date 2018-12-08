package db

import(
    "log"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func init() {
    source := fmt.Sprintf("%s:%s@/%s", "root", "", "quicksilver")

	db, err := sql.Open("mysql", source)
	if err != nil {
		log.Fatal(err)
	}
}

func InitTables() error {
    _, err := db.Exec(`CREATE TABLE keys(
        id INT PRIMARY KEY AUTO_INCREMENT,
        owner VARCHAR(16),
        pub BLOB,
        priv BLOB
        );`)

    return err
}

func ResetTables() error {
    _, err := db.Exec("DROP TABLE IF EXISTS keys")
    if err != nil {
        return err
    }

    return InitTables()
}

func LookupPubKey(owner string) (pub []bytes, err error) {
    row := db.QueryRow(
        `SELECT priv
         FROM keys
         WHERE owner = ?
         ORDER BY id DESC`,
        string, pub)
    err = row.Scan(&pub)
    return
}

func LookupPrivKey(owner string, pub []bytes) (priv []bytes, err error) {
    row := db.QueryRow(
        `SELECT priv
         FROM keys
         WHERE owner = ? AND pub = ?`,
        string, pub)
    err = row.Scan(&priv)
    return
}

func UploadKey(owner int, pub, priv []bytes) error {
    _, err := db.Exec(`INSERT INTO keys (owner, pub, priv) VALUES (?, ?, ?)`, owner, pub, priv)
    return err
}
