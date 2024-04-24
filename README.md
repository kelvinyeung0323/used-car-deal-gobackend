# 二手车交易小程序后端

注：这个自己个人学习而写的项目



## 项目背景：
没有具体的需求，因为学习golang,所以假设给自己一个需求，走一次实际开发过程才好掌据学习的技术，因为后端管理系统是最基础的系统，在这基础上具体的业务需求就可以迅速开展；

项目分三部份:后端程序、前端管理程序、小程序

## 需求
大致对需求进行了梳理
![image](https://github.com/kelvinyeung0323/used-car-deal-gobackend/blob/master/doc/images/micro-app-requirements.png)

![image](https://github.com/kelvinyeung0323/used-car-deal-gobackend/blob/master/doc/images/micro-app-requirement2.png)

## 项目架构
后端程序包括后台管理和交易，因为多年的经验，一般小公司小项目并没有那个体量必须分布式，所以项目以monolith的模式架构，构建一个维护简单的系统；最终程序只分成两部份：系统程序（前端打包进go程序中）+mysql数据库

![image](https://github.com/kelvinyeung0323/used-car-deal-gobackend/blob/master/doc/images/architecture.png)


## 数据库设计
### RABC权限模型
这主要针对后端管理用户，具体业务再另行设计业务用户表；
![image](https://github.com/kelvinyeung0323/used-car-deal-gobackend/blob/master/doc/images/database/user-role-res.png)

## 参数表
用于储存一些参数，如国别、省市区、汽车级别、能源类型等等信息；
![image](https://github.com/kelvinyeung0323/used-car-deal-gobackend/blob/master/doc/images/database/param.png)

### 商品规格表
因为主要针对二手车，所以产品规格主要用于存储各车型的参数配置和图片信息；
![image](https://github.com/kelvinyeung0323/used-car-deal-gobackend/blob/master/doc/images/database/item-spec.png)

### 商品表
用于存储销售的二手车信息，其中检测报告模型时间关系暂未设计
![image](https://github.com/kelvinyeung0323/used-car-deal-gobackend/blob/master/doc/images/database/product.png)
### 订单表及其他
暂未设计



## 技术架构
项目主要使用 go + gin + sqlx 

下面针对以下几项进行说明：
- 项目结构
- 权限验证
- 全局异常处理
- ORM
- 全局事务管理
- Websocket
- 监控线程

### 项目结构
参考java项目习惯将项目分层；  
routes对应java的controller,一些请求参数处理在route里做了，本来想把controller单独出来的，但发现这样做分层太多了，controller做的事本来就不多，所以跟route整合在一起；   
handler对应service层，负责业务逻辑的处理；  
repositories负责数据库操作； 
resources目录存放一些资源文件，使用embbed,打包时嵌入到go程序中；
其中，gin.Context贯穿reoutes->handler->respositories,用于一些全局处理；

### 权限验证
权限验证jwt,编写中间件对请求进行拦截验证权限；这里不多说，网上很多资料；

### 全局异常处理
一般java web应用都会设置一个全局异常处理类来对异常进行统一的处理；
所以项目里使用中间件进行统一异常处理；因为使用error的处理方式个人觉得过于啰唆；  
在repo里使用error的方式处理；因为repo里的逻辑没有太大的业务意义；
然后在handler中处理error,抛panic，这里的panic是有具体内业意义的；   
在中间件中recover panic;然后统一返回错误响应；


### ORM
个人不太喜欢gorm，喜欢java的mybatis的模式；  
本人对SQL比较熟悉，所以gorm的模式太别扭，不如写SQL来的方便；  
所以选择了sqlx+template的方式；   
将SQL写在模板中，通过template编译，这样就跟Mybatis差不多了；

### 全局事务处理
go不像java一样可以使用注解来配置事务处理;     
网上找的解决方案，个人觉得都不太好,过于复杂；在这个项里我把事务处理的逻辑提到中间件里进行处理；  
1.每一个请求到来，在中间件里开启事务，然后把数据库连接放在gin.Context里；    
2.repo层先从gin.Context中把连接合出来再进行操作；    
3.如果数据操作失败或业务逻辑错误，则抛出特定类型的panic;   
4.在中间件进行recover,然后对事务进行rollback;    
5.如果没有panic则commit;    


## Websocket
项目使用websocket跟前端进行实时交互


## 监控线程
监控线程负责定时轮询监控设备，把获取到的数据通过websocket发送到前端并保存相关日志；
通过socket连接设备，每种设备都有各自的交互协议；



