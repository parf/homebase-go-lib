# Example Schemas for any2parquet / any2jsonl

This directory contains example data files demonstrating that **any2parquet** and **any2jsonl** work with ANY schema structure.

## 📋 Available Examples

### 🛒 E-commerce: `products.jsonl`
Product catalog with pricing and inventory tracking.

**Schema:** `product`, `price`, `stock`, `category`

```bash
../../any2parquet products.jsonl products.parquet
../../any2jsonl products.parquet products-out.jsonl
```

### 🌡️ IoT Sensors: `sensors.jsonl`
IoT sensor readings with different measurement types.

**Schema:** `sensor`, `value`, `unit`, `online`, `location`

```bash
../../any2parquet sensors.jsonl sensors.parquet
../../any2jsonl sensors.parquet sensors-out.jsonl
```

### 👥 Users: `users.csv`
User accounts with demographics and subscription info.

**Schema:** `user_id`, `username`, `email`, `age`, `premium`, `credits`, `country`

```bash
../../any2parquet users.csv users.parquet
../../any2jsonl users.parquet users-out.jsonl
```

### 📝 Application Logs: `logs.jsonl`
Application event logs with severity levels.

**Schema:** `timestamp`, `level`, `message`, `user_id`, `service`

```bash
../../any2parquet logs.jsonl logs.parquet
../../any2jsonl logs.parquet logs-out.jsonl
```

### 💳 Transactions: `transactions.jsonl`
Financial transaction records.

**Schema:** `txn_id`, `amount`, `currency`, `status`, `merchant`

```bash
../../any2parquet transactions.jsonl transactions.parquet
../../any2jsonl transactions.parquet transactions-out.jsonl
```

## 🧪 Run All Tests

Test all schemas to verify converters work with ANY structure:

```bash
./test-all-schemas.sh
```

## ✨ Key Features Demonstrated

- ✅ **Universal Schema Support** - No predefined schema needed
- ✅ **Automatic Type Detection** - int64, float64, string, bool
- ✅ **CSV Support** - With header row type inference
- ✅ **Round-trip Conversion** - Data preserved perfectly
- ✅ **Multiple Formats** - JSONL, CSV input support

## 📊 Data Types Supported

| Type | Example | Detection |
|------|---------|-----------|
| Integer | `42`, `1001` | Parsed as int64 |
| Float | `99.99`, `23.5` | Parsed as float64 |
| String | `"text"`, `"Alice"` | Any text value |
| Boolean | `true`, `false` | Parsed as bool |

## 🚀 Use Cases

These examples represent common real-world data patterns:
- **E-commerce:** Product catalogs, inventory management
- **IoT:** Sensor data, telemetry streams
- **User Data:** Customer databases, CRM systems
- **Logs:** Application monitoring, debugging
- **Finance:** Transaction processing, payment systems

All work seamlessly with any2parquet and any2jsonl! 🎉
