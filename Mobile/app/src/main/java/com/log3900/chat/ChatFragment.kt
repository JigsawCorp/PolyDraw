package com.log3900.chat

import android.media.MediaPlayer
import android.os.Bundle
import android.text.InputFilter
import android.text.InputType
import android.util.Log
import android.util.TypedValue
import android.view.Gravity
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.view.inputmethod.EditorInfo
import android.widget.Button
import android.widget.ImageView
import androidx.appcompat.widget.Toolbar
import androidx.core.view.GravityCompat
import androidx.drawerlayout.widget.DrawerLayout
import androidx.fragment.app.DialogFragment
import androidx.fragment.app.Fragment
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.google.android.material.textfield.TextInputEditText
import com.log3900.MainApplication
import com.log3900.R
import com.log3900.chat.ui.MessageAdapter
import com.log3900.settings.sound.SoundManager
import com.log3900.user.account.AccountRepository
import java.util.*

class ChatFragment : Fragment(), ChatView {
    // Services
    private var chatPresenter: ChatPresenter? = null
    // UI elements
    private lateinit var messagesRecyclerView: RecyclerView
    private lateinit var messagesViewAdapter: MessageAdapter
    private lateinit var viewManager: LinearLayoutManager
    private lateinit var sendMessageButton: ImageView
    private lateinit var toolbar: Toolbar
    private lateinit var drawer: DrawerLayout
    private lateinit var messageEditText: TextInputEditText

    private var autoScroll = true

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        val rootView: View = inflater.inflate(R.layout.fragment_chat, container, false)

        setupUiElements(rootView)

        chatPresenter = ChatPresenter(this)

        return rootView
    }

    /**
     * Used to assign all variables holding UI elements for this fragment and setup listeners.
     *
     * @param rootView the root view of this fragment
     */
    private fun setupUiElements(rootView: View) {
        viewManager = LinearLayoutManager(activity)

        setupMessagesRecyclerView(rootView)

        sendMessageButton = rootView.findViewById(R.id.fragment_chat_send_message_button)
        sendMessageButton.setOnClickListener { v -> onSendMessageButtonClick(v) }

        messageEditText = rootView.findViewById(R.id.fragment_chat_new_message_input)
        messageEditText.filters += InputFilter.LengthFilter(1000)
        messageEditText.inputType = InputType.TYPE_TEXT_FLAG_NO_SUGGESTIONS or InputType.TYPE_TEXT_VARIATION_VISIBLE_PASSWORD

        setupToolbar(rootView)

        drawer = rootView.findViewById(R.id.fragment_chat_drawer_layout)

        rootView.findViewById<TextInputEditText>(R.id.fragment_chat_new_message_input).apply {
            setOnClickListener{ v -> onMessageTextInputClick(v)}
            setOnEditorActionListener { textView, actionID, keyEvent ->
                return@setOnEditorActionListener when (actionID) {
                    EditorInfo.IME_ACTION_SEND -> {
                        onSendMessageButtonClick(textView)
                        true
                    }
                    else -> false
                }
            }
        }

        rootView.viewTreeObserver.addOnGlobalLayoutListener {
            val heightDiff = rootView.rootView.height - rootView.height
            val pixels = TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, 200f, context?.resources?.displayMetrics)
            if (heightDiff > pixels) {
                chatPresenter?.onKeyboardChange(true)
            }
            else {
                chatPresenter?.onKeyboardChange(false)
            }
        }

        messagesRecyclerView.addOnScrollListener(object : RecyclerView.OnScrollListener() {
            override fun onScrolled(recyclerView: RecyclerView, dx: Int, dy: Int) {
                super.onScrolled(recyclerView, dx, dy)

                if (dy < 0) {
                    autoScroll = false
                } else {
                    if (!recyclerView.canScrollVertically(1)) {
                        autoScroll = true
                    }
                }
            }
        })

    }

    private fun setupToolbar(rootView: View) {
        toolbar = rootView.findViewById(R.id.fragment_chat_top_layout)
        toolbar.setNavigationIcon(R.drawable.ic_hamburger_menu)

        toolbar.setNavigationOnClickListener {chatPresenter?.handleNavigationDrawerClick()}
    }

    private fun setupMessagesRecyclerView(rootView: View) {
        messagesRecyclerView = rootView.findViewById(R.id.fragment_chat_message_recycler_view)
        messagesViewAdapter = MessageAdapter(LinkedList(), AccountRepository.getInstance().getAccount().username)
        messagesRecyclerView.apply {
            setHasFixedSize(true)
            layoutManager = viewManager
            adapter = messagesViewAdapter
        }

        val scrollListener = object: RecyclerView.OnScrollListener() {
            override fun onScrolled(recyclerView: RecyclerView, dx: Int, dy: Int) {
                super.onScrolled(recyclerView, dx, dy)
                chatPresenter?.scrolledToPositions(viewManager.findFirstVisibleItemPosition(), viewManager.findLastVisibleItemPosition(), dy)
            }
        }

        messagesRecyclerView.addOnScrollListener(scrollListener)
    }

    private fun onSendMessageButtonClick(v: View) {
        val messageInput: TextInputEditText = v.rootView.findViewById(R.id.fragment_chat_new_message_input)
        val messageText = messageInput.text.toString()
        messageInput.text?.clear()
        if (messageText != "" && messageText.trim().length > 0)
        {
            chatPresenter?.sendMessage(messageText.trim())
        }
    }

    private fun onMessageTextInputClick(v: View) {
        messagesViewAdapter.smoothScrollToBottom()
    }

    override fun onResume() {
        super.onResume()
        chatPresenter?.resume()
    }

    override fun openNavigationDrawer() {
        drawer.openDrawer(Gravity.LEFT)
    }

    override fun closeNavigationDrawer() {
        drawer.closeDrawer(Gravity.LEFT)
    }

    override fun isNavigationDrawerOpened(): Boolean {
        return drawer.isDrawerOpen(GravityCompat.START)
    }

    override fun notifyNewMessage() {
        if (autoScroll) {
            scrollMessage(true)
        }
        //messagesViewAdapter.messageInserted()
    }

    override fun setChatMessages(messages: LinkedList<ChatMessage>) {
        messagesViewAdapter.setMessage(messages)
        messagesViewAdapter.notifyDataSetChanged()
    }

    override fun setCurrentChannnelName(name: String) {
        toolbar.setTitle(name)
    }

    override fun scrollMessage(smooth: Boolean) {
        if (smooth) {
            messagesViewAdapter.smoothScrollToBottom()
        } else {
            messagesViewAdapter.scrollToBottom()
        }
    }

    override fun showProgressDialog(dialog: DialogFragment) {
        dialog.show(activity?.supportFragmentManager!!, "progressDialog")
    }

    override fun hideProgressDialog(dialog: DialogFragment) {
        dialog.dismiss()
    }

    override fun notifyMessagesPrepended(count: Int) {
        messagesViewAdapter.prependMessages(count)
    }

    override fun notifyMessagesUpdated(userChanged: UUID) {
        messagesViewAdapter.notifyDataSetChanged()
    }

    override fun onDestroy() {
        super.onDestroy()
        chatPresenter?.destroy()
        chatPresenter = null
    }

}