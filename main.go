package main

import(
    // "fmt"
    "github.com/fabiocolacio/quicksilver/crypto"
    "crypto/elliptic"
    "crypto/rand"
    "math/big"
    "database/sql"
    "os"
    "io/ioutil"
    "encoding/json"
    "log"
    "time"
    "sync"
    "github.com/fabiocolacio/quicksilver/gui"
    "github.com/fabiocolacio/quicksilver/db"
    "github.com/fabiocolacio/quicksilver/api"
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
)

func main() {
    err := db.InitTables()
    if err != nil {
        log.Fatal(err)
    }

    gtk.Init(nil)

    ui, err := gui.UINew()
    if err != nil {
        log.Fatal(err)
    }

    user, passwd, err := gui.LoginDialogRun(ui.Window)
    if err != nil {
        log.Fatal(err)
    }

    jwt, err := api.Login(user, passwd)
    if err != nil {
        log.Fatal(err)
    }

    peer, err := gui.PeerDialogRun(ui.Window)
    if err != nil {
        log.Fatal(err)
    }

    if err = api.LookupUser(peer); err != nil {
        log.Fatal(err)
    }

    ui.Window.SetTitle("Quicksilver - Chatting with " + peer)

    var(
        s        int
        r        int

        myPriv   []byte
        myX     *big.Int
        myY     *big.Int
    )

    myPriv, myX, myY, err = elliptic.GenerateKey(crypto.Curve, rand.Reader)
    if err != nil {
        log.Fatal(err)
    }
    myPub := elliptic.Marshal(crypto.Curve, myX, myY)

    _, err = db.LookupPubKey(peer, user, 0)
    if err == sql.ErrNoRows {
        s = 0
        r = 0

        outfile := os.Getenv("HOME") + "/" + user + ".ecdh"
        err = ioutil.WriteFile(outfile, myPub, 0666)
        if err != nil {
            log.Fatal(err)
        }

        title := "Choose " + peer + "'s public key."
        fc, err := gtk.FileChooserDialogNewWith2Buttons(
            title,
            ui.Window,
            gtk.FILE_CHOOSER_ACTION_OPEN,
            "Open", gtk.RESPONSE_OK,
            "Cancel", gtk.RESPONSE_CLOSE)
        if err != nil {
            log.Fatal(err)
        }

        if fc.Run() == gtk.RESPONSE_OK {
            file := fc.GetFilename()
            fc.Destroy()

            peerKey, err := ioutil.ReadFile(file)
            if err != nil {
                log.Fatal(err)
            }

            err = db.UploadPubKey(peer, user, peerKey, r)
            if err != nil {
                log.Fatal(err)
            }
        } else {
            log.Fatal("No public key was selected.")
        }
    } else if err != nil {
        log.Fatal(err)
    }

    s = db.LatestPrivKey(user, peer)

    err = db.UploadPrivKey(user, peer, myPriv, s)
    if err != nil {
        log.Fatal(err)
    }

    mux := new(sync.Mutex)
    go MessagePoll(jwt, user, peer, &s, &r, mux, ui)

    ui.Callback = func(msg string) {
        nxtPriv, nxtX, nxtY, err := elliptic.GenerateKey(crypto.Curve, rand.Reader)
        if err != nil {
            log.Fatal(err)
        }
        nxtPub := elliptic.Marshal(crypto.Curve, nxtX, nxtY)
        err = db.UploadPrivKey(user, peer, nxtPriv, s + 1)
        if err != nil {
            log.Fatal(err)
        }

        peerKey, err := db.LookupPubKey(peer, user, r)
        if err != nil {
            log.Fatal(err)
        }
        peerX, peerY := elliptic.Unmarshal(crypto.Curve, peerKey)
        if peerX == nil {
            log.Fatal("Invalid key data")
        }

        sessionKey := crypto.ECDH(myPriv, peerX, peerY)
        if err != nil {
            log.Fatal(err)
        }

        mux.Lock()
        c, err := crypto.EncryptMessage([]byte(msg), sessionKey, nxtPub, s, r)
        s += 1
        mux.Unlock()
        myPriv = nxtPriv

        if err != nil {
            log.Println(err)
        }

        payload, err := json.Marshal(c)
        if err != nil {
            log.Println(err)
        }

        err = api.MessageSend(jwt, peer, string(payload))
        if err != nil {
            log.Println(err)
        }
    }

    gtk.Main()
}

func MessagePoll(jwt []byte, user, peer string, s, r *int, mux *sync.Mutex, ui *gui.UI) {
    timestamp := ""
    for {
        messages, err := api.MessageFetch(jwt, peer, timestamp)
        if err != nil {
            log.Println(err)
        } else {
            for i := 0; i < len(messages); i++ {
                message := messages[i]
                timestamp = message.Timestamp

                sender := user == message.Username

                var sessionKey []byte
                if sender {
                    priv, err := db.LookupPrivKey(user, peer, message.Message.Sid)
                    if err != nil {
                        log.Fatal(err)
                    }

                    pub, err := db.LookupPubKey(peer, user, message.Message.Rid)
                    if err != nil {
                        log.Fatal(err)
                    }

                    x, y := elliptic.Unmarshal(crypto.Curve, pub)
                    if x == nil {
                        log.Fatal("Invalid public key.")
                    }

                    sessionKey = crypto.ECDH(priv, x, y)
                } else {
                    priv, err := db.LookupPrivKey(user, peer, message.Message.Rid)
                    if err != nil {
                        log.Fatal(err)
                    }

                    pub, err := db.LookupPubKey(peer, user, message.Message.Sid)
                    if err != nil {
                        log.Fatal(err)
                    }

                    x, y := elliptic.Unmarshal(crypto.Curve, pub)
                    if x == nil {
                        log.Fatal("Invalid public key.")
                    }

                    sessionKey = crypto.ECDH(priv, x, y)
                }

                clearText, nextKey, err := message.Message.Decrypt(sessionKey)
                if err != nil {
                    log.Println(err)
                }

                if !sender {
                    *r += 1
                    db.UploadPubKey(peer, user, nextKey, *r)
                }

                output := map[string]string{
                    "Username": message.Username,
                    "Timestamp": message.Timestamp,
                    "Message": string(clearText),
                }

                glib.IdleAdd(ui.ShowMessage, output)
            }
        }
        time.Sleep(time.Second)
    }
}
