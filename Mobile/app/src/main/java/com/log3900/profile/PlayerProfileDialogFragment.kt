package com.log3900.profile

import android.content.Context
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.DialogFragment
import androidx.fragment.app.FragmentActivity
import com.google.android.material.button.MaterialButton
import com.log3900.R
import com.log3900.utils.ui.getAvatarID
import kotlinx.android.synthetic.main.dialog_fragment_player_profile.*
import java.util.*

class PlayerProfileDialogFragment(userID: UUID) : DialogFragment() {
    private val presenter = PlayerProfilePresenter(this, userID)

    companion object {
        fun show(context: Context, userID: UUID) {
            val fragmentManager = (context as FragmentActivity).supportFragmentManager//activity?.supportFragmentManager!!
            val ft = fragmentManager.beginTransaction()
            fragmentManager.findFragmentByTag("dialog")?.let {
                ft.remove(it)
            }
            ft.addToBackStack(null)

            PlayerProfileDialogFragment(userID).show(ft, "dialog")
        }
    }

    override fun onStart() {
        super.onStart()
        dialog?.window?.setLayout(
            ViewGroup.LayoutParams.WRAP_CONTENT,
            ViewGroup.LayoutParams.WRAP_CONTENT
        )
    }

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        val rootView = inflater.inflate(R.layout.dialog_fragment_player_profile, container, false)
        setUpUi(rootView)
        return rootView
    }

    private fun setUpUi(root: View) {
        presenter.fetchUserInfo()
        root.findViewById<MaterialButton>(R.id.close_dialog_button).setOnClickListener {
            dismiss()
        }
    }

    fun setUsername(username: String) {
        username_text_view.text = username
    }

    fun setAvatar(avatarID: Int) {
        avatar_image_view.setImageResource(getAvatarID(avatarID))
    }
}