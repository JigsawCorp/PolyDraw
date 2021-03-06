package com.log3900.utils.ui

import com.log3900.R
import com.log3900.user.account.Account

/**
 * Returns the ID of the pictureID from the corresponding index
 * @param index The picture index (from [1..16] included)
 * @return avatar_1 if index out of range
 */
fun getAvatarID(index: Int): Int {
    return when(index) {
        0           -> R.drawable.ic_robot  // Virtual player
        in 1..16    -> R.mipmap::class.java.getField("avatar_${index}").getInt(null)
        else        -> R.mipmap.avatar_1
    }
}

/**
 * Returns the ID of the given [Account]'s pictureID linking to the corresponding mipmap
 * @sample avatarView.setImageResource(getAccountAvatarID(myAccount))
 */
fun getAccountAvatarID(account: Account): Int = getAvatarID(account.pictureID)
