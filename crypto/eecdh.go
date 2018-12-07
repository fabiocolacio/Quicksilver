package crypto

import(
    "crypto/elliptic"
    "crypto/cipher"
    "crypto/aes"
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "golang.org/x/crypto/pbkdf2"
    "errors"
    "math/big"
    "fmt"
)

const(
    aesKeySize = 32
    hmacKeySize = 32
    keyHashIters = 4096
)

var(
    secureHash = sha256.New

    Curve = elliptic.P521()

    ErrUnexpectedMAC = errors.New("Computed and expected MAC tags do not match.")
)

type EncryptedMessage struct {
    Ax  []byte // Elliptic-Curve X-value of the message sender
    Ay  []byte // Elliptic-Curve Y-value of the message sender
    Bx  []byte // Elliptic-Curve X-value of the message receiver
    By  []byte // Elliptic-Curve Y-value of the message receiver
    IV  []byte // AES IV used to encrypt the message and HMAC key
    Msg []byte // AES encrypted message data
    Key []byte // AES encrypted HMAC key
    Tag []byte // HMAC integrity tag
}

// DeriveKey creates a key of size keysize from binary data.
func DeriveKey(mother []byte, keysize int) []byte {
    return pbkdf2.Key(mother, nil, 4096, keysize, secureHash)
}

// Encrypt encrypts clearText using a shared secret acquired through an
// elliptic-curve diffie-hellman key exchange.
//
// Your private diffie-hellman information, priv, is used with the peer's
// public diffie-hellman information (bx, by), to create a shared AES session
// key to encrypt clearText with. Returns an EncryptedMessage and an error.
func EncryptMessage(clearText, priv []byte, ax, ay, bx, by *big.Int) (msg *EncryptedMessage, err error) {
    // Create shared secret xp from peer's public key and our private key
    xp, _ := Curve.ScalarMult(bx, by, priv)

    // Derive an AES key from our shared secret
    aesKey := DeriveKey(xp.Bytes(), aesKeySize)
    fmt.Println(aesKey)

    // Create a random HMAC key
    hmacKey := make([]byte, hmacKeySize)
    if _, err := rand.Read(hmacKey); err != nil {
        return nil, err
    }

    // Add PKCS7 padding to clearText
    paddedClearText, err := pkcs7Pad(clearText, aes.BlockSize)
    if err != nil {
        return nil, err
    }

    // Create buffers for ciphertexts
    cipherText := make([]byte, len(paddedClearText))
    encryptedKey := make([]byte, len(hmacKey))

    // Create AES block cipher
    aesCipher, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }

    // Create a random initialization vector for AES encryption
    iv := make([]byte, aes.BlockSize)
    if _, err = rand.Read(iv); err != nil {
        return nil, err
    }

    // Encrypt data with CBC block encrypter
    cbc := cipher.NewCBCEncrypter(aesCipher, iv)
    cbc.CryptBlocks(cipherText, paddedClearText)

    // Encrypt hmac key with CBC block encrypter
    cbc = cipher.NewCBCEncrypter(aesCipher, iv)
    cbc.CryptBlocks(encryptedKey, hmacKey)

    // Generate MAC tag for data
    mac := hmac.New(secureHash, hmacKey)
    mac.Write(cipherText)
    tag := mac.Sum(nil)

    fmt.Println("hmac key", hmacKey)

    msg = &EncryptedMessage{
        Ax: ax.Bytes(),
        Ay: ay.Bytes(),
        Bx: bx.Bytes(),
        By: by.Bytes(),
        IV: iv,
        Msg: cipherText,
        Tag: tag,
        Key: encryptedKey,
    }

    return msg, err
}

func (message *EncryptedMessage) Decrypt(priv []byte, sender bool) ([]byte, error) {
    x := new(big.Int)
    y := new(big.Int)

    if sender {
        x.SetBytes(message.Ax)
        y.SetBytes(message.Ay)
    } else {
        x.SetBytes(message.Bx)
        y.SetBytes(message.By)
    }

    // Create shared secret xp from peer's public key and our private key
    xp, _ := Curve.ScalarMult(x, y, priv)

    // Derive an AES key from our shared secret
    aesKey := DeriveKey(xp.Bytes(), aesKeySize)
    fmt.Println(aesKey)

    // Create AES block cipher
    aesCipher, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }

    fmt.Println("hmac key", message.Key)

    // Decrypt HMAC Key
    cbc := cipher.NewCBCDecrypter(aesCipher, message.IV)
    cbc.CryptBlocks(message.Key, message.Key)
    fmt.Println("hmac key", message.Key)

    // Compare MAC tags
    if !CheckMAC(message.Msg, message.Tag, message.Key) {
        return nil, ErrUnexpectedMAC
    }

    // Decrypt and unpad the payload
    cbc = cipher.NewCBCDecrypter(aesCipher, message.IV)
    cbc.CryptBlocks(message.Msg, message.Msg)
    msg, err := pkcs7Unpad(message.Msg, aes.BlockSize)
    if err != nil {
        return nil, err
    }

    return msg, err
}

func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
