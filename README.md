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
  - **config**: 服务基础配置
  - **common**: 公共模块
  - **constants**: 常量
  - **control**: 控件
  - **error**: 错误码
  - **etcd**: etcd客户端[todo menglc]
  - **example**: 示例[todo menglc]
  - **exec**: 执行器 [todo menglc]
  - **file**: 文件操作[todo menglc]
  - **log**: 日志[todo menglc]
  - **message**: 消息
  - **net**: 网络
  - **packet**: 数据包
  - **pool**: 对象池[todo menglc]
  - **pprof**: 性能分析[todo menglc]
  - **runtime**: 运行时[todo menglc]
  - **server**: 服务[todo menglc]
  - **subpub**: 订阅发布[todo menglc]
  - **time**: 时间管理器[todo menglc]
  - **timer**: 定时器[todo menglc]
  - **util**: 工具类[todo menglc]





| 符号   | 名称/含义         | 用途说明                                   |
|--------|-------------------|--------------------------------------------|
| ⚠️     | 警告（Warning）    | 用于警告开发者注意某些风险或潜在问题         |
| ❗     | 感叹号             | 强调重要信息或警告                         |
| ❕     | 白色感叹号         | 一般提示                                   |
| ❌     | 叉叉（Cross Mark） | 明确禁止某些用法或操作                     |
| ⛔     | 禁止通行           | 表示禁止、不可用                           |
| 🚫     | 禁止（Prohibited） | 表示禁止、不可用                           |
| ✅     | 绿色对勾           | 表示操作成功或推荐的做法                   |
| ✔️     | 对勾               | 表示完成、通过                             |
| ☑️     | 复选框对勾         | 表示选中、确认                             |
| ℹ️     | 信息（Info）        | 提供额外的信息或说明                       |
| 💡     | 灯泡（Idea/Tip）   | 表示提示、建议、灵感                       |
| 🔥     | 火焰（Hot）        | 标记热点、重要、紧急、热修                 |
| ⭐     | 星星（Star）       | 推荐、重点                                 |
| 📝     | 记事本（Note）     | 备注、说明                                 |
| 🛑     | 停止（Stop）       | 停止、终止                                 |
| ➡️     | 右箭头             | 指引、流程                                 |
| 🟢     | 绿色圆点           | 正常、通过                                 |
| 🟡     | 黄色圆点           | 警告、注意                                 |
| 🔴     | 红色圆点           | 错误、异常                                 |
| ▬▬▬▬  | 分隔线             | 分组、分隔                                 |
| 🚀     | 火箭（Rocket）     | 新特性、上线、发布                         |
| 🧩     | 拼图（Puzzle）     | 模块、插件                                 |
| 🕒     | 时钟（Clock）      | 时间相关                                   |
| 🧨     | 爆炸（Boom）       | 危险、大变动                               |