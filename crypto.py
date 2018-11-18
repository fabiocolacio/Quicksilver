#!/usr/bin/python

import os, sys, json
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives.asymmetric import rsa, padding as apadding 
from cryptography.hazmat.primitives import hashes, hmac, padding, serialization
from cryptography.hazmat.backends import default_backend
from cryptography.exceptions import InvalidSignature

AES_BLOCKSIZE = 16
AES_KEYSIZE = 32
HMAC_KEYSIZE = 32

def encrypt(plaintext, keydata):
    """
    Encrypts and encapsulates a plaintext message into a JSON object.

    The process is as follows:

    1. Encrypt the plaintext with a AES-256-CBC with a random key.
    2. Generate a MAC with HMAC-SHA-256 with a random key.
    3. Concatenate the AES and HMAC keys.
    4. Encrypt the concatenated keys with an RSA public key.

    The output is a JSON object with the following structure:

    {
        "Key": <encrypted concatenated AES and HMAC keys>,
        "IV":  <IV used by the CBC>,
        "Msg": <encrypted message>,
        "Tag": <MAC tag for the encrypted message>
    }

    All of the values in the JSON object are base64-encoded, and must
    be decoded before they can be used.

    Args:
        plaintext (str): The message to encrypt
        keydata (str): The pem-encoded data for an RSA public key

    Returns:
        str: A string representation of the encrypted message encapsulated in a JSON object.
    """
    
    # Decode pem-encoded public key
    pubkey = serialization.load_pem_public_key(keydata, backend=default_backend())

    # Pad the message with PKCS7
    padder = padding.PKCS7(AES_BLOCKSIZE * 8).padder()
    padded_plaintext = padder.update(plaintext)
    padded_plaintext += padder.finalize()

    # Encrypt the Plaintext with AES-256-CBC
    aes_key = os.urandom(AES_KEYSIZE)
    iv = os.urandom(AES_BLOCKSIZE)
    cipher = Cipher(algorithms.AES(aes_key), modes.CBC(iv), backend=default_backend())
    encryptor = cipher.encryptor()
    ciphertext = encryptor.update(padded_plaintext) + encryptor.finalize()

    # Create a MAC tag using HMAC
    hmac_key = os.urandom(HMAC_KEYSIZE)
    mac = hmac.HMAC(hmac_key, hashes.SHA256(), backend=default_backend())
    mac.update(ciphertext)
    tag = mac.finalize()

    # Concatenate and encrypt AES and HMAC keys
    keys = aes_key + hmac_key
    encrypted_keys = pubkey.encrypt(
        keys,
        apadding.OAEP(
            mgf=apadding.MGF1(algorithm=hashes.SHA256()),
            algorithm=hashes.SHA256(),
            label=None))

    # Encode message into a JSON structure
    out = {
        "Key": encrypted_keys.encode("base64"),
        "IV": iv.encode("base64"),
        "Tag": tag.encode("base64"),
        "Msg": ciphertext.encode("base64")
    }

    return json.dumps(out)

def decrypt(ciphertext, keydata):
    """
    This function decrypts and decodes a JSON-encoded message that has been
    encrypted with the matching encrypt() function.

    The process is as follows:

    1. Decode all base64-encoded keys of the JSON object.
    2. Decrypt the concatenated AES and HMAC keys using the provided RSA private key.
    3. Verify message integrity by generating an HMAC tag and comparing it with.
       the one provided by the JSON object (raises InvalidSignature exception if this fails).
    4. Decrypt the ciphertext using the AES key and IV.
    5. Unpad the decrypted message with PKCS7.

    Args:
        ciphertext (str): The JSON-encoded message to decrypt
        keydata (str): The pem-encoded data for an RSA private key

    Returns:
        str: The decrypted message as a string

    Raises:
        InvalidSignature: The HMAC tag provided by the JSON object did not match ours.
    """

    # Decode pem-encoded RSA private key
    privkey = serialization.load_pem_private_key(keydata, password=None, backend=default_backend())

    # Decode the base64 encoded JSON object values
    struct = json.loads(ciphertext)
    encrypted_keys = struct["Key"].decode("base64")
    tag = struct["Tag"].decode("base64")
    iv = struct["IV"].decode("base64")
    ciphertext = struct["Msg"].decode("base64")

    # Decrypt and extract the encrypted AES and HMAC keys
    keys = privkey.decrypt(
        encrypted_keys,
        apadding.OAEP(
            mgf=apadding.MGF1(algorithm=hashes.SHA256()),
            algorithm=hashes.SHA256(),
            label=None))
    aes_key = keys[:AES_KEYSIZE]
    hmac_key = keys[AES_KEYSIZE:]

    # Verify the integrity of the ciphertext message with HMAC
    mac = hmac.HMAC(hmac_key, hashes.SHA256(), backend=default_backend())
    mac.update(ciphertext)
    try:
        mac.verify(tag)
    except InvalidSignature:
        raise

    # Decrypt the message with the AES key and IV
    cipher = Cipher(algorithms.AES(aes_key), modes.CBC(iv), backend=default_backend())
    decryptor = cipher.decryptor()
    padded_plaintext = decryptor.update(ciphertext) + decryptor.finalize()

    # Unpad the message with PKCS7
    unpadder = padding.PKCS7(AES_BLOCKSIZE * 8).unpadder()

    return unpadder.update(padded_plaintext) + unpadder.finalize()

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("USAGE: {} <message> <key> e/d".format(sys.argv[0]))
        sys.exit(0)

    plaintext = sys.argv[1]
    keypath = sys.argv[2]
    mode = sys.argv[3]
    keyfile = open(keypath, "rb")
    keydata = keyfile.read() 
    keyfile.close()

    if mode == "d":
        print(decrypt(plaintext, keydata))
    else:
        print(encrypt(plaintext, keydata))
 
