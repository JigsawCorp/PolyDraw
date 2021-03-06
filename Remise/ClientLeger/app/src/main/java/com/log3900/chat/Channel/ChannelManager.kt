package com.log3900.chat.Channel

import com.log3900.MainActivity
import com.log3900.MainApplication
import com.log3900.R
import com.log3900.chat.ChatMessage
import com.log3900.shared.architecture.DialogEventMessage
import com.log3900.shared.architecture.EventType
import com.log3900.shared.architecture.MessageEvent
import com.log3900.user.UserRepository
import com.log3900.user.account.AccountRepository
import org.greenrobot.eventbus.EventBus
import org.greenrobot.eventbus.Subscribe
import org.greenrobot.eventbus.ThreadMode
import java.util.*

class ChannelManager {
    var activeChannel: Channel? = null
    lateinit var availableChannels: ArrayList<Channel>
    lateinit var joinedChannels: ArrayList<Channel>
    private var unreadMessages: HashMap<UUID, Int> = hashMapOf()
    private var unreadMessagesTotal: Int = 0
    var previousChannel: Channel? = null

    constructor() {
    }

    fun init() {
        joinedChannels = ChannelRepository.instance?.getJoinedChannels(AccountRepository.getInstance().getAccount().sessionToken)?.blockingGet()!!
        availableChannels = ChannelRepository.instance?.getAvailableChannels(AccountRepository.getInstance().getAccount().sessionToken)?.blockingGet()!!
        activeChannel = joinedChannels.find {
            it.ID == Channel.GENERAL_CHANNEL_ID
        }!!
        EventBus.getDefault().register(this)
    }

    fun changeSubscriptionStatus(channel: Channel) {
        if (channel.ID == Channel.GENERAL_CHANNEL_ID) {
            return
        }

        if (availableChannels.contains(channel)) {
            ChannelRepository.instance?.subscribeToChannel(channel)
        } else if (joinedChannels.contains(channel)){
            ChannelRepository.instance?.unsubscribeFromChannel(channel)
        } else {
            // TODO: Handle this incoherent state
        }

    }

    fun createChannel(channelName: String): Boolean {
        ChannelRepository.instance?.createChannel(channelName)

        return true
    }

    fun deleteChannel(channel: Channel) {
        ChannelRepository.instance?.deleteChannel(channel)
    }

    @Subscribe(threadMode = ThreadMode.BACKGROUND)
    fun onMessageEvent(event: MessageEvent) {
        when(event.type) {
            EventType.CHANNEL_CREATED -> {
                onChannelCreated(event.data as Channel)
            }
            EventType.CHANNEL_DELETED -> {
                onChannelDeleted(event.data as ChannelDeletedMessage)
            }
            EventType.RECEIVED_MESSAGE -> {
                onMessageReceived(event.data as ChatMessage)
            }
            EventType.SUBSCRIBED_TO_CHANNEL -> {
                onSubscribedToChannel(event.data as Channel)
            }
            EventType.UNSUBSCRIBED_FROM_CHANNEL -> {
                onUnsubscribedFromChannel(event.data as Channel)
            }
        }
    }

    fun onChannelCreated(channel: Channel) {
        if (channel.users.get(0).ID == AccountRepository.getInstance().getAccount().ID) {
            changeSubscriptionStatus(channel)
            changeActiveChannel(channel)
        }
    }

    fun onChannelDeleted(channelDeleted: ChannelDeletedMessage) {
        if (activeChannel?.ID == channelDeleted.channelID) {
            val newActiveChannel = joinedChannels.find {
                it.ID == Channel.GENERAL_CHANNEL_ID
            }!!
            previousChannel = newActiveChannel
            changeActiveChannel(newActiveChannel)
            EventBus.getDefault().post(MessageEvent(EventType.SHOW_ERROR_MESSAGE, DialogEventMessage(
                MainApplication.instance.getContext().getString(R.string.warning),
                MainApplication.instance.getContext().getString(R.string.current_channel_got_deleted, channelDeleted.username, activeChannel?.name),
                null,
                null
            )))
        } else if (activeChannel == null && previousChannel?.ID == channelDeleted.channelID) {
            previousChannel = joinedChannels.find {
                it.ID == Channel.GENERAL_CHANNEL_ID
            }!!
        }

        if (unreadMessages.containsKey(channelDeleted.channelID)) {
            unreadMessagesTotal -= unreadMessages.get(channelDeleted.channelID)!!
            unreadMessages.remove(channelDeleted.channelID)
            EventBus.getDefault().post(MessageEvent(EventType.UNREAD_MESSAGES_CHANGED, unreadMessagesTotal))
        }
    }

    fun changeActiveChannel(channel: Channel?) {
        activeChannel = channel
        if (channel != null && unreadMessages.containsKey(channel.ID)) {
            unreadMessagesTotal -= unreadMessages.get(channel.ID)!!
            unreadMessages[channel.ID] = 0

            EventBus.getDefault().post(MessageEvent(EventType.UNREAD_MESSAGES_CHANGED, unreadMessagesTotal))
        }
        EventBus.getDefault().post(MessageEvent(EventType.ACTIVE_CHANNEL_CHANGED, activeChannel))
    }

    fun getUnreadMessage(): HashMap<UUID, Int> {
        return unreadMessages
    }

    fun getDefaultChannel(): Channel {
        return joinedChannels.find {
            it.ID == Channel.GENERAL_CHANNEL_ID
        }!!
    }

    fun delete() {
        EventBus.getDefault().unregister(this)
    }

    private fun onMessageReceived(message: ChatMessage) {
        if (message.channelID != activeChannel?.ID) {
            if (joinedChannels.find { it.ID == message.channelID } == null) {
                return
            }
            if (!unreadMessages.containsKey(message.channelID)) {
                unreadMessages.put(message.channelID, 0)
            }

            unreadMessages.put(message.channelID, unreadMessages[message.channelID]!! + 1)
            unreadMessagesTotal += 1
            EventBus.getDefault().post(MessageEvent(EventType.UNREAD_MESSAGES_CHANGED, unreadMessagesTotal))
        } else {
            EventBus.getDefault().post(MessageEvent(EventType.ACTIVE_CHANNEL_MESSAGE_RECEIVED, message))
        }
    }

    private fun onSubscribedToChannel(channel: Channel) {
        if (channel.isGame) {
            if (activeChannel != null) {
                changeActiveChannel(channel)
            } else {
                previousChannel = channel
            }
        }
    }

    private fun onUnsubscribedFromChannel(channel: Channel) {
        if (activeChannel == channel) {
            val newActiveChannel = joinedChannels.find {
                it.ID == Channel.GENERAL_CHANNEL_ID
            }!!
            changeActiveChannel(newActiveChannel)
            previousChannel = newActiveChannel
        } else if (activeChannel == null && previousChannel?.ID == channel.ID) {
            previousChannel = joinedChannels.find {
                it.ID == Channel.GENERAL_CHANNEL_ID
            }!!
        }

        if (unreadMessages.containsKey(channel.ID)) {
            unreadMessagesTotal -= unreadMessages.get(channel.ID)!!
            unreadMessages[channel.ID] = 0

            EventBus.getDefault().post(MessageEvent(EventType.UNREAD_MESSAGES_CHANGED, unreadMessagesTotal))
        }
    }


}