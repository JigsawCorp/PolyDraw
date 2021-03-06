package com.log3900.game.match

import android.animation.ObjectAnimator
import android.animation.PropertyValuesHolder
import android.graphics.Color
import android.util.Log
import android.view.View
import android.widget.ImageView
import android.widget.TextView
import androidx.core.content.res.ResourcesCompat
import androidx.recyclerview.widget.RecyclerView
import com.google.android.material.chip.Chip
import com.log3900.MainApplication
import com.log3900.R
import com.log3900.game.group.Player
import com.log3900.game.waiting_room.PlayerViewHolder
import com.log3900.profile.PlayerProfileDialogFragment
import com.log3900.user.UserRepository
import com.log3900.utils.ui.getAvatarID
import io.reactivex.android.schedulers.AndroidSchedulers
import io.reactivex.schedulers.Schedulers

class PlayerViewHolder : RecyclerView.ViewHolder {
    private lateinit var player: Player
    private var currentPosition: Int = 0
    private var score: Int = 0

    // UI
    private var rootView: View
    private var positionTextView: TextView
    private var playerChip: Chip
    private var playerStatusImageView: ImageView
    private var scoreTextView: TextView
    private var scoreChangedTextView: TextView


    constructor(itemView: View) : super(itemView) {
        rootView = itemView
        positionTextView = itemView.findViewById(R.id.list_item_active_match_player_text_view_position)
        playerChip = itemView.findViewById(R.id.list_item_active_match_player_chip_player)
        playerStatusImageView = itemView.findViewById(R.id.list_item_active_match_player_image_view_status)
        scoreTextView = itemView.findViewById(R.id.list_item_active_match_player_text_view_score)
        scoreChangedTextView = itemView.findViewById(R.id.list_item_active_match_player_text_view_updated_score)

        playerChip.setOnClickListener {
            PlayerProfileDialogFragment.show(itemView.context, player.ID)
        }
    }

    fun bind(player: Player, position: Int, score: Int, statusRes: Int? = null) {
        this.player = player
        this.currentPosition = position
        this.score = score

        positionTextView.text = "#$position"
        playerChip.text = player.username
        if (statusRes == null) {
            playerStatusImageView.visibility = View.INVISIBLE
        } else {
            playerStatusImageView.setImageResource(statusRes)
            playerStatusImageView.visibility = View.VISIBLE
        }

        scoreTextView.text = "Score: $score"

        if (player.isCPU) {
            playerChip.chipIcon = ResourcesCompat.getDrawable(
                MainApplication.instance.resources,
                R.drawable.ic_robot,
                null
            )
        } else {
            UserRepository.getInstance().getUser(player.ID)
                .subscribeOn(Schedulers.io())
                .observeOn(AndroidSchedulers.mainThread())
                .subscribe(
                    {
                        playerChip.chipIcon = ResourcesCompat.getDrawable(
                            MainApplication.instance.resources,
                            getAvatarID(it.pictureID),
                            null
                        )
                    },
                    {
                    }
                )
        }
    }

    fun showScoreChangedAnimation(scoreChangedValue: String, isPositive: Boolean) {
        if (isPositive) {
            scoreChangedTextView.setTextColor(Color.parseColor("#FF008000"))
        } else {
            scoreChangedTextView.setTextColor(Color.RED)
        }

        scoreChangedTextView.text = scoreChangedValue
        scoreChangedTextView.bringToFront()

        val scaleUpAnimator = ObjectAnimator.ofPropertyValuesHolder(
            scoreChangedTextView,
            PropertyValuesHolder.ofFloat(View.SCALE_X, 1f, 2f),
            PropertyValuesHolder.ofFloat(View.SCALE_Y, 1f, 2f)
        )
        scaleUpAnimator.duration = 2000

        val alphaChangeAnimator = ObjectAnimator.ofPropertyValuesHolder(
            scoreChangedTextView,
            PropertyValuesHolder.ofFloat(View.ALPHA, 0f, 1f)
        )
        alphaChangeAnimator.repeatCount = 1
        alphaChangeAnimator.repeatMode = ObjectAnimator.REVERSE
        alphaChangeAnimator.duration = 1000
        alphaChangeAnimator.start()
        scaleUpAnimator.start()
    }
}