package com.example.quicksilver

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.widget.TextView
import java.security.Security


class MainActivity : AppCompatActivity() {

    companion object {
        private val TAG = MainActivity::class.qualifiedName
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        Security.addProvider(org.spongycastle.jce.provider.BouncyCastleProvider())

        var textDisplay = findViewById<TextView>(R.id.text_display)
        var displayString = ""
        var keyPair = SecurityHelpers.generateECKeyPair("secp256r1")
        if (keyPair == null) {
            displayString = "Unable to generate key pair."
        }
        else {
            displayString =
                "Public key: " +
                SecurityHelpers.getEcPublicKeyAsHex(keyPair.public) +
                "\n" +
                "Private key: " + SecurityHelpers.getEcPrivateKeyAsHex(keyPair.private)
        }

        Log.d(TAG, displayString);
        textDisplay.setText(displayString)
    }

}
