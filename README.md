# xlib
游戏服务器 lib

## 项目初始化
- go mod init github.com/75912001/xlib
- 清理依赖
  - go mod tidy
- 检查依赖
  - go mod verify

## git 提交标签
- \[feat\]: 新功能（feature）
- \[fix\]: 修复问题（bug fix）
- \[docs\]: 文档的变更
- \[style\]: 代码样式的变更（不影响代码运行的变动）
- \[refactor\]: 重构代码
- \[test\]: 添加或修改测试
- \[chore\]: 构建过程或辅助工具的变动


- **lib**: 公共库
  - **bench**: 服务基础配置[todo menglc]
  - **common**: 公共模块[todo menglc]
  - **constants**: 常量[todo menglc]
  - **control**: 控件
  - **error**: 错误码[todo menglc]
  - **etcd**: etcd客户端[todo menglc]
  - **example**: 示例[todo menglc]
  - **exec**: 执行器 [todo menglc]
  - **file**: 文件操作[todo menglc]
  - **log**: 日志[todo menglc]
  - **message**: 消息[todo menglc]
  - **net**: 网络
  - **packet**: 数据包[todo menglc]
  - **pool**: 对象池[todo menglc]
  - **pprof**: 性能分析[todo menglc]
  - **runtime**: 运行时[todo menglc]
  - **server**: 服务[todo menglc]
  - **subpub**: 订阅发布[todo menglc]
  - **time**: 时间管理器[todo menglc]
  - **timer**: 定时器[todo menglc]
  - **util**: 工具类[todo menglc]