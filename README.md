## XcelFlow
配置表导出工具，目前处于早期阶段，有的功能已经实现

## 支持类型
- [x] 支持多种基础类型 
  - 整形：`int` `int8` `int16` `int32` `int64` 默认值为0
  - 浮点型：`float` `float32` `flat64` 默认值为0
  - 字符串：`str` 默认值为空字符串 
  - 布尔型：`bool` 表示形式`true|TRUE` `false|FALSE`，默认值为false
- [ ] 支持多语言
  - `lang`
- [x] 支持多种结构类型，支持泛型、可以嵌套 
  - 列表：`list` `list(T)`
    - `list` 弱类型列表，列表中的值可以时任何类型 
      - `list` [1, 0.1, true, "中文", (1, 0.1, true, "中文")]
    - `list(T)` 强类型列表，列表中的值必须为指定类型
      - `list(int)` [1, 2, 3]
      - `list(tuple(int))` [(1, 2), (3, 4), (5, 6)]
  - 元组：`tuple` `tuple(T)`
    - tuple 弱类型元组，列表中的值可以时任何类型
      - `tuple` (1, 0.1, true, "中文", [1, 0.1, true, "中文"])
    - `tuple(T)` 强类型元组，列表中的值必须为指定类型
      - `tuple(int)` (1,2,3)
      - `tuple(list(int))` ([1, 2], [3, 4], [5, 6])
  - 字典：`dict` `dict(KT,VT)` 
    - dict 弱类型字典，键是基础类型，值可以是任何类型
      - `dict` {1 = [(1, 2), {2 = 3}, "中文"], abc = def}
    - `dict(KT,VT)` 强类型字典，键和值必须值指定类型
      - `dict(int, list(int))` {1 = [1, 2, 3], 2 = [4, 5, 6]}
- [x] 支持任意类型 
  - `any` 支持填入以上所有类型数据 

## 支持装饰器
用于修饰本列，或全表数据，一列可以指定使用多个装饰器，每个装饰器占一行
- [x] `p_key` 主键装饰器
  - 设置当前列为主键，可以在多列同时设置形成多主键，用于唯一性检查，重复则报错，第一列默认为主键列
  - 用法 `p_key`
- [ ] f_key(表名,列名[,值作用域]) 值引用装饰器
  - 设置当前列为引用列，可以引用其他表中的列数据，用于主外键检查，外键不存在则报错
  - 用法 `f_key(item,id) f_key(item,id, $.e) f_key(item,id, $.e[0]) f_key(item,id, $.k)`
- [x] `u_key` 唯一装饰器 
    - 设置当前列单元格中的值在本列唯一，否则报错
    - 用法 `u_key`
- [x] `not_null` 非空装饰器 
  - 设置当前列非空，检查当前列是否有空值，存在则报错
  - 用法 `not_null`
- [x] `default(默认值)` 默认值装饰器 
  - 设置列默认值，如果此列单元格未填写内容，第一默认值为列类型的默认值，如果需要更改可以通过使用此装饰器修改
  - 用法 `default(1) default([]) default({})`
- [x] `range(最小值,最大值[,值作用域])` 范围装饰器 
  - 设置单元格中的值的取值范围，超出则报错，目前仅可应用与整形与浮点型的检测
  - 用法 `range(1,1000) range(1,1000,$.e) range(1,1000,$.e[0]) range(1,1000,$.k)`
- [x] `macro(值字段列名|列号[,描述字段列名|列号])` 宏定义装饰器 
  - 设置当前列为宏定义名称列，并根据指定的字段列名和描述字段列名进行生成配置
  - 用法 `macro(id, desc) macro(id) macro(1,2)`
- [x] `resource(资源路径[,值作用域])` 资源装饰器 
  - 设置当前列所引用的本地资源路径，检查文件是否存在，否则报错
  - 用法 `resource(samples/icon) resource(samples/icon,$.e) resource(samples/icon,$.e[0]) resource(samples/icon,$.v)`
- [] `ref_table(表名,索引字段名)` 引用表装饰器
  - 设置当前列为引用表，可以引用其他表中的多行数据，具体用法待定
  - 用法 `ref_table(item, id)`
- [] `skip_table([模式名, ...])` 跳过表生成装饰器
  - 设置指定模式的表不生成配置数据，不输出任何内容，为空则跳过任何模式的配表生成
  - 用法 `skip_table([erlang])`

## 以支持导出模板
- [x] Erlang
- [x] JSON
- [x] Flatbuffers
- [x] TypeScript

### TODO
- 建立sqlite数据库
- 支持f_key，引用表表修改检查、被引用修改检查反查
- 支持ref_table、skip_table装饰器
- 支持自定义模板
- 支持多语言
- 创建结构单例化 reader、parser等
- 整理错误码

## 扩展

