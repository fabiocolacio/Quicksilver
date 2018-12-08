package api

import(
    "net/http"
    "encoding/json"
    "errors"
    "bytes"
    "io/ioutil"
    "crypto/tls"
    "net/url"
    "log"
)

var(
    host = "https://localhost:9090"

    ErrLoginFailed = errors.New("Login failed.")
    ErrNoSuchUser = errors.New("No such user.")

    client = &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
        },
    }
)

type Message map[string]string

func MessageFetch(jwt []byte, peer, since string) ([]Message, error) {
    uri := host + "/get?peer=" + url.QueryEscape(peer)

    if len(since) > 0 {
        uri = uri + "&since=" + url.QueryEscape(since)
    }

    request, err := http.NewRequest("GET", uri, nil)
    if err != nil {
        return nil, err
    }
    request.Header.Set("Session", string(jwt))
    res, err := client.Do(request)
    if err != nil {
        return nil, err
    }

    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    if res.StatusCode != 200 {
        return nil, errors.New(string(body))
    }

    var messages []Message
    err = json.Unmarshal(body, &messages)
    if err != nil {
        return nil, err
    }

    return messages, err
}

func MessageSend(jwt []byte, peer, message string) error {
    uri := host + "/send?to=" + url.QueryEscape(peer)
    request, err := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(message)))
    if err != nil {
        return err
    }
    request.Header.Set("Session", string(jwt))
    res, err := client.Do(request)
    if err != nil {
        return err
    }

    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return err
    }

    if res.StatusCode != 200 {
        return errors.New(string(body))
    }

    return nil
}

func SetHost(newHost string) {
    host = newHost
}

func LookupUser(user string) error {
    res, err := client.Get(host + "/lookup?user=" + user)
    if err != nil {
        return err
    }

    if res.StatusCode != 200{
        return ErrNoSuchUser
    }

    return nil
}

func Register(user, passwd string) error {
    creds := map[string]string{
        "Username": user,
        "Password": passwd,
    }

    payload, err := json.Marshal(creds)
    if err != nil {
        return err
    }

    res, err := client.Post(host + "/register", "text/javascript", bytes.NewBuffer(payload))
    if err != nil {
        return err
    }
    if res.StatusCode != 200 {
        return ErrLoginFailed
    }

    return nil
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
    if res.StatusCode != 200 {
        return nil, ErrLoginFailed
    }

    defer res.Body.Close()
    jwt, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    return jwt, err
}
