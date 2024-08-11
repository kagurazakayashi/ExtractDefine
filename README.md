![icon](ico/icon.ico)

# [ExtractDefine](https://github.com/kagurazakayashi/ExtractDefine)

[ English | [中文](#C项目生效宏定义提取工具)]

## C project effective macro definition extraction tool

This tool can help you find out which macros are defined in a C program.

Recognizable keywords:

- `*.c`, `*.h`:
  - `#define`, `#undef`
  - `#if`, `defined`, `#endif`, `||`, `&&`
- `CMakeLists.txt`:
  - `set`
  - `includes`
  - `srcs`
  - `SDKCONFIG_DEFAULTS`
  - `EXTRA_COMPONENT_DIRS`
  - `idf_component_register`

### Parameters

- `-c`: Specify a configuration file (.yaml format, will override command line parameters).

- `-i`: Specify a CMakeLists.txt file, and start searching for macro definitions from here.
- `-d`: log display level: `0`=debug, `1`=info, `2`=success, `3`=warning, `4`=error, `5`=failure, `6`=no log output
- `-f`: only need information of these macros (separated by ,)
- `-h`: show command help

### Configuration file

- See the sample configuration file: [ExtractDefine.yaml](ExtractDefine.yaml)
- Language text file: `(executable file name).i18n.ini`

### Output

- log information output to: Stderr
- result output to: Stdout
  - Result format: `macroName=macroContent`, each line

### License

[Mulan Permissive Software License，Version 2 (Mulan PSL v2)](http://license.coscl.org.cn/MulanPSL2)

## C项目生效宏定义提取工具

这个小工具可以帮你找出一个C程序中都定义了哪些宏。

能够识别的关键词：

- `*.c`, `*.h`:
  - `#define`, `#undef`
  - `#if`, `defined`, `#endif`, `||`, `&&`
- `CMakeLists.txt`:
  - `set`
  - `includes`
  - `srcs`
  - `SDKCONFIG_DEFAULTS`
  - `EXTRA_COMPONENT_DIRS`
  - `idf_component_register`

### 启动参数

- `-c`: 指定一个配置文件 (.yaml 格式, 会覆盖命令行参数)。
- `-i`: 指定一个 CMakeLists.txt 文件，将从这里开始搜索宏定义。
- `-d`: 日志显示级别: `0`=调试, `1`=信息, `2`=成功, `3`=警告, `4`=错误, `5`=失败, `6`=不输出日志
- `-f`: 只需要这些宏的信息(用 , 分隔)
- `-h`: 显示命令帮助

### 配置文件

- 示例配置文件见: [ExtractDefine.yaml](ExtractDefine.yaml)
- 语言文本文件: `(可执行文件名).i18n.ini`

### 输出

- 日志信息输出到: Stderr
- 结果输出到: Stdout
  - 结果格式: `宏名称=宏内容`, 每个一行

### 协议

[木兰宽松许可证， 第2版](http://license.coscl.org.cn/MulanPSL2)
