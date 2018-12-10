package api

import(
    "net/http"
    "encoding/json"
    "errors"
    "bytes"
    "io/ioutil"
    "crypto/tls"
    "crypto/hmac"
    "crypto/sha256"
    "net/url"
    "log"
    "github.com/fabiocolacio/quicksilver/crypto"
    "golang.org/x/crypto/scrypt"
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

// Message represents Mercury's internal representation
// of a message
type Message struct {
    Username  string
    Timestamp string
    Message   crypto.EncryptedMessage
}

// MessageFetch checks for messages from peer through the Mercury Server
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

// MessageSend sends message to peer using the Mercrury server
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

// SetHost changes the host server that will be used by Quicksilver
func SetHost(newHost string) {
    host = newHost
}

// LookupUser checks if a user exists in Mercury's database
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

// Register registers a user with the Mercury server
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

// Login attempts to login to the Mercury Server.
// If it is successful, it returns the JWT.
func Login(user, passwd string) ([]byte, error) {
    res, err := client.Get(host + "/login?user=" + user)
    challenge := make(map[string][]byte)
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(body, &challenge)
    res.Body.Close()

    saltedHash, err := scrypt.Key([]byte(passwd), challenge["S"], 32768, 8, 1, 32)
    if err != nil {
        return nil, err
    }

    mac := hmac.New(sha256.New, saltedHash)
    mac.Write(challenge["C"])
    payload := mac.Sum(nil)

    res, err = client.Post(host + "/auth?user=" + user, "text/javascript", bytes.NewBuffer(payload))
    if err != nil {
        return nil, err
    }
    if res.StatusCode != 200 {
        log.Println(res.StatusCode)
        return nil, ErrLoginFailed
    }

    defer res.Body.Close()
    jwt, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    return jwt, err
}
