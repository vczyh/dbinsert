## Features

- 快速批量插入
- 通过 `Schmea` 自定义表结构
- 支持库、表拷贝，实现快速生成多个库和表
- 内置默认模板，方便快速插入

## Support

- [MySQL](#MySQL)
- [PostgreSQL](#PostgreSQL)

## Schema

`Schema` 定义了表结构，不同的数据库会有不同，但对于关系数据库大体相同。

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

### Usage

```shell
dbinsert mysql \
  --host 100.100.1.194 \ 
  --port 3306 \
  --username cloudos \
  --password Zggyy2019! \
  --create-databases \
  --create-tables
```

### Flags

| 名称  | 说明                                         |
|-----|--------------------------------------------|
| `--host`    | 域名或IP                                      |
|  `--port`   | 端口                                         |
|  `--username`   | 用户                                         |
|     `--password`          | 密码                                         |
|   `--create-databases`            | 创建库如果不存在                                   |
|           `--create-tables`                      | 创建表如果不存在                                   |
|      `--table-size`         | 表记录条数，优先级高于 `Schema` 中 `size`              |
|    `--timeout`                       | 超时时间，`3m` 表示3分钟结束运行                        |
|         `--db-repeat`                            | 数据库重复次数，会生成多个库 [`dbinsert_1`，`dbinsert_2`] |

## PostgreSQL

### PostgreSQL Schema

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
### Usage

### Usage

```shell
dbinsert postgres \
  --host 100.100.1.194 \ 
  --port 3306 \
  --username cloudos \
  --password Zggyy2019! \
  --create-databases \
  --create-tables
```

### Flags

| 名称  | 说明                                         |
|-----|--------------------------------------------|
| `--host`    | 域名或IP                                      |
|  `--port`   | 端口                                         |
|  `--username`   | 用户                                         |
|     `--password`          | 密码                                         |
|   `--create-databases`            | 创建库如果不存在                                   |
|           `--create-tables`                      | 创建表如果不存在                                   |
|      `--table-size`         | 表记录条数，优先级高于 `Schema` 中 `size`              |
|    `--timeout`                       | 超时时间，`3m` 表示3分钟结束运行                        |
|         `--db-repeat`                            | 数据库重复次数，会生成多个库 [`dbinsert_1`，`dbinsert_2`] |

