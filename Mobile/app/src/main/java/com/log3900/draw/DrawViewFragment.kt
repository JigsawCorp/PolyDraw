package com.log3900.draw

import android.animation.Animator
import android.animation.AnimatorListenerAdapter
import android.annotation.SuppressLint
import android.graphics.Color
import android.graphics.Paint
import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.MotionEvent
import android.view.View
import android.view.ViewGroup
import android.widget.SeekBar
import androidx.fragment.app.Fragment
import com.log3900.R
import com.log3900.draw.divyanshuwidget.DrawMode
import com.log3900.draw.divyanshuwidget.PaintOptions
import com.log3900.game.match.MatchManager
import kotlinx.android.synthetic.main.fragment_draw_tools.*
import kotlinx.android.synthetic.main.fragment_draw_view.*
import kotlinx.android.synthetic.main.view_draw_color_palette.*
import nl.dionsegijn.konfetti.KonfettiView
import nl.dionsegijn.konfetti.models.Shape
import nl.dionsegijn.konfetti.models.Size
import java.util.*

// See https://github.com/divyanshub024/AndroidDraw
// and https://android.jlelse.eu/a-guide-to-drawing-in-android-631237ab6e28

class DrawViewFragment(private var canDraw: Boolean = true) : Fragment() {
    lateinit var drawView: DrawViewBase
    var currentColor: Int = 0
    private lateinit var konfettiView: KonfettiView

    private val DEFAULT_WIDTH_PROGRESS = 8

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val root = inflater.inflate(R.layout.fragment_draw_view, container, false)
        setUpUi(root)
        return root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        setUpFab()
        setUpToolButtons()
        setUpWidthSeekbar()
        setUpColorButtons()

