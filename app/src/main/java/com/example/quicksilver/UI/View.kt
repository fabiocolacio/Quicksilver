package com.example.quicksilver.UI

import android.graphics.drawable.AnimatedVectorDrawable
import android.util.Log
import android.view.View.INVISIBLE
import android.view.View.VISIBLE
import android.widget.ImageView
import android.widget.TextView
import androidx.vectordrawable.graphics.drawable.AnimatedVectorDrawableCompat
import com.example.quicksilver.MainActivity
import com.example.quicksilver.R
import java.text.SimpleDateFormat
import java.util.*

class View (ctx : MainActivity, var model : Model) : Observer {
    companion object {
        val TAG = View.javaClass.canonicalName
    }

    class OldModelInfo(var loggedIndex : Int)

    private var oldModelInfo : OldModelInfo
    private var textDisplay : TextView
    private var imageCheck : ImageView

    init {
        oldModelInfo = OldModelInfo(-1)
        textDisplay = ctx.findViewById(R.id.text_display)
        imageCheck = ctx.findViewById(R.id.image_check)
        imageCheck.visibility = INVISIBLE
    }

    override fun update(o: Observable?, arg: Any?) {
        updateUI()
    }

    private fun updateUI() {
        if (model.loginFinished) {
            animateImageCheck()
        }
        for (i in oldModelInfo.loggedIndex+1 until model.logMessages.size) {
            Log.d(TAG,"Displaying model log index $i")
            logMessage(model.logMessages[i], textDisplay)
            oldModelInfo.loggedIndex++
        }
    }

    private fun logMessage(msg : String, textDisplay : TextView) {
        val s = SimpleDateFormat("MM:dd:yyyy:hh:mm:ss")
        val format = s.format(Date())
        val display = "--- $format ---\n$msg\n---\n"
        textDisplay.append(display)
    }

    private fun animateImageCheck() {
        imageCheck.visibility = VISIBLE
        if (imageCheck.drawable is AnimatedVectorDrawable) {
            val avd = imageCheck.drawable as AnimatedVectorDrawable
            avd.start()
        }
        else if (imageCheck is AnimatedVectorDrawableCompat) {
            val avd = imageCheck.drawable as AnimatedVectorDrawableCompat
            avd.start()
        }
    }
}