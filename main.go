package main

import(
    // "fmt"
    // "github.com/fabiocolacio/quicksilver/crypto"
    // "crypto/elliptic"
    // "crypto/rand"
    "os"
    "log"
    "time"
    "github.com/fabiocolacio/quicksilver/gui"
    "github.com/fabiocolacio/quicksilver/api"
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
)

func main() {
    config := os.Getenv("HOME") + "/.config/quicksilver"
    err := os.MkdirAll(config, 0666)
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

    // sess, err := api.UnwrapJWT(jwt)
    // if err != nil {
    //     log.Fatal(err)
    // }

    peer, err := gui.PeerDialogRun(ui.Window)
    if err != nil {
        log.Fatal(err)
    }

    if err = api.LookupUser(peer); err != nil {
        log.Fatal(err)
    }

    // priv, x, y, err := elliptic.GenerateKey(crypto.Curve, rand.Reader)
    // if err != nil {
    //     log.Fatal(err)
    // }

    go MessagePoll(jwt, peer, ui)
    
    ui.Callback = func(msg string) {
        err := api.MessageSend(jwt, peer, msg)
        if err != nil {
            log.Println(err)
        }
    }

    gtk.Main()
}

func MessagePoll(jwt []byte, peer string, ui *gui.UI) {
    timestamp := ""
    for {
        messages, err := api.MessageFetch(jwt, peer, timestamp)
        if err != nil {
            log.Println(err)
        } else {
            for i := 0; i < len(messages); i++ {
                message := messages[i]
                timestamp = message["Timestamp"]
                glib.IdleAdd(ui.ShowMessage, message)
            }
        }
        time.Sleep(2 * time.Second)
    }
}
