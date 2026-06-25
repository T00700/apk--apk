package com.example.runexedemo

import android.content.Context
import android.util.Base64
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONObject
import java.net.HttpURLConnection
import java.net.URL

object GithubConfigSync {
    private const val OWNER = "pzx521521"
    private const val REPO = "qdapi"
    private const val BRANCH = "master"
    private const val PATH = "config.json"
    private const val API_URL = "https://api.github.com/repos/$OWNER/$REPO/contents/$PATH"

    // 从 GitHub 下载 config.json，并覆盖本地配置文件。
    suspend fun downloadCurrentConfig(context: Context): Result<String> = withContext(Dispatchers.IO) {
        runCatching {
            val token = BuildConfig.GITHUB_TOKEN
            require(token.isNotBlank()) { "缺少 github.token，请先配置 local.properties" }

            val response = request(
                method = "GET",
                token = token,
                query = "ref=$BRANCH"
            )
            val remoteContent = response.getString("content")
            val decodedContent = String(
                Base64.decode(remoteContent, Base64.DEFAULT),
                Charsets.UTF_8
            )
            ConfigManager.writeConfig(context, decodedContent)
            decodedContent
        }
    }

    // 将本地 config.json 上传到 GitHub Contents API，并覆盖目标文件。
    suspend fun syncCurrentConfig(context: Context): Result<String> = withContext(Dispatchers.IO) {
        runCatching {
            val token = BuildConfig.GITHUB_TOKEN
            require(token.isNotBlank()) { "缺少 github.token，请先配置 local.properties" }

            val content = ConfigManager.readConfig(context)
            val sha = fetchRemoteSha(token)
            val encodedContent = Base64.encodeToString(
                content.toByteArray(Charsets.UTF_8),
                Base64.NO_WRAP
            )
            val body = JSONObject()
                .put("message", "sync config from Android app")
                .put("content", encodedContent)
                .put("sha", sha)
                .put("branch", BRANCH)
                .toString()

            val response = request(
                method = "PUT",
                token = token,
                body = body
            )
            "同步成功：${response.optJSONObject("content")?.optString("html_url").orEmpty()}"
        }
    }

    // 查询远端文件当前 sha，GitHub 更新文件时必须携带该值。
    private fun fetchRemoteSha(token: String): String {
        val response = request(
            method = "GET",
            token = token,
            query = "ref=$BRANCH"
        )
        return response.getString("sha")
    }

    // 统一执行 GitHub API 请求，并在非 2xx 时抛出可读错误。
    private fun request(
        method: String,
        token: String,
        query: String? = null,
        body: String? = null
    ): JSONObject {
        val target = if (query == null) API_URL else "$API_URL?$query"
        val connection = (URL(target).openConnection() as HttpURLConnection).apply {
            requestMethod = method
            connectTimeout = 15_000
            readTimeout = 15_000
            setRequestProperty("Authorization", "Bearer $token")
            setRequestProperty("Accept", "application/vnd.github+json")
            setRequestProperty("X-GitHub-Api-Version", "2022-11-28")
            if (body != null) {
                doOutput = true
                setRequestProperty("Content-Type", "application/json; charset=utf-8")
            }
        }

        try {
            if (body != null) {
                connection.outputStream.use { output ->
                    output.write(body.toByteArray(Charsets.UTF_8))
                }
            }
            val responseText = if (connection.responseCode in 200..299) {
                connection.inputStream.bufferedReader().use { it.readText() }
            } else {
                connection.errorStream?.bufferedReader()?.use { it.readText() }.orEmpty()
            }
            if (connection.responseCode !in 200..299) {
                throw IllegalStateException("GitHub API ${connection.responseCode}: $responseText")
            }
            return JSONObject(responseText)
        } finally {
            connection.disconnect()
        }
    }
}
