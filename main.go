package main

import(
    "fmt"
    "github.com/fabiocolacio/quicksilver/crypto"
    "crypto/elliptic"
    "crypto/rand"
)

func main() {
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
}
