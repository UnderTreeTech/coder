## XO模板代码生成

0、XO生成的模板只针对标准SQL，理论上不使用各数据库独特语法如:replace、limit，
都可适用于基于标准SQL的数据。包括MySQL、MSSQL、Oracle、PostgreSQL、SQLite，不支持NOSQL

1、依赖安装
- 如未安装goimports则执行以下命令

`go get -u golang.org/x/tools/cmd/goimports`

- 安装XO

**xo项目已更换项目名，目前叫dbtpl。当前的xo依赖特定提交，安装命令如下。可以install后自行rename为xo。**

`go install github.com/xo/dbtpl@9a3ddc1e1407243ea7d30a3956073fcbeaa8d7bc`

2、生成代码操作步骤：
可将所有基于标准SQL的表都在MySQL里创建一下，不论是kingbase、达梦还是其他数据库

- 本机新建目录model/dao，model/iface，model/model
  
- 生成model

`xo mysql://root@127.0.0.1:3306/db_global -o ./model --template-path /Users/sunqiang1/export/apps/golang/360/conversation/templates/model`

- 生成dao实现

`xo mysql://root@127.0.0.1:3306/db_global -o ./dao --template-path /Users/sunqiang1/export/apps/golang/360/conversation/templates/dao`

- 生成iface定义

`xo mysql://root@127.0.0.1:3306/db_global -o ./iface --template-path /Users/sunqiang1/export/apps/golang/360/conversation/templates/iface`

- 将生成的文件copy至项目对应文件夹

## 注意事项

- 关于XO使用

模板只定义了常用的函数，基本能满足80%的场景。
超出场景需要定制的情况，文件命名与xo保持一致，去掉xo即可。不可混在一起，xo文件可
反复生成，混在一起将被覆盖自定义函数
  
- 关于DAO操作规范

禁止联表查询，复杂查询必须拆分成单条简单查询，数据库可以快速返回无压力。
数据聚合的压力应由业务服务器来承担，分摊算力。


