package com.example.quicksilver

import java.security.*
import java.security.spec.ECGenParameterSpec

class SecurityHelpers {
    companion object {
        fun generateECKeyPair(curve : String) : KeyPair? {
            val ecGenSpec = ECGenParameterSpec(curve)
            try {
                val g = KeyPairGenerator.getInstance("EC")
                g.initialize(ecGenSpec, SecureRandom())
                return g.generateKeyPair()
            }
            catch (e : NoSuchAlgorithmException) {
                e.printStackTrace()
            }
            catch (e : InvalidAlgorithmParameterException) {
                e.printStackTrace()
            }
            catch (e : NullPointerException) {
                e.printStackTrace()
            }
            return null
        }

        fun generateECSignature(key : PrivateKey, data : ByteArray) : ByteArray {
            val sig = Signature.getInstance("SHA256withECDSA")
            sig.initSign(key)
            sig.update(data)
            return sig.sign()
        }

        fun validateECSignature(key : PublicKey, data : ByteArray, signature : ByteArray) : Boolean {
            val sig = Signature.getInstance("SHA256withECDSA")
            sig.initVerify(key)
            sig.update(data)
            return sig.verify(signature)
        }
    }
}
