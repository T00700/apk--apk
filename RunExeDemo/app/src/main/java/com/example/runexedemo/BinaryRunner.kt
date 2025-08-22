package com.example.runexedemo

import android.content.Context
import android.os.Build
import android.util.Log
import java.io.File
import java.io.InputStreamReader
import java.util.concurrent.Executors

object BinaryRunner {
    private const val LIB_FILE_NAME = "app_arm64.so"

    data class Result(
        val exitCode: Int,
        val combinedOutput: String
    )

    class StreamingHandle internal constructor(
        private val process: Process?
    ) {
        fun stop() {
            try {
                process?.destroy()
            } catch (_: Throwable) {}
        }
    }

    /**
     * Execute and stream stdout+stderr lines in real time.
     * onLine will be invoked on a background thread; if UI updates are needed, post to main thread.
     */
    fun runStreaming(
        context: Context,
        libName: String? = null,
        vararg args: String,
        onLine: (String) -> Unit,
        onExit: (Int) -> Unit
    ): StreamingHandle {
        val nativeLibDir = context.applicationInfo.nativeLibraryDir
        val targetLibName = libName ?: LIB_FILE_NAME
        val libCandidate = File(nativeLibDir, targetLibName)
        if (!libCandidate.exists()) {
            try { onLine("[not found] ${'$'}{libCandidate.absolutePath}") } catch (_: Throwable) {}
            try { onLine("请将二进制放到 jniLibs/arm64-v8a/libexample_arm64.so 并重建") } catch (_: Throwable) {}
            try { onExit(-1) } catch (_: Throwable) {}
            return StreamingHandle(null)
        }

        // Ensure --config is passed unless already provided
        val hasConfigArg = args.asList().indexOf("--config") >= 0
        val finalArgs: List<String> = if (hasConfigArg) {
            args.asList()
        } else {
            val cfgPath = ConfigManager.ensureConfig(context).absolutePath
            listOf("--config", cfgPath) + args.asList()
        }

        val command = buildList {
            add(libCandidate.absolutePath)
            addAll(finalArgs)
        }
        val process = ProcessBuilder(command)
            .directory(context.filesDir)
            .redirectErrorStream(true)
            .start()

        val executor = Executors.newSingleThreadExecutor()
        executor.execute {
            try {
                InputStreamReader(process.inputStream).buffered().use { reader ->
                    var line: String?
                    while (reader.readLine().also { line = it } != null) {
                        try {
                            onLine(line!!)
                        } catch (_: Throwable) {}
                    }
                }
            } catch (t: Throwable) {
                try { onLine("[stream error] $t") } catch (_: Throwable) {}
            } finally {
                val exit = try { process.waitFor() } catch (_: Throwable) { -1 }
                try { onExit(exit) } catch (_: Throwable) {}
                try { executor.shutdown() } catch (_: Throwable) {}
            }
        }

        return StreamingHandle(process)
    }
}

