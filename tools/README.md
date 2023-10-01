# Tools

## 关于 common

common.BaseTool 包含在一个 "工具结构体" 的 Base 字段

## 声明一个工具

除 common 文件夹外，其中的每一个文件夹都被当作一个"工具"，工具的名称就是文件夹名称。

以 screencap 举例。

screencap/base.go 用于实现 common.BaseTool 接口。

screencap/screencap.go 用于实现主要功能。

-   其中 IScreencap 接口是此工具自有的一些函数，与工具的具体功能相关联。
-   ScreencapImpl 结构体是实现功能的结构体。Impl 结构体必须包含一个 Base common.BaseTool 字段，并且实现对应的 Ixxx 接口。如果需要注册命令行子命令便实现 common. UseCommand 接口。
-   需要创建一个 Impl 实例作为全局入口 (var Screencap \*ScreencapImpl = New())
