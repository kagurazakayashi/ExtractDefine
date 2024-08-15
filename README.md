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
- `-d`: Log display level: `0`=debug, `1`=info, `2`=success, `3`=warning, `4`=error, `5`=failure, `6`=no log output
- `-f`: Only need information of these macros (separated by ,)
- `-e`: Pre-set these macros at the beginning of the run (format as key=value and separated by , )
- `-h`: Show command help

### Configuration file

- See the sample configuration file: [ExtractDefine.yaml](ExtractDefine.yaml)
- Language text file: `(executable file name).i18n.ini`

#### Variable default values

During parsing, the `EXTRACTDEFINE` variable is provided to distinguish when using this tool.

- If `CMAKE_C_FLAGS` is not specified,
- The default value will be: empty string.
- If `CMAKE_SOURCE_DIR`, `CMAKE_CURRENT_LIST_DIR`, `ROOT_DIR`, `__root_dir` are not specified,
- The default value will be: the path to the folder where the root `CMakeLists.txt` is located.
- If `CMAKE_CURRENT_SOURCE_DIR`, `COMPONENT_DIR` are not specified,
- The default value will be: the path to the folder where the currently processed `CMakeLists.txt` is located.

If still not found, abort the processing of the current configuration.

### Output

- log information output to: Stderr
- result output to: Stdout
  - Result format: `macroName=macroContent`, each line

#### Example

Assuming there is a configuration file `ExtractDefine.yaml`, to output the run log to `run.log` and the run result to `result.txt`, please use the following command:

`./ExtractDefine -c ExtractDefine.yaml >result.txt 2>run.log`

If you want to output only to `result.txt`, please use the following command:

`./ExtractDefine -c ExtractDefine.yaml >result.txt`

## Compile

`go build .`

For full platform compilation, refer to [build.bat](build.bat) .

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
- `-e`: 在运行开始时预置这些宏(格式为 键=值 并用 , 分隔)
- `-h`: 显示命令帮助

### 配置文件

- 示例配置文件见: [ExtractDefine.example.yaml](ExtractDefine.example.yaml)
- 语言文本文件: `(可执行文件名).i18n.ini`

### `CMakeLists.txt` 中的变量

如果遇到使用 `${}` 定义的变量，将试图从已收集的宏定义中查找，如果没有，则：

#### 变量的默认值

解析时会提供 `EXTRACTDEFINE` 变量，用于区分在使用本工具。

- 如果没有指定 `CMAKE_C_FLAGS`,
  - 将采用默认值: 空字符串。
- 如果没有指定 `CMAKE_SOURCE_DIR`, `CMAKE_CURRENT_LIST_DIR`, `ROOT_DIR`, `__root_dir`,
  - 将采用默认值: 根 `CMakeLists.txt` 所在的文件夹路径。
- 如果没有指定 `CMAKE_CURRENT_SOURCE_DIR`, `COMPONENT_DIR`,
  - 将采用默认值: 当前处理的 `CMakeLists.txt` 所在的文件夹路径。

如果仍然没有找到，则中止处理当前这条配置。

### 输出

- 日志信息输出到: Stderr
- 结果输出到: Stdout
  - 结果格式: `宏名称=宏内容`, 每个一行

#### 示例

假定有配置文件 `ExtractDefine.yaml`, 将运行日志输出到 `run.log`, 将运行结果输出到 `result.txt`, 请使用下面的命令:

`./ExtractDefine -c ExtractDefine.yaml >result.txt 2>run.log`

如果只输出到 `result.txt`, 请使用下面的命令:

`./ExtractDefine -c ExtractDefine.yaml >result.txt`

## 编译

`go build .`

全平台编译参考 [build.bat](build.bat) 。

### 协议

[木兰宽松许可证， 第2版](http://license.coscl.org.cn/MulanPSL2)
