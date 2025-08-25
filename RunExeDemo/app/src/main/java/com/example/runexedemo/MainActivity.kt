package com.example.runexedemo

import android.content.Context
import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material3.Button
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.PlayArrow
import androidx.compose.material.icons.outlined.Settings
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.compose.foundation.text.selection.SelectionContainer
import androidx.compose.ui.tooling.preview.Preview
import com.example.runexedemo.ui.theme.RunExeDemoTheme
import java.io.File
import java.io.FileOutputStream
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.DisposableEffect
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.runtime.rememberCoroutineScope
import kotlinx.coroutines.launch

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            RunExeDemoTheme {
                var selectedTab by remember { mutableStateOf(0) }
                Scaffold(
                    modifier = Modifier.fillMaxSize(),
                    bottomBar = {
                        NavigationBar(containerColor = MaterialTheme.colorScheme.surfaceVariant) {
                            NavigationBarItem(
                                selected = selectedTab == 0,
                                onClick = { selectedTab = 0 },
                                label = { Text("Run") },
                                icon = { Icon(Icons.Outlined.PlayArrow, contentDescription = "Run") },
                                colors = NavigationBarItemDefaults.colors(
                                    selectedIconColor = MaterialTheme.colorScheme.onPrimaryContainer,
                                    selectedTextColor = MaterialTheme.colorScheme.onSurface,
                                    indicatorColor = MaterialTheme.colorScheme.primaryContainer,
                                    unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                                    unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant
                                )
                            )
                            NavigationBarItem(
                                selected = selectedTab == 1,
                                onClick = { selectedTab = 1 },
                                label = { Text("Config") },
                                icon = { Icon(Icons.Outlined.Settings, contentDescription = "Config") },
                                colors = NavigationBarItemDefaults.colors(
                                    selectedIconColor = MaterialTheme.colorScheme.onPrimaryContainer,
                                    selectedTextColor = MaterialTheme.colorScheme.onSurface,
                                    indicatorColor = MaterialTheme.colorScheme.primaryContainer,
                                    unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                                    unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant
                                )
                            )
                        }
                    }
                ) { innerPadding ->
                    when (selectedTab) {
                        0 -> RunScreen(modifier = Modifier.padding(innerPadding))
                        else -> ConfigScreen(modifier = Modifier.padding(innerPadding))
                    }
                }
            }
        }
    }
}

@Composable
fun RunScreen(modifier: Modifier = Modifier) {
    val ctx = LocalContext.current
    var output by remember { mutableStateOf("") }
    var running by remember { mutableStateOf(false) }
    var handle by remember { mutableStateOf<BinaryRunner.StreamingHandle?>(null) }
    val scrollState = rememberScrollState()
    val coroutineScope = rememberCoroutineScope()
    // discover .so files from nativeLibraryDir
    val nativeLibDir = ctx.applicationInfo.nativeLibraryDir
    val exeFiles = remember(nativeLibDir) {
        File(nativeLibDir)
            .listFiles { f -> f.isFile }
            ?.map { it.name }
            ?.filter { name -> !name.startsWith("libandroidx") }
            ?: emptyList()
    }
    var expanded by remember { mutableStateOf(false) }
    var selectedSo by remember { mutableStateOf(exeFiles.firstOrNull() ?: "") }
    Column(modifier = modifier.padding(16.dp)) {
        Box(modifier = Modifier.fillMaxWidth()) {
            OutlinedTextField(
                readOnly = true,
                value = selectedSo,
                onValueChange = {},
                label = { Text("选择 .so") },
                trailingIcon = {
                    Icon(
                        imageVector = Icons.Outlined.Settings,
                        contentDescription = null,
                        modifier = Modifier
                            .padding(4.dp)
                            .clickable { expanded = !expanded }
                    )
                },
                modifier = Modifier
                    .fillMaxWidth()
                    .clickable { expanded = !expanded }
            )
            DropdownMenu(expanded = expanded, onDismissRequest = { expanded = false }) {
                exeFiles.forEach { name ->
                    DropdownMenuItem(text = { Text(name) }, onClick = {
                        selectedSo = name
                        expanded = false
                    })
                }
            }
        }
        Button(onClick = {
            if (!running) {
                output = ""
                running = true
                val cfg = ConfigManager.ensureConfig(ctx)
                handle = BinaryRunner.runStreaming(
                    context = ctx,
                    libName = selectedSo,
                    "--config", cfg.absolutePath,
                    onLine = { line ->
                        output += if (output.isEmpty()) line else "\n$line"
                    },
                    onExit = { code ->
                        output += if (output.isEmpty()) "exit=$code" else "\nexit=$code"
                        running = false
                    }
                )
            } else {
                handle?.stop()
            }
        }) {
            Text(if (running) "Stop" else "Run")
        }
        SelectionContainer {
            Text(output, modifier = Modifier
                .padding(top = 12.dp)
                .verticalScroll(scrollState)
            )
        }
        LaunchedEffect(output) {
            if (output.isNotEmpty()) {
                coroutineScope.launch {
                    scrollState.animateScrollTo(scrollState.maxValue)
                }
            }
        }
    }
}

@Preview(showBackground = true)
@Composable
fun RunScreenPreview() {
    RunExeDemoTheme {
        RunScreen()
    }
}

@Composable
fun ConfigScreen(modifier: Modifier = Modifier) {
    val ctx = LocalContext.current
    var configText by remember { mutableStateOf("") }
    LaunchedEffect(Unit) {
        configText = ConfigManager.readConfig(ctx)
    }
    Column(modifier = modifier.fillMaxSize().padding(16.dp)) {
        OutlinedTextField(
            value = configText,
            onValueChange = {
                configText = it
                ConfigManager.writeConfig(ctx, configText)
            },
            label = { Text("config.json") },
            modifier = Modifier
                .padding(top = 12.dp)
                .fillMaxWidth()
                .weight(1f)
        )
    }
}