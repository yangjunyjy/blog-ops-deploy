# 五、kafka配置消息压缩

## 1.更改原有主题为压缩格式

```bash
/data/oneops1/kafka/kafka_2.13-3.5.0/bin/kafka-configs.sh --bootstrap-server 134.95.172.17:9092 --entity-type topics --entity-name test --alter --add-config compression.type=zstd
```

## 2.更改压缩等级

```bash
# 编辑每个 Broker 的 server.properties 文件
vim /data/oneops1/kafka/kafka_2.13-3.5.0/config/server.properties

# 设置压缩类型
compression.type=zstd

# 设置 zstd 压缩级别（1-22，推荐 3-6）
compression.zstd.level=5

# 重启kafka
./bin/kafka-server-start.sh -daemon   ./config/kraft/server.properties
```

## 3.新增主题采用压缩格式

```shell
/data/oneops1/kafka/kafka_2.13-3.5.0/bin/kafka-topics.sh --bootstrap-server 134.95.172.17:9092 --create --topic compressed-topic-test --config compression.type=zstd 
```

## 4.验证主题是否使用压缩格式

```bash
/data/oneops1/kafka/kafka_2.13-3.5.0/bin/kafka-configs.sh --bootstrap-server 134.95.172.17:9092 --entity-type topics --entity-name compressed-topic-test --describe
```

