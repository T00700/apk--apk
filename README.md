# 通过网页和文件生成一个简单app
之前发过一个使用[apktool](https://github.com/pzx521521/html2apk)版本的,简单来说就是一个对apk反编译和编译的过程
但是apktool还是有点问题的:
+ [x] 需要下载jdk/jre,有环境依赖
+ [x] 速度较慢
+ [x] 需要中间调用apktool,不容易部署到服务器/docker
# 使用
## 修改apk显示的网页
+ 在线网址
```shell
./apkEditor -o="/Users/parapeng/Downloads/app-new.apk" https://www.example.com 
```
+ 指定输出路径
```shell
./apkEditor -o=demo.apk https://www.example.com 
```
+ 本地文件 仅一个index.html
```shell
apkEditor <yourpath>/index.html
```
+ 本地文件夹 包含html+css+js
```shell
apkEditor <your-dir>
```
+ 本地zip文件 包含html+css+js的zip文件
```shell
apkEditor <your-dir>/demo.zip
```

## 修改其他信息
在修改显示的网页的基础上添加
+ label
  对应application.label  
  用于显示软件名
+ versionCode
  对应manifest.android:versionCode  
  版本号用于更新软件
+ versionName
  对应manifest.android:versionName  
  用于显示软件版本号
+ package  
  ~~对应manifest.android:package    
  软件包名~~
  修改这个要改的东西太多了,请使用apktool
+ 生成默认的webview并修改信息
```shell
./apkEditor -versionCode=222 -versionName="2.2.2" -label="NewApp" -o="/Users/parapeng/Downloads/app-new.apk" https://www.example.com
```

# 原理
## 反编译apk正常的流程是:
+ 解压apk  
    ```shell
    unzip origin.apk -d origin/
    ```
+ 修改文件
    ```shell  
    echo "1234" > /tmp/cpid
    cp /tmp/cpid origin/
    ```
+ 重新打包apk
  aapt2+aidl+Renderscript+Javac+DEX+zipflinger
  + 如果不需要修改代码的话是aapt2+zipflinger
  + aapt2会对一些资源做特殊处理,如AndroidManifest.xml会变为一个二进制文件,所以只使用zip是不行的
+ [签名apk](https://android.googlesource.com/platform/build/+/refs/heads/main/tools/signapk/)
  + v1签名 jarsigner
  + v2签名 (Android11+) apksigner, 
    + 必须先 zipalign 在 apksigner,否则会导致签名失效
  + 对于签名v1、v2、zipalign对齐的顺序是：v1----zipalign对齐----v2
+ [zipalign](https://developer.android.com/tools/zipalign?hl=zh-cn)对齐  
  + Android12+（API31+）必须的
  + 源码[cpp](https://android.googlesource.com/platform/build/+/refs/heads/main/tools/zipalign/)

apktool 相当于把`重新打包apk`简化了,并多了解析dex文件
## 实现原理
由于并不需要修改dex中的东西  
解压->修改->压缩->对齐->签名  
zipalign 困了我好久最好找到源码才知道怎么实现,然后还不知道问什么这么实现
[源码](https://github.com/pzx521521/apkEditor)
[在线体验](https://github.com/pzx521521/apkEditor)
# 参考引用:  
[zipmerge](https://github.com/rsc/zipmerge)  
[signv2](https://github.com/morrildl/playground-android)
# todo
+ [ ] 包名的修改
+ [ ] 图标修改
+ [ ] 桌面App(UI)
+ [ ] 对其他的app的修改