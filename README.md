# 路径管理
在所有endpoints通过握手流程完成了协商后，且所有endpoints的 multipath特性开启，endpoints可以使用multipath。
该建议为路径管理增加了一个frame
- PATH_STATUS 帧，用于接收端告知path的状态和优先级
  
所有的MP帧都是在1-RTT报文里面发送

## Path Identifier和Connection ID


## 数据流控制

QUIC提供了两种层面上的数据流控制方案：

- Stream 流量控制，通过限制在任何 stream 上可以发送的最大绝对字节偏移量，防止单个 stream 消耗连接（quic connection）的全部接收缓冲。
- Connection流量控制，通过限制所有STREAM帧的数据总字节数，防止发送方超过接收方的连接缓冲容量。