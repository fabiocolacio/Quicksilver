package main

import(
    "fmt"
    
//    "github.com/fabiocolacio/quicksilver/crypto"
//    "crypto/elliptic"
//    "crypto/rand"

    "log"
    "github.com/fabiocolacio/quicksilver/gui"
    "github.com/fabiocolacio/quicksilver/api"
    "github.com/gotk3/gotk3/gtk"
)

func main() {
    jwt, err := api.Login("fabio", "fabio")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(jwt))

    gtk.Init(nil)

    ui, err := gui.UINew()
    if err != nil {
        log.Fatal(err)
    }

    ui.Window.ShowAll()

    gtk.Main()

    /*
    pa, ax, ay, err := elliptic.GenerateKey(crypto.Curve, rand.Reader)
    if err != nil {
        fmt.Println(err)
    }

    pb, bx, by, err := elliptic.GenerateKey(crypto.Curve, rand.Reader)
    if err != nil {
        fmt.Println(err)
    }

    msg := "I think, therefore I am"
    fmt.Println(msg)

    cipher, err := crypto.EncryptMessage([]byte(msg), pa, ax, ay, bx, by)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(cipher.Msg))

    clear, err := cipher.Decrypt(pb, true)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(clear))
    */
}
