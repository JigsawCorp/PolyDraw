package com.log3900.game.match

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.log3900.R
import com.log3900.game.group.Group
import com.log3900.game.group.Player
import java.util.*
import kotlin.collections.ArrayList

class PlayerAdapter: RecyclerView.Adapter<PlayerViewHolder>() {
    private var players: ArrayList<Player> = arrayListOf()
    private lateinit var group: Group


    fun setGroup(group: Group) {
        this.group = group
    }

    fun setPlayers(players: ArrayList<Player>) {
        this.players = players
    }

    override fun getItemCount(): Int {
        return players.size
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): PlayerViewHolder {
        val view = LayoutInflater.from(parent.context).inflate(R.layout.list_item_active_match_player, parent, false)

        return PlayerViewHolder(view)
    }

    override fun onBindViewHolder(holder: PlayerViewHolder, position: Int) {
        holder.bind(players[position], position + 1, 0)
    }
}