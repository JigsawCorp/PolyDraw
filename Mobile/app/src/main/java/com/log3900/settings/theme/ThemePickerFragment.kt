package com.log3900.settings.theme

import android.app.Activity
import android.app.AlertDialog
import android.app.Dialog
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.DialogFragment
import androidx.fragment.app.Fragment
import androidx.recyclerview.widget.GridLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.log3900.R

class ThemePickerFragment : DialogFragment() {
    private lateinit var themesRecyclerView: RecyclerView
    private lateinit var themesAdapter: ThemeAdapter
    private var themes = ThemeManager.getThemesAsArrayList()

    override fun onCreateDialog(savedInstanceState: Bundle?): Dialog {
        val dialogBuilder = AlertDialog.Builder(activity)
            .setTitle("Theme Picker")
            .setPositiveButton("Save") { _, _ ->
                ThemeManager.changeTheme(themesAdapter.selectedTheme)
            }
            .setNegativeButton("Cancel") { _, _ ->

            }

        val view = activity?.layoutInflater?.inflate(R.layout.fragment_theme_picker, null)

        themesRecyclerView = view?.findViewById(R.id.fragment_theme_picker_recycler_view)!!

        setupRecyclerView()

        dialogBuilder.setView(view)
        return dialogBuilder.create()
    }

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        val rootView = inflater.inflate(R.layout.fragment_theme_picker, container, false)

        return rootView
    }

    private fun setupRecyclerView() {
        themesAdapter = ThemeAdapter(themes)
        themesRecyclerView.apply {
            setHasFixedSize(true)
            layoutManager = GridLayoutManager(this.context, 4)
            adapter = themesAdapter
        }
    }
}