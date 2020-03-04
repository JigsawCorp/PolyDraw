package com.log3900.settings.language

import android.content.Context
import android.content.res.Resources
import android.util.Log
import com.log3900.MainApplication
import com.log3900.R
import com.log3900.user.account.AccountRepository
import io.reactivex.Completable
import java.util.*
import kotlin.collections.ArrayList

class LanguageManager {
    enum class LANGUAGE {
        SYSTEM,
        ENGLISH,
        FRENCH
    }

    companion object {
        private val languages: ArrayList<Language> = arrayListOf(
            Language(LANGUAGE.SYSTEM.ordinal, "", R.string.language_system_title),
            Language(LANGUAGE.ENGLISH.ordinal, "en", R.string.language_english_title),
            Language(LANGUAGE.FRENCH.ordinal, "fr", R.string.language_french_title)
        )
        fun getAvailableLanguages(): ArrayList<Language> {
            return languages
        }
        fun getCurrentLanguage(): Language {
            return languages[AccountRepository.getInstance().getAccount().languageID]
        }

        fun changeLanguage(language: Language): Completable {
            val currentAccount = AccountRepository.getInstance().getAccount()
            currentAccount.languageID = language.index
            return AccountRepository.getInstance().updateAccount(currentAccount)
        }

        fun applySavedLanguage(context: Context) {
            val res = context.resources
            val displayMetric = res.displayMetrics
            val configuration = res.configuration
            if (getCurrentLanguage().index != LANGUAGE.SYSTEM.ordinal) {
                configuration.setLocale(Locale(getCurrentLanguage().languageCode))
            } else {
                configuration.setLocale(Resources.getSystem().configuration.locales[0])
            }
            res.updateConfiguration(configuration, displayMetric)
        }
    }
}