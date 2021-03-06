package com.log3900.socket

import android.os.Handler
import android.os.Looper
import android.util.Log
import java.io.*
import java.net.InetSocketAddress
import java.net.Socket
import java.net.SocketException
import java.util.*
import java.util.concurrent.CountDownLatch
import java.util.concurrent.atomic.AtomicBoolean
import java.util.concurrent.atomic.AtomicReference
import kotlin.concurrent.timerTask

enum class Request {
    SEND_MESSAGE,
    START_READING,
    STOP_READING,
    CONNECT,
    DISCONNECT,
    SET_MESSAGE_LISTENER
}

enum class State {
    CONNECTED,
    DISCONNECTED,
    DISCONNECTING,
    ERROR
}

object SocketHandler {
    private lateinit var socket: Socket
    private lateinit var inputStream: DataInputStream
    private lateinit var outputStream: DataOutputStream
    private var requestHandler: Handler? = null
    private var messageReadListener: Handler? = null
    private var connectionErrorListener: Handler? = null
    private var disconnectionErrorListener: Handler? = null
    private var readMessages = AtomicBoolean(false)
    private var socketHealthcheckTimer: Timer = Timer()
    public var state: AtomicReference<State> = AtomicReference(State.DISCONNECTED)

    fun connect() {
        socket = Socket()
        socket.connect(InetSocketAddress("log3900-203.canadacentral.cloudapp.azure.com", 5001), 10000)
        inputStream = DataInputStream(socket.getInputStream())
        outputStream = DataOutputStream(socket.getOutputStream())
        state.set(State.CONNECTED)
    }

    fun createRequestHandler() {
        val lock = CountDownLatch(1)
        Thread(Runnable {
            Looper.prepare()
            requestHandler = Handler {
                handleRequest(it)
                true
            }
            lock.countDown()
            Looper.loop()
        }).start()
        lock.await()
    }

    fun setMessageReadListener(handler: Handler?) {
        messageReadListener = handler
    }

    fun setConnectionErrorListener(handler: Handler?) {
        connectionErrorListener = handler
    }

    fun setDisconnectionListener(handler: Handler?) {
        disconnectionErrorListener = handler
    }


    fun sendRequest(message: android.os.Message) {
        requestHandler?.sendMessage(message)
    }

    private fun onWriteMessage(message: Message) {
        try {
            outputStream.writeByte(message.type.eventType.toInt())
            outputStream.writeShort(message.data.size)
            outputStream.write(message.data)
        } catch (e: Exception) {
            handlerError()
        }
    }

    fun onDisconnect() {
        if (state.get() == State.CONNECTED) {
            state.set(State.DISCONNECTING)
            socketHealthcheckTimer.cancel()
            socket.close()
            outputStream.close()
            inputStream.close()
        }
    }

    private fun handleRequest(message: android.os.Message) {
        when (message.what) {
            Request.SEND_MESSAGE.ordinal -> {
                if (message.obj is Message) {
                    onWriteMessage(message.obj as Message)
                }
            }
            Request.START_READING.ordinal -> {
                onReadMessage()
            }
            Request.STOP_READING.ordinal -> {
                onStopReadingMessages()
            }
            Request.DISCONNECT.ordinal -> {
                onDisconnect()
            }
            Request.SET_MESSAGE_LISTENER.ordinal -> {
                if (message.obj is Handler) {
                    messageReadListener = message.obj as Handler
                }
            }
        }
    }

    fun onStopReadingMessages() {
        readMessages.compareAndSet(true, false)
    }

    fun onReadMessage() {
        if (!readMessages.get()) {
            readMessages.compareAndSet(false, true)
            Thread(Runnable {
                while (readMessages.get()) {
                    readMessage()
                }
            }).start()
        }
    }

    fun readMessage() {
        try {
            val typeByte = inputStream.readUnsignedByte()

            val length = inputStream.readUnsignedShort()

            var values = ByteArray(length.toInt())
            var totalReadBytes = 0
            while (totalReadBytes < length) {
                val amountRead = inputStream.read(values, totalReadBytes, length - totalReadBytes)
                totalReadBytes += amountRead
            }

            val type = Event.values().find { it.eventType == typeByte.toInt() }
                ?: return

            val message = Message(type, values)
//            Log.d("DRAW", message.toString())

            if (message.type == Event.HEALTH_CHECK_SERVER) {
                Log.d("Healthcheck", "Received server healthcheck")
                onWriteMessage(Message(Event.HEALTH_CHECK_CLIENT, byteArrayOf()))
                Log.d("Healthcheck", "Sent healthcheck response")
                socketHealthcheckTimer.cancel()
                socketHealthcheckTimer = Timer()
                socketHealthcheckTimer.schedule( timerTask {
                    Log.d("Healthcheck", "Timer expired")
                    handlerError()
                }, 6000)
            }
            else if (messageReadListener != null) {
                Log.d("POTATO", "New socket message of type ${message.type}")
                val msg = android.os.Message()
                msg.obj = message
                messageReadListener?.sendMessage(msg)
            }

        } catch (e: SocketException){
            Log.d("POTATO", "Caught socket exception when reading = $e")
            val sw = StringWriter()
            val pw = PrintWriter(sw)
            e.printStackTrace(pw)
            handlerError()
        } catch (e: EOFException) {
            Log.d("POTATO", "Caught socket exception when reading = $e")
            handlerError()
        }
    }

    private fun handlerError() {
        Log.d("Healthcheck", "handlerError()")
        if (state.get() == State.DISCONNECTING) {
            Log.d("Healthcheck", "State.Disconnecting")
            state.set(State.DISCONNECTED)
            readMessages.set(false)
            disconnectionErrorListener?.sendEmptyMessage(SocketEvent.DISCONNECTED.ordinal)
        }
        else if (state.get() == State.CONNECTED) {
            Log.d("Healthcheck", "State.conntected")
            state.set(State.ERROR)
            socketHealthcheckTimer.cancel()
            readMessages.set(false)
            socket.close()
            val message = android.os.Message()
            message.what = SocketEvent.CONNECTION_ERROR.ordinal
            connectionErrorListener?.sendMessage(message)
        }
    }
}