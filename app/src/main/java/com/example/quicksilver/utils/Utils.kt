package com.example.quicksilver.utils

class Utils {
    companion object {
        fun byteArrayToHexString(arr : ByteArray) : String {
            var ret = StringBuilder()
            for (byte in arr) {
                ret.append(String.format("%02X", byte))
            }
            return ret.toString()
        }
    }
}