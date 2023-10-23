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




