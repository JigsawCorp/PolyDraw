package com.log3900

import android.content.Intent
import android.media.Image
import android.os.Bundle
import android.os.Handler
import android.view.MenuItem
import android.view.View
import android.widget.ImageView

import androidx.appcompat.app.AppCompatActivity
import androidx.appcompat.widget.Toolbar
import androidx.core.view.GravityCompat
import androidx.drawerlayout.widget.DrawerLayout
import androidx.fragment.app.Fragment
import androidx.navigation.findNavController
import androidx.navigation.ui.navigateUp
import androidx.navigation.ui.AppBarConfiguration
import androidx.navigation.ui.setupActionBarWithNavController
import androidx.navigation.ui.setupWithNavController
import com.google.android.material.floatingactionbutton.FloatingActionButton
import com.google.android.material.navigation.NavigationView
import com.log3900.draw.DrawViewFragment
import com.log3900.login.LoginActivity
import com.log3900.profile.ProfileFragment
import com.log3900.settings.SettingsActivity
import com.log3900.settings.theme.ThemeManager


import com.log3900.socket.SocketService
import com.log3900.tutorial.TutorialActivity
import com.log3900.ui.home.HomeFragment
import com.log3900.user.account.AccountRepository

open class MainActivity : AppCompatActivity() {

    private lateinit var appBarConfiguration: AppBarConfiguration
    private lateinit var drawerLayout: DrawerLayout

    override fun onCreate(savedInstanceState: Bundle?) {
        ThemeManager.applyTheme(this)
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        setupUI()

        val toolbar:Toolbar = findViewById(R.id.toolbar)
        setSupportActionBar(toolbar)

        val hideShowMessage: FloatingActionButton = findViewById(R.id.hideShowMessage)
        hideShowMessage.setOnClickListener{
            var chatView = (supportFragmentManager.findFragmentById(R.id.fragment) as Fragment).view
            when(chatView!!.visibility){
                View.INVISIBLE -> chatView.visibility = View.VISIBLE
                View.VISIBLE -> chatView.visibility = View.INVISIBLE
            }
        }


        drawerLayout = findViewById(R.id.drawer_layout)
        val navView: NavigationView = findViewById(R.id.nav_view)
        val navController = findNavController(R.id.nav_host_fragment)

        appBarConfiguration = AppBarConfiguration(
            setOf(
                R.id.nav_home
            ), drawerLayout)

        setupActionBarWithNavController(navController, appBarConfiguration)

        navView.setupWithNavController(navController)
        navView.menu.findItem(R.id.nav_home).setOnMenuItemClickListener {
            startFragment("HOME_VIEW_FRAGMENT", HomeFragment())
            true
        }

        navView.menu.findItem(R.id.nav_draw).setOnMenuItemClickListener {
            startFragment("DRAW_VIEW_FRAGMENT", DrawViewFragment())
            true
        }

        navView.menu.findItem(R.id.logoutButton).setOnMenuItemClickListener { item: MenuItem? ->
            when (item!!.itemId) {
                R.id.logoutButton -> logout()
            }
            true
        }

        navView.menu.findItem(R.id.nav_profile).setOnMenuItemClickListener {
            startFragment("PROFILE_VIEW_FRAGMENT_TAG", ProfileFragment())
            true
        }

        // Header
        val header = navView.getHeaderView(0)
        val avatar: ImageView = header.findViewById(R.id.nav_header_avatar)
        avatar.setOnClickListener {
            startFragment("PROFILE_VIEW_FRAGMENT_TAG", ProfileFragment())
        }

        if (!AccountRepository.getInstance().getAccount().tutorialDone) {
            startActivity(Intent(applicationContext, TutorialActivity::class.java))
        }
    }

    override fun onResume() {
        super.onResume()
        if (ThemeManager.hasActivityThemeChanged(this)) {
            println("Recreating activity")
            this.recreate()
        }
    }

    override fun onSupportNavigateUp(): Boolean {
        val navController = findNavController(R.id.nav_host_fragment)

        return navController.navigateUp(appBarConfiguration) || super.onSupportNavigateUp()
    }

    private fun setupUI() {
        findViewById<ImageView>(R.id.app_bar_main_settings).setOnClickListener {
            startActivity(Intent(this, SettingsActivity::class.java))
        }
        findViewById<ImageView>(R.id.app_bar_main_tutorial).setOnClickListener {
            startActivity(Intent(this, TutorialActivity::class.java))
        }
    }

    private fun logout() {
        SocketService.instance?.disconnectSocket(Handler {
            intent = Intent(this, LoginActivity::class.java)
            intent.flags = Intent.FLAG_ACTIVITY_CLEAR_TOP or Intent.FLAG_ACTIVITY_CLEAR_TASK or Intent.FLAG_ACTIVITY_NEW_TASK
            startActivity(intent)
            finish()
            true
        })
    }

    override fun onBackPressed() {
        super.onBackPressed()
        logout()
    }

    private fun startFragment(tag: String, fragment: Fragment) {
        drawerLayout.closeDrawer(GravityCompat.START)
        val fragmentManager = supportFragmentManager
        if (fragmentManager.findFragmentByTag(tag) != null)
            return

        fragmentManager.beginTransaction().apply {
            replace(R.id.nav_host_fragment, fragment, tag)
            commit()
        }
    }
}
