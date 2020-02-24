package com.log3900.profile

import android.app.Activity
import android.app.Dialog
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.GridLayout
import android.widget.ImageView
import android.widget.TextView
import androidx.core.view.children
import androidx.fragment.app.DialogFragment
import androidx.fragment.app.FragmentActivity
import com.google.android.material.button.MaterialButton
import com.log3900.R
import com.log3900.utils.ui.getAvatarID

class ModifyAvatarDialog : DialogFragment() {

    companion object {
        fun start(activity: FragmentActivity) {
            val fragmentManager = activity.supportFragmentManager
            val ft = fragmentManager.beginTransaction()
            fragmentManager.findFragmentByTag("avatar_dialog")?.let {
                ft.remove(it)
            }
            ft.addToBackStack(null)

            val avatarDialogFragment = ModifyAvatarDialog()
            avatarDialogFragment.show(ft, "avatar_dialog")
        }
    }

    var selectedAvatarID = -1
    lateinit var selectedAvatarView: ImageView
    lateinit var selectedTitle: TextView
    lateinit var confirmButton: MaterialButton

    override fun onStart() {
        super.onStart()
        dialog?.window?.setLayout(
            ViewGroup.LayoutParams.WRAP_CONTENT,
            ViewGroup.LayoutParams.WRAP_CONTENT
        )
    }

    override fun onCreateDialog(savedInstanceState: Bundle?): Dialog {
        val dialog = super.onCreateDialog(savedInstanceState)
        dialog.setCanceledOnTouchOutside(false)
        return dialog
    }

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        val rootView = inflater.inflate(R.layout.dialog_modify_avatar, container, false)
        setUpUi(rootView)
        return rootView
    }

    private fun setUpUi(root: View) {
        val closeButton: MaterialButton = root.findViewById(R.id.close_dialog_button)
        closeButton.setOnClickListener {
            dismiss()
        }
        confirmButton = root.findViewById(R.id.confirm_dialog_button)

        selectedAvatarView = root.findViewById(R.id.selected_avatar)
        selectedTitle = root.findViewById(R.id.selected_title)

        val avatarContainer: GridLayout = root.findViewById(R.id.avatar_container)
        avatarContainer.children
            .sortedBy { it.id }
            .forEachIndexed { index, avatar ->
                avatar.setOnClickListener {
                    changeSelectedAvatar(index)
                }
        }
    }

    private fun changeSelectedAvatar(index: Int) {
        selectedTitle.visibility = View.VISIBLE
        selectedAvatarView.visibility = View.VISIBLE
        confirmButton.isEnabled = true

        selectedAvatarID = index + 1
        selectedAvatarView.setImageResource(getAvatarID(selectedAvatarID))
    }
}