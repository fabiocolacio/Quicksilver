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

    var(
        myPriv   []byte
        myX     *big.Int
        myY     *big.Int
        peerX   *big.Int
        peerY   *big.Int
    )

    myPriv, myX, myY, err = elliptic.GenerateKey(crypto.Curve, rand.Reader)
    if err != nil {
        log.Fatal(err)
    }
    myPub := elliptic.Marshal(crypto.Curve, myX, myY)

    peerKey, err := db.LookupPubKey(peer)
    if err == sql.ErrNoRows {
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

            peerKey, err = ioutil.ReadFile(file)
            if err != nil {
                log.Fatal(err)
            }

            err = db.UploadKey(peer, peerKey, nil)
            if err != nil {
                log.Fatal(err)
            }
        } else {
            log.Fatal("No public key was selected.")
        }
    } else if err != nil {
        log.Fatal(err)
    }

    err = db.UploadKey(user, myPub, myPriv)
    if err != nil {
        log.Fatal(err)
    }

    peerX, peerY = elliptic.Unmarshal(crypto.Curve, peerKey)
    if peerX == nil {
        log.Fatal("Invalid key data")
    }

    go MessagePoll(jwt, user, peer, ui)

    ui.Callback = func(msg string) {
        c, err := crypto.EncryptMessage([]byte(msg), myPriv, myX, myY, peerX, peerY)
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

        myPriv, myX, myY, err = elliptic.GenerateKey(crypto.Curve, rand.Reader)
        if err != nil {
            log.Fatal(err)
        }

        myPub := elliptic.Marshal(crypto.Curve, myX, myY)
        err = db.UploadKey(user, myPub, myPriv)
        if err != nil {
            log.Fatal(err)
        }
    }

    gtk.Main()
}

func MessagePoll(jwt []byte, user, peer string, ui *gui.UI) {
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

                x := new(big.Int)
                y := new(big.Int)

                if sender {
                    x.SetBytes(message.Message.Ax)
                    y.SetBytes(message.Message.Ay)
                } else {
                    x.SetBytes(message.Message.Bx)
                    y.SetBytes(message.Message.By)
                }

                pubKey := elliptic.Marshal(crypto.Curve, x, y)

                if pubKey == nil {
                    log.Println("booty")
                }

                privKey, err := db.LookupPrivKey(user, pubKey)
                if err != nil {
                    log.Println(err)
                }

                clearText, err := message.Message.Decrypt(privKey, sender)
                if err != nil {
                    log.Println(err)
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
