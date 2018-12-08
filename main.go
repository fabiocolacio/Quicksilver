package main

import(
    // "fmt"
    // "github.com/fabiocolacio/quicksilver/crypto"
    // "crypto/elliptic"
    // "crypto/rand"
    // "os"
    "log"
    "github.com/fabiocolacio/quicksilver/gui"
    // "github.com/fabiocolacio/quicksilver/api"
    "github.com/gotk3/gotk3/gtk"
)

func main() {
    gtk.Init(nil)

    ui, err := gui.UINew()
    if err != nil {
        log.Fatal(err)
    }

    user, passwd, err := gui.LoginDialogRun(ui.Window)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(user, passwd)

    gtk.Main()
}
