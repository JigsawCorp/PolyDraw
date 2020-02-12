package com.log3900.chat.Message

import android.app.Service
import android.content.Intent
import android.os.Binder
import android.os.Handler
import android.os.IBinder
import android.os.Looper
import com.daveanthonythomas.moshipack.MoshiPack
import com.google.gson.JsonObject
import com.log3900.chat.ChatRestService
import com.log3900.socket.Message
import com.log3900.socket.SocketService
import com.log3900.utils.format.moshi.TimeStampAdapter
import com.log3900.utils.format.moshi.UUIDAdapter
import com.squareup.moshi.JsonAdapter
import com.squareup.moshi.Moshi
import com.squareup.moshi.Types
import retrofit2.Call
import java.util.*
import java.util.concurrent.ConcurrentHashMap
import kotlin.collections.ArrayList
import retrofit2.Callback
import retrofit2.Response

class MessageRepository : Service() {
    enum class Event {
        MESSAGE_RECEIVED
    }

    // Service
    private val binder = MessageRepositoryBinder()
    private var socketService: SocketService? = null
    private var subscribers: ConcurrentHashMap<Event, ArrayList<Handler>> = ConcurrentHashMap()

    // Data
    private val cachedMessages: ConcurrentHashMap<UUID, LinkedList<ReceivedMessage>> = ConcurrentHashMap()

    companion object {
        var instance: MessageRepository? = null
    }

    fun getChannelMessages(channelID: String, sessionToken: String, startIndex: Int, endIndex: Int): LinkedList<ReceivedMessage> {
        var messages: LinkedList<ReceivedMessage>? = null
        if (true) {
            if (!cachedMessages.containsKey(UUID.fromString(channelID))) {
                cachedMessages[UUID.fromString(channelID)] = LinkedList<ReceivedMessage>()
            }
            messages = cachedMessages[UUID.fromString(channelID)]
            println("repository messages = " + messages)
        } else {
            val call = ChatRestService.service.getChannelMessages(
                sessionToken,
                "EN",
                channelID,
                startIndex,
                endIndex
            )
            call.enqueue(object : Callback<JsonObject> {
                override fun onResponse(call: Call<JsonObject>, response: Response<JsonObject>) {
                    when (response.code()) {
                        200 -> {
                            val moshi = Moshi.Builder()
                                .add(UUIDAdapter())
                                .build()
                            val adapter: JsonAdapter<LinkedList<ReceivedMessage>> = moshi.adapter(
                                Types.newParameterizedType(
                                    ArrayList::class.java,
                                    ReceivedMessage::class.java
                                )
                            )
                            messages =
                                adapter.fromJson(response.body()!!.getAsJsonArray("Messages").toString())
                        }
                        else -> {

                        }
                    }
                }

                override fun onFailure(call: Call<JsonObject>, t: Throwable) {
                }
            })
        }

        return messages!!
    }

    fun sendMessage(message: SentMessage) {
        socketService?.sendSerializedMessage(com.log3900.socket.Event.MESSAGE_SENT, message)
    }

    fun sendMessage(messageText: String) {
        val message = SentMessage(messageText, UUID.randomUUID())
        sendMessage(message)
    }

    fun subscribe(event: Event, handler: Handler) {
        if (!subscribers.containsKey(event)) {
            subscribers[event] = ArrayList()
        }

        subscribers[event]?.add(handler)
    }

    override fun onBind(intent: Intent?): IBinder? {
        return binder
    }

    override fun onCreate() {
        super.onCreate()
        instance = this
        socketService = SocketService.instance
        cachedMessages[UUID.fromString("00000000-0000-0000-0000-000000000000")] = LinkedList()

        Thread(Runnable {
            Looper.prepare()
            socketService?.subscribeToMessage(com.log3900.socket.Event.MESSAGE_RECEIVED, Handler {
                receiveMessage(it.obj as Message)
                true
            })

            Looper.loop()
        }).start()
    }

    override fun onDestroy() {
        super.onDestroy()
        socketService = null
    }

    private fun receiveMessage(message: Message) {
        val tempMessage = android.os.Message()
        tempMessage.what = Event.MESSAGE_RECEIVED.ordinal
        val moshi = MoshiPack({
            add(TimeStampAdapter())
            add(UUIDAdapter())
        })
        tempMessage.obj = moshi.unpack(message.data) as ReceivedMessage
        addMessageToCache(tempMessage.obj as ReceivedMessage)
        notifySubscribers(Event.MESSAGE_RECEIVED, tempMessage)
    }

    private fun addMessageToCache(message: ReceivedMessage) {
        if (!cachedMessages.containsKey(message.channelID)) {
            cachedMessages[message.channelID] = LinkedList<ReceivedMessage>()
        }
        cachedMessages[message.channelID]?.addLast(message)
        //cachedMessages[UUID.fromString("00000000-0000-0000-0000-000000000000")]?.addLast(message)
        //println("cahced message array = " + cachedMessages[UUID.fromString("00000000-0000-0000-0000-000000000000")])
        //println("added message to cache, messages = " + cachedMessages[UUID.fromString("00000000-0000-0000-0000-000000000000")])
    }

    private fun notifySubscribers(event: Event, message: android.os.Message) {
        if (subscribers.containsKey(event)) {
            val handlers = subscribers[event]
            for (handler in handlers.orEmpty()) {
                val messageCopy = android.os.Message()
                messageCopy.what = message.what
                messageCopy.obj = message.obj
                handler.sendMessage(messageCopy)
            }
        }
    }

    inner class MessageRepositoryBinder : Binder() {
        fun getService(): MessageRepository = this@MessageRepository
    }
}