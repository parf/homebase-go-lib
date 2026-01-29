# Test Results: ANY Schema Support

## ✅ any2parquet - Supports ANY Schema

### Test 1: Product Schema
```jsonl
{"product":"Laptop","price":999.99,"stock":50}
{"product":"Mouse","price":29.99,"stock":100}
```

**Result:** ✅ Success
- Input: 2 records (product, price, stock)
- Output: 683 bytes Parquet
- Schema auto-inferred: price (float64), product (string), stock (int64)

### Test 2: Sensor Schema (Completely Different!)
```jsonl
{"sensor":"temp-01","value":23.5,"unit":"celsius","online":true}
{"sensor":"pressure-02","value":1013.25,"unit":"hPa","online":false}
```

**Result:** ✅ Success
- Input: 2 records (sensor, value, unit, online)
- Output: 768 bytes Parquet
- Schema auto-inferred: online (bool), sensor (string), unit (string), value (float64)

## ✅ any2jsonl - Supports ANY Schema

### Round-trip Test 1: Product Schema
```bash
any2parquet test-schema1.jsonl test-schema1.parquet
any2jsonl test-schema1.parquet test-schema1-out.jsonl
```

**Result:** ✅ Perfect round-trip
```jsonl
{"price":999.99,"product":"Laptop","stock":50}
{"price":29.99,"product":"Mouse","stock":100}
```

### Round-trip Test 2: Sensor Schema
```bash
any2parquet test-schema2.jsonl test-schema2.parquet
any2jsonl test-schema2.parquet test-schema2-out.jsonl
```

**Result:** ✅ Perfect round-trip
```jsonl
{"online":true,"sensor":"temp-01","unit":"celsius","value":23.5}
{"online":false,"sensor":"pressure-02","unit":"hPa","value":1013.25}
```

## ⚠️ any2fb - Fixed Schema ONLY

FlatBuffer requires generated code for each schema, so it cannot support arbitrary schemas.

**Supported Schema (FIXED):**
- id, name, email, age, score, active, category, timestamp

**Recommendation:** Use `any2parquet` for ANY schema needs!

## Summary

| Tool | ANY Schema Support | Notes |
|------|-------------------|-------|
| **any2parquet** | ✅ YES | Auto-infers schema from data |
| **any2jsonl** | ✅ YES | Reads/writes any structure |
| **any2fb** | ❌ NO | FlatBuffer limitation, use Parquet instead |

## Tested Schemas

✅ Product catalog: product, price, stock
✅ Sensor data: sensor, value, unit, online
✅ E-commerce: id, product, category, price, in_stock, quantity, rating
✅ IoT mixed: sensor_id, location, temperature, humidity, detections, battery
✅ Users: user_id, username, email, age, premium, credits, country

**All work perfectly with any2parquet and any2jsonl!**
