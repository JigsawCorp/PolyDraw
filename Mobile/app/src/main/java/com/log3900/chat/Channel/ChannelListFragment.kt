package com.log3900.chat.Channel

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.log3900.R
import io.github.luizgrp.sectionedrecyclerviewadapter.SectionedRecyclerViewAdapter
import java.util.*

class ChannelListFragment : Fragment(), ChannelListView {
    // Services
    private lateinit var channelListPresenter: ChannelListPresenter
    // UI elements
    private lateinit var channelsRecyclerView: RecyclerView
    private lateinit var layoutManager: LinearLayoutManager

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        val rootView: View = inflater.inflate(R.layout.fragment_channel_list, container, false)

        var channelGroups = arrayListOf<ChannelGroup>()
        var channels1 = arrayListOf<Channel>()
        channels1.add(Channel(UUID.randomUUID(), "channel1", arrayOf()))
        channels1.add(Channel(UUID.randomUUID(), "channel2", arrayOf()))
        channels1.add(Channel(UUID.randomUUID(), "channel3", arrayOf()))
        var channels2 = arrayListOf<Channel>()
        channels2.add(Channel(UUID.randomUUID(), "channel4", arrayOf()))
        channels2.add(Channel(UUID.randomUUID(), "channel5", arrayOf()))
        channels2.add(Channel(UUID.randomUUID(), "channel6", arrayOf()))
        channelGroups.add(ChannelGroup("ChannelGroup1", channels1))
        channelGroups.add(ChannelGroup("ChannelGroup2", channels2))

        channelsRecyclerView = rootView.findViewById(R.id.fragment_channel_list_recycler_view_channels)
        layoutManager = LinearLayoutManager(this.context)
        channelsRecyclerView.layoutManager = layoutManager
        val channelAdapter = SectionedRecyclerViewAdapter()
        channelAdapter.addSection(channelGroups.get(0).name, ChannelSection(channelGroups.get(0), channels1, object: ChannelSection.ClickListener {
            override fun onHeaderRootViewClick(group: ChannelGroup) {
                val sectionAdapter = channelAdapter.getAdapterForSection(group.name)
                val section = channelAdapter.getSection(group.name) as ChannelSection
                val count = section.contentItemsTotal

                section.expanded = !section.expanded
                sectionAdapter.notifyHeaderChanged()

                if (!section.expanded) {
                    sectionAdapter.notifyItemRangeRemoved(0, count)
                } else {
                    sectionAdapter.notifyAllItemsInserted()
                }
            }
        }))
        channelAdapter.addSection(channelGroups.get(1).name, ChannelSection(channelGroups.get(1),channels2, object: ChannelSection.ClickListener {
            override fun onHeaderRootViewClick(group: ChannelGroup) {
            }
        }))

        channelListPresenter = ChannelListPresenter(this)
        channelsRecyclerView.adapter = channelAdapter
        return rootView
    }
}