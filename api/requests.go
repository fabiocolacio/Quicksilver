package api

import(
    "net/http"
    "encoding/json"
    "errors"
    "bytes"
    "io/ioutil"
    "crypto/tls"
)

var(
    host = "https://localhost:9090"

    ErrLoginFailed = errors.New("Login failed.")

    client = &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
        },
    }
)

func SetHost(newHost string) {
    host = newHost
}

func Login(user, passwd string) ([]byte, error) {
    creds := map[string]string{
        "Username": user,
        "Password": passwd,
    }

    payload, err := json.Marshal(creds)
    if err != nil {
        return nil, err
    }

    res, err := client.Post(host + "/login", "text/javascript", bytes.NewBuffer(payload))
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    jwt, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    return jwt, err
}
