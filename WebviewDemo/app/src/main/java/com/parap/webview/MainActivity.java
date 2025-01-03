package com.parap.webview;

import android.net.http.SslError;
import android.os.Build;
import android.os.Bundle;
import android.view.View;
import android.view.WindowManager;
import android.webkit.SslErrorHandler;
import android.webkit.WebResourceRequest;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;

import androidx.appcompat.app.AppCompatActivity;
import androidx.swiperefreshlayout.widget.SwipeRefreshLayout;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;

public class MainActivity extends AppCompatActivity {
    private WebView webView;
    private SwipeRefreshLayout swipeRefreshLayout;

    private void loadWebPage() {
        // 首先尝试读取url.txt
        String url = readFileFromAssets("url.txt");
        if (url != null && !url.trim().isEmpty()) {
            webView.loadUrl(url.trim());
            return;
        }

        // 如果url.txt不存在或为空，查找并加载HTML文件
        String htmlFile = findHtmlFile();
        if (htmlFile != null) {
            webView.loadUrl("file:///android_asset/" + htmlFile);
            return;
        }

        // 最后尝试读取default.txt
        String defaultUrl = readFileFromAssets("default.txt");
        if (defaultUrl != null && !defaultUrl.trim().isEmpty()) {
            webView.loadUrl(defaultUrl.trim());
        }
    }

    private String findHtmlFile() {
        try {
            String[] files = getAssets().list("");
            // 优先查找index.html
            for (String file : files) {
                if (file.equals("index.html")) {
                    return file;
                }
            }
            // 如果没有index.html，返回找到的第一个html文件
            for (String file : files) {
                if (file.endsWith(".html")) {
                    return file;
                }
            }
        } catch (IOException e) {
            // 忽略错误
        }
        return null;
    }

    private String readFileFromAssets(String fileName) {
        try {
            BufferedReader reader = new BufferedReader(new InputStreamReader(getAssets().open(fileName)));
            String url = reader.readLine();
            reader.close();
            return url;
        } catch (IOException e) {
            return null;
        }
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        
        // 设置全屏并使用刘海区域
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.P) {
            WindowManager.LayoutParams lp = getWindow().getAttributes();
            lp.layoutInDisplayCutoutMode = WindowManager.LayoutParams.LAYOUT_IN_DISPLAY_CUTOUT_MODE_SHORT_EDGES;
            getWindow().setAttributes(lp);
        }
        
        // 隐藏系统UI
        getWindow().getDecorView().setSystemUiVisibility(
            View.SYSTEM_UI_FLAG_LAYOUT_STABLE
            | View.SYSTEM_UI_FLAG_LAYOUT_HIDE_NAVIGATION
            | View.SYSTEM_UI_FLAG_LAYOUT_FULLSCREEN
            | View.SYSTEM_UI_FLAG_HIDE_NAVIGATION
            | View.SYSTEM_UI_FLAG_FULLSCREEN
            | View.SYSTEM_UI_FLAG_IMMERSIVE_STICKY);
            
        setContentView(R.layout.activity_main);

        webView = findViewById(R.id.webview);
        swipeRefreshLayout = findViewById(R.id.swipeRefresh);

        WebSettings settings = webView.getSettings();
        settings.setJavaScriptEnabled(true);
        // 允许访问文件
        settings.setAllowFileAccess(true);
        // 设置可以访问assets目录
        settings.setAllowFileAccessFromFileURLs(true);
        settings.setAllowUniversalAccessFromFileURLs(true);
        
        // SSL/TLS 相关设置
        settings.setDomStorageEnabled(true);
        settings.setDatabaseEnabled(true);
        settings.setLoadsImagesAutomatically(true);
        // 完全禁用 SSL 证书检查
        settings.setUserAgentString("Mozilla/5.0 (Linux; Android 10; Mobile) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36");
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            settings.setMixedContentMode(WebSettings.MIXED_CONTENT_ALWAYS_ALLOW);
        }
        
        // 添加更多设置
        settings.setJavaScriptCanOpenWindowsAutomatically(true);
        settings.setLoadWithOverviewMode(true);
        settings.setUseWideViewPort(true);
        settings.setSupportMultipleWindows(true);
        settings.setBuiltInZoomControls(true);
        settings.setDisplayZoomControls(false);
        settings.setSaveFormData(true);
        settings.setAppCacheEnabled(true);
        settings.setCacheMode(WebSettings.LOAD_DEFAULT);
        
        webView.setWebViewClient(new WebViewClient() {
            @Override
            public boolean shouldOverrideUrlLoading(WebView view, WebResourceRequest request) {
                // 新版本的实现（API 24及以上）
                return shouldOverrideUrlLoadingCommon(request.getUrl().toString());
            }

            private boolean shouldOverrideUrlLoadingCommon(String url) {
                // 只允许加载 http 和 https 开头的链接
                if (url.startsWith("http://") || url.startsWith("https://")) {
                    return false; // 返回 false 让 WebView 正常加载
                }
                // 拦截其他所有协议的链接
                return true; // 返回 true 表示已处理该链接
            }

            @Override
            public void onReceivedSslError(WebView view, SslErrorHandler handler, SslError error) {
                // 完全忽略所有SSL证书错误
                handler.proceed();
            }


            @Override
            public void onPageFinished(WebView view, String url) {
                swipeRefreshLayout.setRefreshing(false);
            }
        });
        
        // 设置下拉刷新的监听器
        swipeRefreshLayout.setOnRefreshListener(new SwipeRefreshLayout.OnRefreshListener() {
            @Override
            public void onRefresh() {
                webView.reload();
            }
        });
        
        // 加载网页
        loadWebPage();
    }

    @Override
    public void onBackPressed() {
        if (webView.canGoBack()) {
            webView.goBack();
        } else {
            super.onBackPressed();
        }
    }
}