        enableDrawFunctions(canDraw)
    }

    override fun onDestroyView() {
        drawView.onDestroy()
        super.onDestroyView()
    }

    @SuppressLint("ClickableViewAccessibility")
    private fun setUpUi(root: View) {
        drawView = root.findViewById(R.id.draw_view_canvas)
        currentColor = resources.getColor(R.color.color_draw_black, null)
        konfettiView = root.findViewById(R.id.fragment_draw_view_konfetti)
    }

    private fun setUpFab() {
        draw_tools_fab.setOnClickListener {
            draw_tools_view.apply {
                if (visibility == View.GONE) {
                    visibility = View.VISIBLE
                    animate()
                        .translationX(0f)
                        .alpha(1.0f)
                        .setListener(null)
                } else {
                    animate()
                        .translationX(100f)
                        .alpha(0f)
                        .setListener(object : AnimatorListenerAdapter() {
                            override fun onAnimationEnd(animation: Animator?) {
                                super.onAnimationEnd(animation)
                                visibility = View.GONE
                            }
                        })
                }
            }
        }
    }

    private fun setUpToolButtons() {
        // Drawing selected by default
        draw_button.isPressed = true

        draw_button.setOnTouchListener { _, motionEvent ->
            if (motionEvent.action == MotionEvent.ACTION_DOWN) {
                startDrawMode()
            }
            true
        }
        remove_button.setOnTouchListener { view, motionEvent ->
            if (motionEvent.action == MotionEvent.ACTION_DOWN) {
                updateDrawToolButtonPressed(view)
                drawView.setDrawMode(DrawMode.REMOVE)
            }
            true
        }
        erase_button.setOnTouchListener { view, motionEvent ->
            if (motionEvent.action == MotionEvent.ACTION_DOWN) {
                updateDrawToolButtonPressed(view)
                drawView.setDrawMode(DrawMode.ERASE)
            }
            true
        }

        // Circle tip selected by default
        circle_tip_button.isPressed = true
        circle_tip_button.setOnTouchListener { view, motionEvent ->
            if (motionEvent.action == MotionEvent.ACTION_DOWN) {
                updateTipButtonPressed(view)
                drawView.setCap(Paint.Cap.ROUND)
            }
            true
        }

        square_tip_button.setOnTouchListener { view, motionEvent ->
            if (motionEvent.action == MotionEvent.ACTION_DOWN) {
                updateTipButtonPressed(view)
                drawView.setCap(Paint.Cap.SQUARE)
            }
            true
        }
    }

    private fun startDrawMode() {
        if (drawView.getDrawMode() == DrawMode.DRAW)
            return

        updateDrawToolButtonPressed(draw_button)
        drawView.setDrawMode(DrawMode.DRAW)
        changeDrawColor(currentColor)
    }

    private fun setUpWidthSeekbar() {
        seekbar_width.setOnSeekBarChangeListener(object: SeekBar.OnSeekBarChangeListener{
            override fun onProgressChanged(seekBar: SeekBar?, progress: Int, fromUser: Boolean) {
                val newProgress = if (progress < 1) 1 else progress
                drawView.setStrokeWidth(newProgress.toFloat())
            }

            override fun onStartTrackingTouch(seekBar: SeekBar?) {}

            override fun onStopTrackingTouch(seekBar: SeekBar?) {}
        })
    }

    private fun setUpColorButtons() {
        // Black by default
        updateColorScale(color_picker_black)

        color_picker_black.setOnClickListener {
            onColorButtonPressed(R.color.color_draw_black)
            updateColorScale(it)
        }
        color_picker_white.setOnClickListener {
            onColorButtonPressed(R.color.color_draw_white)
            updateColorScale(it)
        }
        color_picker_red.setOnClickListener {
            onColorButtonPressed(R.color.color_draw_red)
            updateColorScale(it)
        }
        color_picker_green.setOnClickListener {
            onColorButtonPressed(R.color.color_draw_green)
            updateColorScale(it)
        }
        color_picker_blue.setOnClickListener {
            onColorButtonPressed(R.color.color_draw_blue)
            updateColorScale(it)
        }
        color_picker_yellow.setOnClickListener {
            onColorButtonPressed(R.color.color_draw_yellow)
            updateColorScale(it)
        }
        color_picker_cyan.setOnClickListener {
            onColorButtonPressed(R.color.color_draw_cyan)
            updateColorScale(it)
        }
        color_picker_magenta.setOnClickListener {
            onColorButtonPressed(R.color.color_draw_magenta)
            updateColorScale(it)
        }
    }

    private fun onColorButtonPressed(color: Int) {
        currentColor = resources.getColor(color, null)
        changeDrawColor(currentColor)
        startDrawMode()
    }

    private fun updateDrawToolButtonPressed(button: View) {
        draw_button.isPressed = false
        remove_button.isPressed = false
        erase_button.isPressed = false

        button.isPressed = true
    }

    private fun updateTipButtonPressed(button: View) {
        circle_tip_button.isPressed = false
        square_tip_button.isPressed = false

        button.isPressed = true
    }

    private fun updateColorScale(colorPicked: View) {
        // Reset all colors
        color_picker_black.scaleX = 1f
        color_picker_black.scaleY = 1f

        color_picker_white.scaleX = 1f
        color_picker_white.scaleY = 1f

        color_picker_red.scaleX = 1f
        color_picker_red.scaleY = 1f

        color_picker_green.scaleX = 1f
        color_picker_green.scaleY = 1f

        color_picker_blue.scaleX = 1f
        color_picker_blue.scaleY = 1f

        color_picker_yellow.scaleX = 1f
        color_picker_yellow.scaleY = 1f

        color_picker_cyan.scaleX = 1f
        color_picker_cyan.scaleY = 1f

        color_picker_magenta.scaleX = 1f
        color_picker_magenta.scaleY = 1f

        // Scale up the selected color
        colorPicked.isPressed = true
        colorPicked.animate()
            .scaleX(1.5f)
            .scaleY(1.5f)
    }

    private fun changeDrawColor(color: Int) {
        drawView.setColor(color)
    }

    fun enableDrawFunctions(enable: Boolean, drawingID: UUID? = null, matchManager: MatchManager? = null) {
        canDraw = enable
        if (canDraw) {
            setDrawToolsVisibility(View.VISIBLE)
            setToDefaultValues()
        } else {
            setDrawToolsVisibility(View.GONE)
        }
        Log.d("DRAW_CANVAS", drawingID.toString())
        drawView.enableCanDraw(canDraw, drawingID)
        drawView.socketDrawingReceiver?.matchManager = matchManager
    }

    fun setToDefaultValues() {
        updateDrawToolButtonPressed(draw_button)
        updateTipButtonPressed(circle_tip_button)
        updateColorScale(color_picker_black)
        seekbar_width.progress = DEFAULT_WIDTH_PROGRESS
        drawView.setOptions(PaintOptions())
    }

    fun showConfetti() {
        konfettiView.build()
            .addColors(Color.YELLOW, Color.GREEN, Color.RED, Color.BLUE)
            .setDirection(0.0, 359.0)
            .setSpeed(1f, 5f)
            .setFadeOutEnabled(true)
            .setTimeToLive(2000L)
            .addShapes(Shape.Square, Shape.Circle)
            .addSizes(Size(12))
            .setPosition(-50f, konfettiView.width + 50f, -50f, -50f)
            .streamFor(300, 5000L)
    }

    fun stopConfetti() {
        konfettiView.getActiveSystems().forEach {
            it.stop()
        }
    }

    private fun setDrawToolsVisibility(visibility: Int) {
        draw_tools_view.visibility = visibility
        draw_tools_fab.visibility = visibility
    }

    fun clearCanvas() {
        drawView.clearCanvas()
    }
}