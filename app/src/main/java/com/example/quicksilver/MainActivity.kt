package com.example.quicksilver

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.os.Handler
import android.os.Message
import android.util.Log
import com.example.quicksilver.UI.Model
import com.example.quicksilver.UI.View
import com.example.quicksilver.testing.FakeServer
import java.security.KeyPair
import java.security.Security
import java.util.*

class MainActivity : AppCompatActivity() {

    companion object {
        private val TAG = MainActivity::class.qualifiedName

        const val FAKE_SERVER_DELAY_MS = 5000
        const val TEST_NAME = "Bob"

        private const val MSG_REGISTER_KEY_RESULT = 0
        private const val MSG_UNREGISTER_KEY_RESULT = 1
        private const val MSG_LOGIN_RESULT = 2

        private const val OP_REGISTER_KEY = "REGISTER_KEY"
        private const val OP_UNREGISTER_KEY = "UNREGISTER_KEY"
        private const val OP_LOGIN = "LOGIN"
    }

    private lateinit var handler : Handler
    private lateinit var model : Model
    private lateinit var view : View

    private lateinit var keyPair : KeyPair

    class OpResult(var opType : String, var userName : String, var result: Boolean)

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        Security.addProvider(org.spongycastle.jce.provider.BouncyCastleProvider())

        // initialize model and view objects
        model = Model(false, ArrayList())
        view = View(this, model)
        model.addObserver(view)

        // generate key pair, initialize fake server
        var genKeyPair = SecurityHelpers.generateECKeyPair("secp256r1")
        if (genKeyPair == null) {
            model.updateNewLogMessage("Failed to generate key pair.")
            return
        }
        else {
            keyPair = genKeyPair
        }
        var fakeServer = FakeServer(mainLooper, FAKE_SERVER_DELAY_MS)

        // set up handler to process results of asynchronous calls
        handler = Handler() {
            val opResult = it.obj as OpResult
            when (it.what) {
                MSG_REGISTER_KEY_RESULT -> {
                    processResult(opResult, fakeServer, model)
                    true
                }
                MSG_UNREGISTER_KEY_RESULT -> {
                    processResult(opResult, fakeServer, model)
                    true
                }
                MSG_LOGIN_RESULT -> {
                    processResult(opResult, fakeServer, model)
                    true
                }
                else -> {
                    model.updateNewLogMessage("Unrecognized message type to main handler: " + it.what)
                    true
                }
            }
        }

        // request to register key with fake server
        model.updateNewLogMessage("Requesting to register key for [$TEST_NAME] at fake server...")
        fakeServer.requestRegisterKey(TEST_NAME, keyPair.public.encoded, keyPair.public.format) {
            val msg = Message()
            msg.what = MSG_REGISTER_KEY_RESULT
            msg.obj = OpResult(OP_REGISTER_KEY, TEST_NAME, it)
            handler.sendMessage(msg)
        }

    }

    private fun processResult(result : OpResult, server : FakeServer, model : Model) {

        logResult(result, model)

        val userName = result.userName

        when (result.opType) {
            OP_REGISTER_KEY -> {
                model.updateNewLogMessage("Requesting to login [$TEST_NAME] at fake server...")
                server.requestLogin(
                    userName,
                    SecurityHelpers.generateECSignature(
                        keyPair.private,
                        FakeServer.LOGIN_TEST_VALUE
                    )
                ) {
                    val msg = Message()
                    msg.what = MSG_LOGIN_RESULT
                    msg.obj = OpResult(OP_LOGIN, userName, it)
                    handler.sendMessage(msg)
                }
            }
            OP_UNREGISTER_KEY -> {
                // nothing
            }
            OP_LOGIN -> {
                if (result.result) {
                    model.updateLoginStatus(true)
                }
            }
            else -> {
                model.updateNewLogMessage("Unrecognized op type: " + result.opType)
            }
        }
    }

    private fun logResult(result : OpResult, model : Model) {
        if (result.result) {
            model.updateNewLogMessage("Finished [" + result.opType + "] for [" +
                    result.userName + "] to fake server.")
        }
        else {
            model.updateNewLogMessage("Failed to [" + result.opType + "] for [" +
                    result.userName + "] to fake server.")
        }
    }

}
