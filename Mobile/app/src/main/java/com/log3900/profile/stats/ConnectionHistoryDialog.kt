package com.log3900.profile.stats

import android.app.Dialog
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.DialogFragment
import androidx.recyclerview.widget.DividerItemDecoration
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.google.android.material.button.MaterialButton
import com.log3900.R
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

class ConnectionHistoryDialog : DialogFragment() {
    override fun onStart() {
        super.onStart()
        dialog?.window?.setLayout(
            1200,
            1000
        )
    }

    override fun onCreateDialog(savedInstanceState: Bundle?): Dialog {
        val dialog = super.onCreateDialog(savedInstanceState)
        dialog.setCanceledOnTouchOutside(false)
        return dialog
    }

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        val rootView = inflater.inflate(R.layout.dialog_connection_history, container, false)
        setUpUi(rootView)
        return rootView
    }

    private fun setUpUi(root: View) {
        val closeButton: MaterialButton = root.findViewById(R.id.close_dialog_button)
        closeButton.setOnClickListener {
            dismiss()
        }

        setUpMessagesRV(root)
    }

    private fun setUpMessagesRV(root: View) {
        GlobalScope.launch {
            withContext(Dispatchers.Main) {
                val rvConnections: RecyclerView = root.findViewById(R.id.rv_connections)
//                    val connections = StatsRepository.getConnectionHistory()
                val connections = getConnectionHistory()
                val connectionAdapter = ConnectionAdapter(connections)

                rvConnections.apply {
                    adapter = connectionAdapter
                    layoutManager = LinearLayoutManager(activity, LinearLayoutManager.VERTICAL, false)
                    // Divider between items
                    val divItem = DividerItemDecoration(rvConnections.context,
                        (layoutManager as LinearLayoutManager).orientation)
                    divItem.setDrawable(context.getDrawable(R.drawable.divider)!!)
                    addItemDecoration(divItem)
                }
            }
        }
    }

    private suspend fun getConnectionHistory(): List<Connection> {
        // TODO: Error handling
        return StatsRepository.getConnectionHistory()
    }
}