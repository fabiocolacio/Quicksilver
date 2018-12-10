package api

import(
    "strings"
    "encoding/json"
    "github.com/fabiocolacio/mercury/server"
    "encoding/base64"
    "errors"
    "bytes"
)

var(
    ErrMalformedJWT = errors.New("Malformed JWT")
)

// UnwrapJWT returns the JSON payload element of a JWT
func UnwrapJWT(jwt []byte) (server.Session, error) {
    var sess server.Session
    elements := strings.Split(string(jwt), ".")
    if len(elements) < 3 {
        return sess, ErrMalformedJWT
    }
    payload := elements[1]
    jsonObj := make([]byte, base64.URLEncoding.DecodedLen(len(payload)))

    _, err := base64.URLEncoding.Decode(jsonObj, []byte(payload))
    if err != nil {
        return sess, err
    }
    jsonObj = bytes.Trim(jsonObj, "\x00")

    err = json.Unmarshal(jsonObj, &sess)
    if err != nil {
        return sess, err
    }

    return sess, err
}
