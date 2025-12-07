# **GoKafkaWarehouse**

GoKafkaWarehouse is a real-time ETL pipeline implemented with **Go**, **Debezium**, **Kafka**, and **MySQL** (as both OLTP source and OLAP target).
The system performs **change data capture**, applies **domain-driven transformations**, and loads the processed data into an OLAP database using multiple Go consumers.

---

## **Architecture Overview**

```
                 +------------------+
                 |     MySQL OLTP   |
                 |   (source DB)    |
                 +---------+--------+
                           |
                           | CDC (binlog)
                           v
                 +----------------------+
                 |      Debezium       |
                 |   MySQL Connector   |
                 +----------+-----------+
                            |
                            | Kafka topics
                            v
                    +-----------------+
                    |     Kafka       |
                    +--------+--------+
                             |
       +---------------------+---------------------+
       |                     |                     |
       v                     v                     v
+--------------+     +---------------+     +----------------+
| Consumer A   |     |  Consumer B   |     |  Consumer C    |
| (ingestion)  |     | (transform)   |     | (load to OLAP) |
+------+-------+     +-------+-------+     +--------+-------+
       |                     |                     |
       v                     v                     v
                         Go Services
                             |
                             v
                    +-----------------+
                    |   MySQL OLAP    |
                    |  (analytics DB) |
                    +-----------------+
```

---

## **Project Goals**

* Populate OLAP tables in **real time** from an OLTP MySQL instance.
* Provide a **fully functional ETL pipeline** demonstrating CDC → Kafka → Go services → OLAP loading.
* Showcase domain-driven design with multiple specialized consumer services.

---

## **Technologies**

| Component        | Technology               |
| ---------------- | ------------------------ |
| Language         | Go                       |
| CDC              | Debezium MySQL Connector |
| Messaging Layer  | Apache Kafka (local)     |
| Source Database  | MySQL (OLTP)             |
| Target Warehouse | MySQL (OLAP)             |
| Architecture     | DDD + multiple consumers |

---

## **Domain-Driven Design Structure**

```
/cmd
   /ingestor        → first consumer that reads raw CDC messages
   /transformer     → applies business rules and enrichment
   /loader          → writes transformed records into the OLAP DB
/internal
   /domain          → entities, aggregates, value objects
   /application     → use cases and services
   /infrastructure  → kafka clients, mysql repositories, config
   /pkg             → shared utilities
```

---

## **How to Run (Local Environment)**

### **Prerequisites**

* Docker + Docker Compose
* Go **>= 1.21**
* Make (optional)

### **1. Start infrastructure**

This brings up MySQL (OLTP), MySQL (OLAP), Kafka, Zookeeper, and Debezium.

```bash
docker compose up -d
```

### **2. Run consumers**

Each consumer is an independent Go service.

```bash
go run cmd/ingestor/main.go
go run cmd/transformer/main.go
go run cmd/loader/main.go
```

---

## **Pipeline Flow**

1. **Debezium** captures MySQL binlog events.
2. Events are published into **Kafka topics**.
3. Each Go consumer performs a specific stage of the ETL pipeline:

   * **Ingestor**: reads raw events.
   * **Transformer**: applies domain rules and transformations.
   * **Loader**: persists the final dataset into the OLAP database.
4. OLAP tables are updated **continuously and in real time**.

---

## **Use Case**

This project is designed as a **portfolio demonstration** for senior engineering roles, showcasing:

* real-time data engineering skills
* streaming ETL architecture
* Go microservices
* Kafka-based integrations
* Domain-Driven Design applied to data pipelines

---

## **License**

This project is released under the **MIT License**.
You are free to use, modify, and distribute it.
