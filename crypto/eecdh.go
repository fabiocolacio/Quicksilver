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
    Sid  int   // The index of this sender's key to use in this diffie-hellman
    Rid  int   // The index of the receiver's key to use in this diffie-hellman
    Nxt []byte // Elliptic-Curve public data for the next message (encrypted)
    IV  []byte // AES IV used to encrypt the message and HMAC key
    Msg []byte // AES encrypted message data
    Key []byte // AES encrypted HMAC key
    Tag []byte // HMAC integrity tag
}

// DeriveKey creates a key of size keysize from binary data.
func DeriveKey(mother []byte, keysize int) []byte {
    return pbkdf2.Key(mother, nil, 4096, keysize, secureHash)
}

func ECDH(priv []byte, x, y *big.Int) []byte {
    // Create shared secret xp from peer's public key and our private key
    xp, _ := Curve.ScalarMult(x, y, priv)

    // Derive an AES key from our shared secret
    return DeriveKey(xp.Bytes(), aesKeySize)
}

// Encrypt encrypts clearText using a shared secret acquired through an
// elliptic-curve diffie-hellman key exchange.
//
// Your private diffie-hellman information, priv, is used with the peer's
// public diffie-hellman information (bx, by), to create a shared AES session
// key to encrypt clearText with. Returns an EncryptedMessage and an error.
func EncryptMessage(clearText, aesKey, nxt []byte, sid, rid int) (msg *EncryptedMessage, err error) {
    // Create a random HMAC key
    hmacKey := make([]byte, hmacKeySize)
    if _, err := rand.Read(hmacKey); err != nil {
        return nil, err
    }

    fmt.Println(nxt)

    // Add PKCS7 padding to clearText
    paddedClearText, err := PKCS7Pad(clearText, aes.BlockSize)
    if err != nil {
        return nil, err
    }

    // Add PKCS7 padding to next key
    nxt, err = PKCS7Pad(nxt, aes.BlockSize)
    if err != nil {
        return nil, err
    }

    // Create buffers for ciphertexts
    cipherText := make([]byte, len(paddedClearText))
    encryptedKey := make([]byte, len(hmacKey))
    encryptedNxt := make([]byte, len(nxt))

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

    // Encrypt nxt key with CBC block encrypter
    cbc = cipher.NewCBCEncrypter(aesCipher, iv)
    cbc.CryptBlocks(encryptedNxt, nxt)

    // Generate MAC tag for data
    mac := hmac.New(secureHash, hmacKey)
    mac.Write(cipherText)
    tag := mac.Sum(nil)

    msg = &EncryptedMessage{
        Sid: sid,
        Rid: rid,
        Nxt: encryptedNxt,
        IV:  iv,
        Msg: cipherText,
        Tag: tag,
        Key: encryptedKey,
    }

    return msg, err
}

func (message *EncryptedMessage) Decrypt(aesKey []byte) (clearText, nextKey []byte, err error) {
    // Create AES block cipher
    aesCipher, err := aes.NewCipher(aesKey)
    if err != nil {
        return
    }

    // Decrypt HMAC Key
    cbc := cipher.NewCBCDecrypter(aesCipher, message.IV)
    cbc.CryptBlocks(message.Key, message.Key)

    // Compare MAC tags
    if !CheckMAC(message.Msg, message.Tag, message.Key) {
        err = ErrUnexpectedMAC
        return
    }

    // Decrypt and unpad the payload
    cbc = cipher.NewCBCDecrypter(aesCipher, message.IV)
    cbc.CryptBlocks(message.Msg, message.Msg)
    msg, err := PKCS7Unpad(message.Msg, aes.BlockSize)
    if err != nil {
        return
    }

    // Decrypt and unpad the next key
    cbc = cipher.NewCBCDecrypter(aesCipher, message.IV)
    cbc.CryptBlocks(message.Nxt, message.Nxt)
    nxt, err := PKCS7Unpad(message.Nxt, aes.BlockSize)
    if err != nil {
        return
    }

    fmt.Println(nxt)

    return msg, nxt, err
}

func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
