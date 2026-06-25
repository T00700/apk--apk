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
import androidx.compose.material3.CircularProgressIndicator
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
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.foundation.horizontalScroll
import androidx.compose.foundation.layout.Row
import androidx.compose.material3.Switch
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.height
import androidx.compose.ui.unit.sp
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.itemsIndexed
import androidx.compose.foundation.background
import androidx.compose.foundation.lazy.rememberLazyListState

private const val RUN_PREFS_NAME = "run_settings"
private const val SELECTED_SO_KEY = "selected_so"

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
    val preferences = remember(ctx) {
        ctx.getSharedPreferences(RUN_PREFS_NAME, Context.MODE_PRIVATE)
    }
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
    var selectedSo by remember(exeFiles) {
        mutableStateOf(
            preferences.getString(SELECTED_SO_KEY, null)
                ?.takeIf { it in exeFiles }
                ?: exeFiles.firstOrNull()
                ?: ""
        )
    }
    var wrap by remember { mutableStateOf(false) }
    var filter by remember { mutableStateOf("") }

    // derive displayed lines with filter
    val lines = remember(output, filter) {
        val src = if (filter.isBlank()) output else output
            .lineSequence()
            .filter { it.contains(filter, ignoreCase = true) }
            .joinToString("\n")
        if (src.isEmpty()) emptyList() else src.split('\n')
    }
    val listState = rememberLazyListState()

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
                        preferences.edit().putString(SELECTED_SO_KEY, name).apply()
                        expanded = false
                    })
                }
            }
        }
        // One-line controls: Run button + wrap toggle + filter field
        Row(modifier = Modifier.padding(top = 8.dp)) {
            Button(onClick = {
                output = "nativeLibraryDir: $nativeLibDir\n"
                if (selectedSo ==""){
                    output += "找不到.so文件"
                    return@Button
                }
                if (!running) {
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
                Text(if (running) "Stop" else "Run", fontSize = 14.sp)
            }
            Spacer(modifier = Modifier.width(12.dp))
            Row {
                Text("换行", fontSize = 12.sp)
                Spacer(modifier = Modifier.width(6.dp))
                Switch(checked = wrap, onCheckedChange = { wrap = it })
            }
            Spacer(modifier = Modifier.width(12.dp))
            OutlinedTextField(
                value = filter,
                onValueChange = { filter = it },
                label = { Text("过滤关键词", fontSize = 12.sp) },
                modifier = Modifier.weight(1f),
                singleLine = true
            )
        }
        SelectionContainer {
            LazyColumn(
                state = listState,
                modifier = Modifier
                    .padding(top = 12.dp)
                    .fillMaxWidth()
            ) {
                itemsIndexed(lines) { index, line ->
                    val bg = if (index % 2 == 0)
                        MaterialTheme.colorScheme.surface
                    else
                        MaterialTheme.colorScheme.surfaceVariant
                    Box(
                        modifier = Modifier
                            .background(bg)
                            .fillMaxWidth()
                            .padding(horizontal = 8.dp, vertical = 2.dp)
                    ) {
                        val lineModifier = if (wrap) Modifier else Modifier.horizontalScroll(rememberScrollState())
                        Text(
                            text = line,
                            softWrap = wrap,
                            fontFamily = FontFamily.Monospace,
                            modifier = lineModifier
                        )
                    }
                }
            }
        }
        LaunchedEffect(lines.size) {
            if (lines.isNotEmpty()) {
                listState.animateScrollToItem(lines.size - 1)
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
    val coroutineScope = rememberCoroutineScope()
    var configText by remember { mutableStateOf("") }
    var working by remember { mutableStateOf(false) }
    var statusMessage by remember { mutableStateOf("") }
    LaunchedEffect(Unit) {
        configText = ConfigManager.readConfig(ctx)
    }
    Column(modifier = modifier.fillMaxSize().padding(16.dp)) {
        Row(modifier = Modifier.fillMaxWidth()) {
            Button(
                enabled = !working,
                onClick = {
                    working = true
                    statusMessage = "正在同步到 GitHub..."
                    ConfigManager.writeConfig(ctx, configText)
                    coroutineScope.launch {
                        val result = GithubConfigSync.syncCurrentConfig(ctx)
                        statusMessage = result.fold(
                            onSuccess = { it },
                            onFailure = { "同步失败：${it.message ?: it::class.java.simpleName}" }
                        )
                        working = false
                    }
                },
                modifier = Modifier.weight(1f)
            ) {
                if (working) {
                    CircularProgressIndicator(
                        modifier = Modifier
                            .padding(end = 8.dp)
                            .height(16.dp)
                            .width(16.dp),
                        strokeWidth = 2.dp
                    )
                }
                Text(if (working) "处理中" else "同步到 GitHub")
            }
            Spacer(modifier = Modifier.width(8.dp))
            Button(
                enabled = !working,
                onClick = {
                    working = true
                    statusMessage = "正在从 GitHub 下载..."
                    coroutineScope.launch {
                        val result = GithubConfigSync.downloadCurrentConfig(ctx)
                        statusMessage = result.fold(
                            onSuccess = {
                                configText = it
                                "下载成功，已更新本地 config.json"
                            },
                            onFailure = { "下载失败：${it.message ?: it::class.java.simpleName}" }
                        )
                        working = false
                    }
                },
                modifier = Modifier.weight(1f)
            ) {
                Text("从 GitHub 下载")
            }
        }
        if (statusMessage.isNotBlank()) {
            Text(
                text = statusMessage,
                modifier = Modifier.padding(top = 8.dp),
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
        }
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
