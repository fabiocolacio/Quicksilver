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
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS pubkeys(
        id INT,
        owner VARCHAR(16),
        peer VARCHAR(16),
        pubkey BLOB,
        primary key (id, owner, peer));`)
    if err != nil {
        return err
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS privkeys(
        id INT,
        owner VARCHAR(16),
        peer VARCHAR(16),
        privkey BLOB,
        primary key (id, owner, peer));`)

    return err
}

func ResetTables() error {
    _, err := db.Exec("DROP TABLE IF EXISTS pubkeys")
    if err != nil {
        return err
    }

    _, err = db.Exec("DROP TABLE IF EXISTS privkeys")
    if err != nil {
        return err
    }

    return InitTables()
}

func LookupPubKey(owner, peer string, id int) (pub []byte, err error) {
    row := db.QueryRow(
        `SELECT pubkey
         FROM pubkeys
         WHERE owner = ? AND peer = ? AND id = ?`,
        owner, peer, id)
    err = row.Scan(&pub)
    return
}

func LatestPubKey(owner, peer string) (id int) {
    row := db.QueryRow(
        `SELECT id
         FROM pubkeys
         WHERE owner = ? AND peer = ?
         ORDER BY id DESC`,
        owner, peer)
    err := row.Scan(&id)
    if err != nil {
        return 0
    }
    return id
}


func LookupPrivKey(owner, peer string, id int) (priv []byte, err error) {
    row := db.QueryRow(
        `SELECT privkey
         FROM privkeys
         WHERE owner = ? AND peer = ? AND id = ?
         ORDER BY id DESC`,
        owner, peer, id)
    err = row.Scan(&priv)
    return
}

func LatestPrivKey(owner, peer string) (id int) {
    row := db.QueryRow(
        `SELECT id
         FROM privkeys
         WHERE owner = ? AND peer = ?
         ORDER BY id DESC`,
        owner, peer)
    err := row.Scan(&id)
    if err != nil {
        return 0
    }
    return id
}

func UploadPubKey(owner, peer string, pubkey []byte, id int) error {
    _, err := db.Exec(`INSERT INTO pubkeys (id, owner, peer, pubkey) VALUES (?, ?, ?, ?)`, id, owner, peer, pubkey)
    return err
}

func UploadPrivKey(owner, peer string, privkey []byte, id int) error {
    _, err := db.Exec(`INSERT INTO privkeys (id, owner, peer, privkey) VALUES (?, ?, ?, ?)`, id, owner, peer, privkey)
    return err
}
