package com.log3900.utils.format

import java.nio.ByteBuffer
import java.util.*



object UUIDUtils {
    fun uuidToByteArray(uuid: UUID): ByteArray {
        val bb = ByteBuffer.wrap(ByteArray(16))
        bb.putLong(uuid.mostSignificantBits)
        bb.putLong(uuid.leastSignificantBits)
        return bb.array()
    }

    fun byteArrayToUUID(byteArray: ByteArray): UUID {
        val bb = ByteBuffer.wrap(byteArray)
        val high = bb.long
        val low = bb.long
        return UUID(high, low)
    }
}