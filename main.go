package main

import(
    "crypto/elliptic"
    "crypto/rand"
    "io/ioutil"
    "encoding/json"
    "database/sql"
    "os"
    "log"
    "time"
    "github.com/fabiocolacio/quicksilver/crypto"
    "github.com/fabiocolacio/quicksilver/gui"
    "github.com/fabiocolacio/quicksilver/db"
    "github.com/fabiocolacio/quicksilver/api"
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
)

func main() {
    // Create database tables if they don't exist
    err := db.InitTables()
    if err != nil {
        log.Fatal(err)
    }

    gtk.Init(nil)

    // Create UI
    ui, err := gui.UINew()
    if err != nil {
        log.Fatal(err)
    }

    // Log the user in
    user, passwd, err := gui.LoginDialogRun(ui.Window)
    if err != nil {
        log.Fatal(err)
    }
    jwt, err := api.Login(user, passwd)
    if err != nil {
        log.Fatal(err)
    }

    // Select a peer to speak with
    peer, err := gui.PeerDialogRun(ui.Window)
    if err != nil {
        log.Fatal(err)
    }
    if err = api.LookupUser(peer); err != nil {
        log.Fatal(err)
    }
    ui.Window.SetTitle("Quicksilver - Chatting with " + peer)

    // Generate first ECDH key pair
    myPriv, myX, myY, err := elliptic.GenerateKey(crypto.Curve, rand.Reader)
    if err != nil {
        log.Fatal(err)
    }
    myPub := elliptic.Marshal(crypto.Curve, myX, myY)

    // Check if We have a key for our peer already
    _, err = db.LookupPubKey(peer, user, 0)

    // If we don't have our peer's key, ask the user to select it
    if err == sql.ErrNoRows {
        // Create our diffie public key file to give to our peer
        outfile := os.Getenv("HOME") + "/" + user + ".ecdh"
        err = ioutil.WriteFile(outfile, myPub, 0666)
        if err != nil {
            log.Fatal(err)
        }

        // Create File Chooser Dialog
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

        // Prompt user to choose peer's public key
        if fc.Run() == gtk.RESPONSE_OK {
            file := fc.GetFilename()
            fc.Destroy()

            // Read the key data
            peerKey, err := ioutil.ReadFile(file)
            if err != nil {
                log.Fatal(err)
            }

            // Store the key in the database
            err = db.UploadPubKey(peer, user, peerKey, 0)
            if err != nil {
                log.Fatal(err)
            }
        } else {
            log.Fatal("No public key was selected.")
        }
    } else if err != nil {
        log.Fatal(err)
    }

    // Store our private key in the database
    err = db.UploadPrivKey(user, peer, myPriv, db.LatestPrivKey(user, peer))
    if err != nil {
        log.Fatal(err)
    }

    // Poll for messages in parallel
    go MessagePoll(jwt, user, peer, ui)

    // Callback function for sending messages when the user pressed Enter
    ui.Callback = func(msg string) {
        // Get peer's latest public key id
        r := db.LatestPubKey(peer, user)
        // Get our latest private key id
        s := db.LatestPrivKey(user, peer)

        // Generate the next ECDH key
        nxtPriv, nxtX, nxtY, err := elliptic.GenerateKey(crypto.Curve, rand.Reader)
        if err != nil {
            log.Fatal(err)
        }
        nxtPub := elliptic.Marshal(crypto.Curve, nxtX, nxtY)

        // Store the private component for ourselves
        err = db.UploadPrivKey(user, peer, nxtPriv, s + 1)
        if err != nil {
            log.Fatal(err)
        }

        // Get the peer's public key
        peerKey, err := db.LookupPubKey(peer, user, r)
        if err != nil {
            log.Println("e")
            log.Fatal(err)
        }
        peerX, peerY := elliptic.Unmarshal(crypto.Curve, peerKey)
        if peerX == nil {
            log.Fatal("Invalid key data")
        }

        // Get our latest private key
        myKey, err := db.LookupPrivKey(user, peer, s)
        if err != nil {
            log.Println("f")
            log.Fatal(err)
        }

        // Perform ECDH to make session key
        sessionKey := crypto.ECDH(myKey, peerX, peerY)

        // Encrypt the message with the session key
        c, err := crypto.EncryptMessage([]byte(msg), sessionKey, nxtPub, s, r)
        if err != nil {
            log.Println(err)
        }

        // Convert the message into a JSON string
        payload, err := json.Marshal(c)
        if err != nil {
            log.Println(err)
        }

        // Send the message to the server
        err = api.MessageSend(jwt, peer, string(payload))
        if err != nil {
            log.Println(err)
        }
    }

    gtk.Main()
}

// MessagePoll polls the server for new messages, decrypts them as they come,
// and shows them in the GUI. It also stores new public keys into the database.
func MessagePoll(jwt []byte, user, peer string, ui *gui.UI) {
    timestamp := ""
    for {
        // Get list of new messages
        messages, err := api.MessageFetch(jwt, peer, timestamp)
        if err != nil {
            log.Println(err)
        } else {
            // Process each message in the list
            for i := 0; i < len(messages); i++ {
                message := messages[i]

                timestamp = message.Timestamp

                // Check if we are the sender of the message
                sender := user == message.Username

                var(
                    privId int // Our private key id to use
                    pubId  int // Peer's public key id to use
                )

                // Fetch appropriate ECDH parameters
                if sender {
                    privId = message.Message.Sid
                    pubId = message.Message.Rid
                } else {
                    privId = message.Message.Rid
                    pubId = message.Message.Sid
                }
                priv, err := db.LookupPrivKey(user, peer, privId)
                if err != nil {
                    log.Println("a")
                    log.Fatal(err)
                }
                pub, err := db.LookupPubKey(peer, user, pubId)
                if err != nil {
                    log.Println("b")
                    log.Fatal(err)
                }
                x, y := elliptic.Unmarshal(crypto.Curve, pub)
                if x == nil {
                    log.Fatal("Invalid public key.")
                }

                // Perform ECDH to get session key
                sessionKey := crypto.ECDH(priv, x, y)

                // Decrypt message with session key
                clearText, nextKey, err := message.Message.Decrypt(sessionKey)
                if err != nil {
                    log.Println(err)
                }

                // If we are the receiver, store the next public key
                if !sender {
                    err = db.UploadPubKey(peer, user, nextKey, message.Message.Sid + 1)
                    if err != nil {
                        log.Println(err)
                    }
                }

                // Send the message to the UI thread to be displayed
                output := map[string]string{
                    "Username": message.Username,
                    "Timestamp": message.Timestamp,
                    "Message": string(clearText),
                }
                glib.IdleAdd(ui.ShowMessage, output)
            }
        }

        // Sleep for 1 second before polling again
        time.Sleep(time.Second)
    }
}
