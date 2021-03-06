package com.log3900.login

import android.content.Context
import android.content.SharedPreferences
import android.util.Log
import com.log3900.MainApplication
import com.log3900.R
import com.log3900.shared.architecture.EventType
import com.log3900.shared.architecture.MessageEvent
import org.greenrobot.eventbus.EventBus
import org.greenrobot.eventbus.Subscribe
import org.greenrobot.eventbus.ThreadMode

object UserPrefsManager {
    private val sharedPrefs: SharedPreferences = with (MainApplication.instance) {
        getSharedPreferences(resources.getString(R.string.preference_file_key), Context.MODE_PRIVATE)
    }
    private val bearerPrefKey = MainApplication.instance.resources.getString(R.string.preference_file_bearer_token_key)
    private val usernamePrefKey = MainApplication.instance.resources.getString(R.string.preference_file_username_key)

    init {
        EventBus.getDefault().register(this)
    }

    /// Bearer
    fun getBearer(): String? = sharedPrefs.getString(bearerPrefKey, null)
    fun storeBearer(token: String) = store(bearerPrefKey, token)
    fun resetBearer() = reset(bearerPrefKey)

    /// Username
    fun getUsername(): String? = sharedPrefs.getString(usernamePrefKey, null)
    fun storeUsername(username: String) = store(usernamePrefKey, username)
    fun resetUsername() = reset(usernamePrefKey)

    fun resetAll() {
        resetBearer()
        resetUsername()
    }

    private fun store(key: String, value: String) {
        with(sharedPrefs.edit()) {
            putString(key, value)
            commit()
        }
    }

    private fun reset(key: String) {
        with(sharedPrefs.edit()) {
            putString(key, null)
            commit()
        }
    }

    @Subscribe(threadMode = ThreadMode.MAIN)
    fun onEvent(event: MessageEvent) {
        when (event.type) {
            EventType.LOGOUT -> resetAll()
            else -> return
        }
    }
}