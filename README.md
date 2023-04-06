# gateway
基于go设计一款简易的网关服务


## 网关
功能点：  
1.反向代理 ✅    
2.负载均衡(支持随机，轮询，加权轮询，一致性hash) ✅  
3.header头转换  
4.strip_uri  
5.url重写  
6.ip白名单和黑名单控制  
7.流量统计  
8.漏桶限流控制  
9.jwt认证  

## 参考
1. [load balance简易逻辑](https://kasvith.me/posts/lets-create-a-simple-lb-go/)
2. [Bilibili 简易负载均衡开发](https://www.bilibili.com/video/BV1P7411B7W2/)