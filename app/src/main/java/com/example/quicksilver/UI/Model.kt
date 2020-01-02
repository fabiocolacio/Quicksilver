package com.example.quicksilver.UI

import android.util.Log
import java.util.*

class Model (var loginFinished : Boolean, var logMessages : ArrayList<String>) : Observable() {

    companion object {
        private val TAG = Model.javaClass.canonicalName
    }

    init {
        loginFinished = false
        logMessages = ArrayList()
    }

    fun updateLoginStatus(newLoginStatus : Boolean) {
        loginFinished = newLoginStatus
        notifyObservers()
    }

    fun updateNewLogMessage(msg : String) {
        Log.d(TAG, "updateNewLogMessage called for message: $msg")
        logMessages.add(msg)
        notifyObservers()
    }

    override fun notifyObservers() {
        setChanged()
        super.notifyObservers()
    }
}