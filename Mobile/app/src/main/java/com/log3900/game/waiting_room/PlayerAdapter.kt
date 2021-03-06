package com.log3900.game.waiting_room

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.log3900.R
import com.log3900.game.group.Group
import com.log3900.game.group.MatchMode
import com.log3900.game.group.Player
import com.log3900.user.account.AccountRepository
import java.util.*
import kotlin.collections.ArrayList

class PlayerAdapter: RecyclerView.Adapter<PlayerViewHolder> {
    private lateinit var players: ArrayList<Player>
    private lateinit var playersCopy: ArrayList<Player?>
    private lateinit var group: Group
    private var listener: PlayerViewHolder.Listener

    constructor(listener: PlayerViewHolder.Listener) {
        this.listener = listener
        playersCopy = arrayListOf()
    }

    fun setGroup(group: Group) {
        this.group = group

        for (i in 1..group.playersMax) {
            playersCopy.add(null)
        }
    }

    fun setPlayers(players: ArrayList<Player>) {
        this.players = players

        refreshLists()
    }

    fun playerAdded(playerID: UUID) {
        refreshLists()
    }

    fun playerRemoved(playerID: UUID) {
        refreshLists()
    }

    private fun refreshLists() {
        playersCopy.clear()
        if (group.gameType == MatchMode.FFA) {
            playersCopy = players.clone() as ArrayList<Player?>
        } else {
            players.forEach {
                if (!it.isCPU) {
                    playersCopy.add(it)
                }
            }
        }

        for (i in playersCopy.size until group.playersMax) {
            playersCopy.add(null)
        }
        notifyDataSetChanged()
    }

    override fun getItemCount(): Int {
        return playersCopy.size
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): PlayerViewHolder {
        val view = LayoutInflater.from(parent.context).inflate(R.layout.list_item_match_waiting_room_player, parent, false)

        return PlayerViewHolder(view, listener)
    }

    override fun onBindViewHolder(holder: PlayerViewHolder, position: Int) {
        holder.bind(
            playersCopy[position],
            playersCopy[position] == null || (playersCopy[position]!!.isCPU && group.gameType != MatchMode.FFA),
            group.ownerID == playersCopy[position]?.ID,
            AccountRepository.getInstance().getAccount().ID == group.ownerID,
            group.gameType == MatchMode.FFA)
    }
}