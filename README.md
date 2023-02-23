## Features

- 根据字段类型自动生成数据
- 快速批量插入
- 通过 `schmea` 自定义表结构
- 支持库拷贝，实现快速生成多个库
- 内置多种 `schema`，方便快速插入

## Support

- [MySQL](#MySQL)
- [PostgreSQL](#PostgreSQL)
- [Redis](#Redis)

## Install

```shell
go install github.com/vczyh/dbinsert@latest
```

## Schema

`schema` 定义了表结构，不同的数据库会有不同，但对于关系数据库大体相同，可以使用 `--schema` 指定内置或者自定义 `schema`
，如果不指定则使用默认 `schema`。

支持 `schema` 的数据库：

- MySQL
- PostgreSQL

## MySQL

### MySQL Schema

```json5
[
  {
    // 数据库名称
    "database": "dbinsert",
    // 表名称
    "table": "tbl",
    // 插入 10000 条数据
    "size": 10000,
    // 主键
    "primaryKeyFieldNames": [
      "id"
    ],
    // 所有列
    "fields": [
      {
        "name": "id",
        "type": "INT",
        // 是否自增
        "autoIncrement": true
      },
      {
        "name": "name",
        "type": "CHAR(60)"
      }
    ]
  },
  // ... 
]
```

### 内置 schema

- [`sysbench`](./relation/schema/sysbench_mysql.json)

### Usage

```shell
dbinsert mysql \
  --host x.x.x.x \ 
  --port 3306 \
  --user xxx  \
  --password xxx \
  --create-databases \
  --create-tables
```

### Flags

| 名称                   | 默认                                                  | 说明                                         |
|----------------------|-----------------------------------------------------|--------------------------------------------|
| `--schema`           | [`sysbench`](./relation/schema/sysbench_mysql.json) | 内置或者自定义 `schema`                           |
| `--host`             | `127.0.0.1`                                         | 域名或IP                                      |
| `--port`             | `3306`                                              | 端口                                         |
| `--username`         | `root`                                              | 用户                                         |
| `--password`         | `""`                                                | 密码                                         |
| `--create-databases` | `false`                                             | 创建库如果不存在                                   |
| `--create-tables`    | `false`                                             | 创建表如果不存在                                   |
| `--table-size`       | `0`                                                 | 表记录条数，覆盖 `Schema` 中 `size`，`0` 不生效         |
| `--timeout`          | `10h`                                               | 超时时间，`3m` 表示3分钟结束运行                        |
| `--db-repeat`        | `0`                                                 | 数据库重复次数，会生成多个库 [`dbinsert_1`，`dbinsert_2`] |

## PostgreSQL

### PostgreSQL schema

```json5
[
  {
    // 数据库名称
    "database": "dbinsert",
    // 表名称
    "table": "tbl",
    // 插入 10000 条数据
    "size": 10000,
    // 主键
    "primaryKeyFieldNames": [
      "id"
    ],
    // 所有列
    "fields": [
      {
        "name": "id",
        "type": "serial"
      },
      {
        "name": "name",
        "type": "CHAR(60)"
      }
    ]
  }
]
```

### 内置 Schema

- [`sysbench`](./relation/schema/sysbench_postgres.json)

### Usage

```shell
dbinsert postgres \
  --host 100.100.1.194 \ 
  --port 3306 \
  --user xxx \
  --password xxx \
  --create-databases \
  --create-tables
```

### Flags

| 名称                   | 默认                                                     | 说明                                         |
|----------------------|--------------------------------------------------------|--------------------------------------------|
| `--schema`           | [`sysbench`](./relation/schema/sysbench_postgres.json) | 内置或者自定义 `schema`                           |
| `--host`             | `127.0.0.1`                                            | 域名或IP                                      |
| `--port`             | `5432`                                                 | 端口                                         |
| `--username`         | `""`                                                   | 用户                                         |
| `--password`         | `""`                                                   | 密码                                         |
| `--create-databases` | `false`                                                | 创建库如果不存在                                   |
| `--create-tables`    | `false`                                                | 创建表如果不存在                                   |
| `--table-size`       | `0`                                                    | 表记录条数，覆盖 `Schema` 中 `size`，`0` 不生效         |
| `--timeout`          | `10h`                                                  | 超时时间，`3m` 表示3分钟结束运行                        |
| `--db-repeat`        | `0`                                                    | 数据库重复次数，会生成多个库 [`dbinsert_1`，`dbinsert_2`] |

## Redis

### Usage

```shell
dbinsert redis \
  --host 127.0.0.1 \
  --port 6379 \
  --user xxx \
  --password xxx \  
  --key-count 10000 
```

支持集群：

```shell
dbinsert redis \
  --cluster \
  --addrs 100.100.5.222:6379 \
  --user xxx \
  --password xxx  \
  --key-count 30000
```

### Flags

| 名称            | 默认               | 支持类型  | 说明                  |
|---------------|------------------|-------|---------------------|
| `--user`      | `default`        | 单点和集群 | 用户                  |
| `--password`  | `""`             | 单点和集群 | 密码                  |
| `--timeout`   | `10h`            | 单点和集群 | 超时时间，`3m` 表示3分钟结束运行 |
| `--key-count` | `0`              | 单点和集群 | 插入的记录数, `0` 代表无限    |
| `--value-len` | `50`             | 单点和集群 | `value` 字符串的长度      |
| `--host`      | `127.0.0.1`      | 单点    | 域名或IP               |
| `--port`      | `6379`           | 单点    | 端口                  |
| `--cluster`   | `false`          | 集群    | 开启集群模式              |
| `--addrs`     | `127.0.0.1:6379` | 集群    | 逗号分隔的集群地址列表         |