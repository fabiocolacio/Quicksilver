package com.example.quicksilver.testing

import android.os.Handler
import android.os.Looper
import android.os.Message
import android.os.SystemClock.uptimeMillis
import android.util.Log
import com.example.quicksilver.SecurityHelpers
import java.security.KeyFactory
import java.security.PublicKey
import java.security.spec.KeySpec
import java.security.spec.PKCS8EncodedKeySpec
import java.security.spec.X509EncodedKeySpec

class FakeServer(mainThreadLooper: Looper, simulatedDelay: Int) {

    companion object {

        private val TAG = FakeServer.javaClass.canonicalName

        val LOGIN_TEST_VALUE = ByteArray(10) { init : Int -> (init+1).toByte()}

        // Messages
        private const val MSG_REGISTER_KEY = 0
        private const val MSG_UNREGISTER_KEY = 1
        private const val MSG_LOGIN = 2
    }

    class RequestInfo(var s1: String?, var s2: String?, var o1: Any?,
                      var callback : (Boolean) -> Unit)

    var keys : HashMap<String, PublicKey> = HashMap()
    var kf : KeyFactory = KeyFactory.getInstance("EC")
    var handler : Handler
    var simulatedDelayMs : Int = simulatedDelay

    init {
        handler = Handler(mainThreadLooper) {
            val ri = it.obj as RequestInfo
            when (it.what) {
                MSG_REGISTER_KEY -> {
                    registerKey(
                        ri.s1 as String,
                        ri.o1 as ByteArray,
                        ri.s2 as String,
                        ri.callback
                    )
                    true
                }
                MSG_UNREGISTER_KEY -> {
                    unregisterKey(
                        ri.s1 as String,
                        ri.callback
                    )
                    true
                }
                MSG_LOGIN -> {
                    login(
                        ri.s1 as String, ri.o1 as ByteArray,
                        ri.callback
                    )
                    true
                }
                else -> {
                    Log.e(TAG, "Unrecognized message type: " + it.what)
                    true
                }
            }
        }
    }

    fun requestRegisterKey(userName : String, publicKeyBytes : ByteArray, format : String,
                           callback : (Boolean) -> Unit) {
        val msg = Message()
        msg.what = MSG_REGISTER_KEY
        msg.obj = RequestInfo(userName, format, publicKeyBytes, callback)
        handler.sendMessageAtTime(msg, uptimeMillis() + simulatedDelayMs)
    }

    fun requestUnregisterKey(userName : String, callback : (Boolean) -> Unit) {
        val msg = Message()
        msg.what = MSG_UNREGISTER_KEY
        msg.obj = RequestInfo(userName, null,null, callback)
        handler.sendMessageAtTime(msg, uptimeMillis() + simulatedDelayMs)
    }

    fun requestLogin(userName : String, testSignature : ByteArray,
                     callback : (Boolean) -> Unit) {
        val msg = Message()
        msg.what = MSG_LOGIN
        msg.obj = RequestInfo(userName, null, testSignature, callback)
        handler.sendMessageAtTime(msg, uptimeMillis() + simulatedDelayMs)
    }

    private fun registerKey(userName : String, publicKeyBytes : ByteArray, format : String,
                            callback : (Boolean) -> Unit) {
        if (keys.containsKey(userName)) {
            callback(false)
            return
        }
        val spec : KeySpec
        if (format == "X.509") {
            spec = X509EncodedKeySpec(publicKeyBytes)
        }
        else if (format == "PKCS#8") {
            spec = PKCS8EncodedKeySpec(publicKeyBytes)
        }
        else {
            callback(false)
            return
        }
        keys.put(userName, kf.generatePublic(spec))
        callback(true)
    }

    private fun unregisterKey(userName : String, callback : (Boolean) -> Unit) {
        if (!keys.containsKey(userName)) {
            callback(false)
            return
        }
        keys.remove(userName)
        callback(true)
    }

    private fun login(userName : String, testSignature : ByteArray,
                      callback : (Boolean) -> Unit) {
        val key: PublicKey? = keys[userName]
        if (key == null) {
            callback(false)
            return
        }
        if (SecurityHelpers.validateECSignature(
                key, LOGIN_TEST_VALUE, testSignature)) {
            callback(true)
            return
        }
        callback(false)
    }


}