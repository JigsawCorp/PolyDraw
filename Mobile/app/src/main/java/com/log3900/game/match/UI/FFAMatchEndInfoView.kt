package com.log3900.game.match.UI

import android.content.Context
import android.util.AttributeSet
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.constraintlayout.widget.ConstraintLayout
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.log3900.MainApplication
import com.log3900.R

class FFAMatchEndInfoView(context: Context, attributeSet: AttributeSet) : ConstraintLayout(context, attributeSet) {
    private var layout: ConstraintLayout
    private var winnerTextView: TextView
    private var playersRecyclerView: RecyclerView

    private lateinit var playersAdapter: PlayerAdapter

    init {
        layout = View.inflate(context, R.layout.view_ffa_match_end_info, this) as ConstraintLayout
        winnerTextView = layout.findViewById(R.id.view_ffa_match_end_info_text_view_winner)
        playersRecyclerView = layout.findViewById(R.id.view_ffa_match_end_info_recycler_view_players)

        setupRecyclerView()
    }

    fun setupRecyclerView() {
        playersAdapter = PlayerAdapter()
        playersRecyclerView.apply {
            setHasFixedSize(true)
            layoutManager = LinearLayoutManager(context)
            adapter = playersAdapter
        }
    }

    fun setWinner(winnerName: String, isUser: Boolean) {
        if (isUser) {
            winnerTextView.text = MainApplication.instance.getContext().getString(R.string.ffa_match_end_inf_winner_is_you)
        } else {
            winnerTextView.text = MainApplication.instance.getContext().getString(R.string.ffa_match_end_info_winner_is, winnerName)
        }
    }

    fun setPlayers(players: ArrayList<Pair<String, Int>>) {
        playersAdapter.setPlayers(players)
    }

    class PlayerAdapter : RecyclerView.Adapter<PlayerViewHolder>() {
        private var players: ArrayList<Pair<String, Int>> = arrayListOf()

        fun setPlayers(players: ArrayList<Pair<String, Int>>) {
            this.players = players
            notifyDataSetChanged()
        }

        override fun getItemCount(): Int {
            return players.size
        }

        override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): PlayerViewHolder {
            val view = LayoutInflater.from(parent.context).inflate(R.layout.list_item_ffa_match_end_info_player, parent, false)

            return PlayerViewHolder(view)
        }

        override fun onBindViewHolder(holder: PlayerViewHolder, position: Int) {
            holder.bind(players[position], position + 1)
        }
    }


    class PlayerViewHolder : RecyclerView.ViewHolder {
        private var rootView: View
        private var positionTextView: TextView
        private var nameTextView: TextView
        private var scoreTextView: TextView

        private lateinit var player: Pair<String, Int>

        constructor(itemView: View) : super(itemView) {
            rootView = itemView
            positionTextView = itemView.findViewById(R.id.list_item_ffa_match_end_info_text_view_position)
            nameTextView = itemView.findViewById(R.id.list_item_ffa_match_end_info_text_view_name)
            scoreTextView = itemView.findViewById(R.id.list_item_ffa_match_end_info_text_view_points)
        }

        fun bind(player: Pair<String, Int>, position: Int) {
            this.player = player
            positionTextView.text = "#$position"
            nameTextView.text = player.first
            scoreTextView.text = player.second.toString()
        }
    }
}