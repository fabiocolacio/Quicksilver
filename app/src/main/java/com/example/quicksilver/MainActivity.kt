package com.example.quicksilver

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.TextView
import com.example.quicksilver.testing.FakeServer
import java.security.Security
import java.text.SimpleDateFormat
import java.util.*

class MainActivity : AppCompatActivity() {

    companion object {
        private val TAG = MainActivity::class.qualifiedName

        const val FAKE_SERVER_DELAY_MS = 5000
        const val TEST_NAME = "Bob"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        Security.addProvider(org.spongycastle.jce.provider.BouncyCastleProvider())

        var textDisplay = findViewById<TextView>(R.id.text_display)
        var keyPair = SecurityHelpers.generateECKeyPair("secp256r1")
        if (keyPair == null) {
            logMessage("Unable to generate key pair.", textDisplay)
            return
        }

        var fakeServer = FakeServer(mainLooper, FAKE_SERVER_DELAY_MS)

        logMessage("Requesting to register key for $TEST_NAME at fake server...", textDisplay)
        fakeServer.requestRegisterKey(TEST_NAME, keyPair.public.encoded, keyPair.public.format) {
            logMessage("Finished registering key for $TEST_NAME to fake server.", textDisplay)
        }
    }

    private fun logMessage(msg : String, textDisplay : TextView) {
        val s = SimpleDateFormat("MM:dd:yyyy:hh:mm:ss")
        val format = s.format(Date())
        val display = "--- $format ---\n$msg\n---\n"
        textDisplay.append(display)
    }

}
