package com.log3900.chat.Channel

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import com.log3900.R

class ChannelListFragment : Fragment(), ChannelListView {
    // Services
    private lateinit var channelListPresenter: ChannelListPresenter
    // UI elements


    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        val rootView: View = inflater.inflate(R.layout.fragment_channel_list, container, false)

        channelListPresenter = ChannelListPresenter(this)

        return rootView
    }
}