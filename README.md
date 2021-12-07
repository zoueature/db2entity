# db2entity
同步数据库表到strcut实体


### Usage 
```shell
./main --help
同步数据库表到go实体

Usage:
  db2entity [flags]

Flags:
      --config string        config file (default is $HOME/.db2struct.yaml)
  -d, --database string      数据库名
      --destination string   生成是实体目录 (default ".")
      --help                 usage help
  -h, --host string          数据库host (default "127.0.0.1")
      --package string       包名 (default "entity")
  -p, --password string      数据库密码
  -P, --port string          数据库端口号 (default "3306")
      --prefix string        需要去除的表前缀
  -u, --username string      数据库用户名 (default "root")

```