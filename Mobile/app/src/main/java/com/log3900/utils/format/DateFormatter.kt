package com.log3900.utils.format

import android.text.format.DateUtils
import java.text.SimpleDateFormat
import java.util.*

class DateFormatter {
    companion object {

        fun formatDate(date: Date): String {
            var dateFormat: SimpleDateFormat? = null
            if (DateUtils.isToday(date.time)) {
                dateFormat = SimpleDateFormat("HH:mm:ss")
            }
            else {
                dateFormat = SimpleDateFormat("dd/MM/yyyy HH:mm:ss")
            }
            val dateString = dateFormat.format(date)
            return dateString
        }

        fun formatDateToTime(date: Date): String {
            var dateFormat: SimpleDateFormat = SimpleDateFormat("mm:ss")
            return dateFormat.format(date)
        }

        fun formatFullDate(date: Date): String = SimpleDateFormat("dd/MM/yyyy HH:mm:ss").format(date)

        fun formatFullTime(date: Date): String {
            val formatter = SimpleDateFormat("HH:mm:ss")
            formatter.timeZone = TimeZone.getTimeZone("UTC")
            return formatter.format(date)
        }
    }
}