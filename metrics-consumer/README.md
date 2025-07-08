
# Homework 9: Kafka Metrics with Partitioned Events and Replication

## 📝 Goal
Build a small distributed system with:

- `cart-service` and `stock-service` producing structured metrics events to Kafka.
- `metrics-consumer` that consumes and logs these events.
  - Have its own Dockerfile.
  - Be built into a Docker image.
  - Be pushed to Docker Hub as a public image, so that your docker-compose.yml can pull it by tag (e.g. yourname/metrics-consumer:hw9).
- Kafka topic with **2 partitions and replication factor 2** running on **2 Kafka brokers**.
- Using Docker Compose for Kafka (in a **separate Compose file with shared network**).

---

## 📂 Project structure

```
.
├── kafka/
│   └── docker-compose.yml   
│     # Kafka + Zookeeper cluster (2 brokers) + kafka-ui
│   └── Dockerfile   
├── cart/
│
├── stock/
│   
├── metrics-consumer/
├── go.work                 # Go workspace to link all modules
├── README.md
```

---

## ⚙️ Kafka setup
- Located under `kafka/docker-compose.yml`.
- Runs:
  - 2 Kafka brokers
  - 1 Zookeeper
  - kafka-ui
- Creates topic `metrics` with:
  - **2 partitions** (partition 0 for cart events, 1 for stock events)
  - **replication.factor=2**

### 📌 Shared network
Kafka Compose uses:
```yaml
networks:
  shared-net:
    external: true
```
so your services can join this `shared-net` and reach Kafka by:
```
kafka1:9092, kafka2:9092
```

---

## 🚀 Services

| Service           | Description                                    | Writes to Kafka        |
|-------------------|------------------------------------------------|-------------------------|
| `cart-service`    | Simulates adding items to cart                 | Always writes to **partition 0** |
| `stock-service`   | Simulates SKU creation & stock changes         | Always writes to **partition 1** |
| `metrics-consumer`| Subscribes to `metrics` topic, logs all events | Reads both partitions |

All services run in the same `shared-net`.

---

## 📚 Event structure
All Kafka messages must be JSON in this format:

| Field       | Type   | Description                           |
|-------------|--------|---------------------------------------|
| `type`      | string | Event type, e.g. `cart_item_added`    |
| `service`   | string | `"cart"` or `"stock"`                 |
| `timestamp` | string | ISO8601 UTC timestamp                 |
| `payload`   | object | Event-specific data                  |

### 🛒 Cart events
#### `cart_item_added`
```json
{
  "type": "cart_item_added",
  "service": "cart",
  "timestamp": "2025-07-08T19:15:32Z",
  "payload": {
    "cartId": "xyz123",
    "sku": "A123",
    "count": 2,
    "status": "success"
  }
}
```
#### `cart_item_failed`
```json
{
  "type": "cart_item_failed",
  "service": "cart",
  "timestamp": "2025-07-08T19:16:04Z",
  "payload": {
    "cartId": "xyz123",
    "sku": "A123",
    "count": 5,
    "status": "failed",
    "reason": "not enough stock"
  }
}
```

### 📦 Stock events
#### `sku_created`
```json
{
  "type": "sku_created",
  "service": "stock",
  "timestamp": "2025-07-08T19:20:17Z",
  "payload": {
    "sku": "A123",
    "price": 12.5,
    "count": 100
  }
}
```
#### `stock_changed`
```json
{
  "type": "stock_changed",
  "service": "stock",
  "timestamp": "2025-07-08T19:21:50Z",
  "payload": {
    "sku": "A123",
    "count": 12,
    "price": 12.5
  }
}
```

---

## 🚀 How to run
1️⃣ **Start Kafka cluster**
```bash
cd kafka
docker-compose up -d
```
Creates brokers, zookeeper, topic `metrics` with 2 partitions & replication=2.

2️⃣ **Start services**
Each service’s Dockerfile joins the same `shared-net`, connects to Kafka by:
```
KAFKA_BROKERS=kafka1:9092,kafka2:9092
TOPIC=metrics
```
```bash
docker-compose up -d cart stock metrics-consumer
```
(Or use `docker run` manually with `--network shared-net`).

---

## 🔍 How to test replication
- After all services are producing & consuming:
```bash
docker stop kafka1
```
- You should still see metrics flowing thanks to replication on `kafka2`.

---

## 📑 What to submit
✅ Source code for:
- `cart-service`
- `stock-service`
- `metrics-consumer`
- `kafka/docker-compose.yml` for brokers + topic creation

✅ A short demo log showing:
- `cart` writes to partition 0, `stock` writes to partition 1
- `metrics-consumer` prints all events

✅**Note**: Please make sure to update or fix your existing tests after adding the Kafka integration.

---


🎉 **Good luck!**  
