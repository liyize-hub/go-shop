# 基于Iris框架开发的本地购物系统
## 毕业设计
## 项目结构
| 目录或文件 | 说明   |  
|:---------|:-------:|
| config         |  配置|
| controllers    |  控制器|
| log            |  日志文件|
| dao            |  数据库操作|
| datasource     |  数据库配置|
| models         |  数据模型|
| router         |  路由层|
| services       |  服务层|
| utils          |  实用工具|
| main.go  | go主程序入口|


* config：项目配置文件及读取配置文件的相关功能
* controller：控制器目录、项目各个模块的控制器及业务逻辑处理的所在目录
* datasource：实现mysql连接和操作、封装操作mysql数据库的目录。
* model：数据实体目录，主要是项目中各业务模块的实体对象的定义
* service：服务层目录。用于各个模块的基础功能接口定义及实现，是各个模块的数据层。
* static：配置项目的静态资源目录。
* util：提供通用的方法封装。
* main.go：项目程序主入口
* config.json：项目配置文件。

controller层与通信协议耦合，与api框架耦合，api未必只有一套：http, api, grpc, graphql, websocket...
service层实现具体的业务代码
dao层完成对数据库的操作 数据库未必只有一种，mysql, redis...