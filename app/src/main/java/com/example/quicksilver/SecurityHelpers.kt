package com.example.quicksilver

import android.util.Log
import org.spongycastle.jce.ECNamedCurveTable;

import java.math.BigInteger;
import java.security.InvalidAlgorithmParameterException;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.SecureRandom;
import java.security.interfaces.ECPrivateKey;
import java.security.interfaces.ECPublicKey;
import java.security.spec.ECGenParameterSpec
import java.security.spec.ECParameterSpec

public class SecurityHelpers {

    companion object {
        private val TAG = SecurityHelpers::class.qualifiedName

        fun generateECKeyPair(curve : String) : KeyPair? {
            var ecGenSpec = ECGenParameterSpec(curve)
            try {
                var g = KeyPairGenerator.getInstance("EC")
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

        fun getEcPublicKeyAsHex(publicKey : PublicKey) : String {

            var ecPublicKey = publicKey as ECPublicKey
            var ecPoint = ecPublicKey.getW();

            Log.d(TAG, "Length of x: " + ecPoint.getAffineX().toByteArray().size);
            Log.d(TAG, "Length of y: " + ecPoint.getAffineY().toByteArray().size);

            return byteArrayToHexString(ecPoint.getAffineX().toByteArray()) +
                    byteArrayToHexString(ecPoint.getAffineY().toByteArray())
        }

        fun getEcPrivateKeyAsHex(privateKey : PrivateKey) : String {

            var ecPrivateKey = privateKey as ECPrivateKey;
            var ecPoint = ecPrivateKey.getS();

            var privateKeyBytes = ecPoint.toByteArray();

            return byteArrayToHexString(privateKeyBytes)

        }

        fun byteArrayToHexString(arr : ByteArray) : String {
            var ret = StringBuilder()
            for (byte in arr) {
                ret.append(String.format("%02X", byte))
            }
            return ret.toString()
        }
    }

}
