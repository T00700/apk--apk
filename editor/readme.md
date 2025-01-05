你说的也很有道理
在仔细对比了其他语言的实现并重新看了一下里面的逻辑之后
我也不是很确定这样修改对不对或者说需不需要修改
我应该关闭这pr吗?

我应该怎么使用go语言判断文件是否存在呢?
github上搜`fileexist language:Go`
大部分都用到了os.IsNotExist
所以你问ai:"golang file exist",ai会回答:
```
func FileExists(path string) bool {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    }
    return err == nil
}
```
为什么需要`os.IsNotExist(err)`呢?下面的写法不是更好吗?
```go
func FileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
```
因为err != nil 的时候才能使 os.IsNotExist(err) retrun true




之前我没有仔细看golang里面的代码,
只是凭直觉认为判断一个文件是否存在应该用
if _, err := os.Stat();os.IsNotExist(err){}
因为其他语言都是这样判断的(like python, cpp, java)
其实应该用fileopen或者 if _, err := os.Stat, err != nil {}
只能说os.IsNotExist 这个名字太迷惑人了

所以我参考了其他语言中的处理
有两种方式判断一个文件是否存在
+ os.Open(i don't like it 因为要file.close())
+ os.Stat

在判断文件是否存在的时候,os.Open 相当于 
```
if _, err := os.Stat, err != nil {}
```
比如: ErrPermission时os.Open也会返回错误

|                 | python         | nodejs    | cpp  | java |
|-----------------|----------------|-----------|------|------|
| use             | stat           | stat      | file | file |

应为cpp, java 都是用file来判断的, 这里不做详细说明
让我们来看下python,nodejs

+ python 3.13
判断文件实现的方式,推荐的使用 pathlib 
```
from pathlib import Path
print(Path("http://www.google.com").exists())
```
`3.11/lib/python3.11/pathlib.py Line.1230`
```
    def exists(self):
        """
        Whether this path exists.
        """
        try:
            self.stat()
        except OSError as e:
            if not _ignore_error(e):
                raise
            return False
        except ValueError:
            # Non-encodable path
            return False
        return True
```
可以看出所有错误都会返回false, 
相当于golang中 `if _, err := os.Stat, err != nil {}`


+ nodejs
```
const fs = require('fs');

const filePath = 'http://www.google.com'; 

fs.stat(filePath, (err, stats) => {
  if (err) {
    console.log(err);
    console.log(`The file '${filePath}' does not exist.`);
  } else {
    console.log(`The file '${filePath}' exists.`);
  }
});
```
output
```
[Error: ENOENT: no such file or directory, stat 'http://www.google.com'] {
errno: -2,
code: 'ENOENT',
syscall: 'stat',
path: 'http://www.google.com'
}
The file 'http://www.google.com' does not exist.
```
注意nodejs,为了兼容unix, 在windows下, 以下会返回错误码'ENOENT'
https://github.com/libuv/libuv/blob/ec5a4b54f7da7eeb01679005c615fee9633cdb3b/src/win/error.c#L138
```
    case ERROR_BAD_PATHNAME:                return UV_ENOENT;
    case ERROR_DIRECTORY:                   return UV_ENOENT;
    case ERROR_ENVVAR_NOT_FOUND:            return UV_ENOENT;
    case ERROR_FILE_NOT_FOUND:              return UV_ENOENT;
    case ERROR_INVALID_NAME:                return UV_ENOENT;
    case ERROR_INVALID_DRIVE:               return UV_ENOENT;
    case ERROR_INVALID_REPARSE_DATA:        return UV_ENOENT;
    case ERROR_MOD_NOT_FOUND:               return UV_ENOENT;
    case ERROR_PATH_NOT_FOUND:              return UV_ENOENT;
    case WSAHOST_NOT_FOUND:                 return UV_ENOENT;
    case WSANO_DATA:                        return UV_ENOENT;
```
