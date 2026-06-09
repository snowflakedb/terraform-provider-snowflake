# Iceberg table test fields

This directory serves as preloaded data for iceberg tables (on files and delta types). This directory cannot be uploaded to the external stage in test setup because of the following limitation:
```
091003 (22000): Failure using stage area. Cause: [GMDFOYIT_CDEECDC2_69C2_E6B3_0B0C_B6B39DFA61C8 GET and PUT commands are not supported with external stage]
```
This is why they are uploaded outside of the test workflow. You can use `aws` cli to view and upload the needed files. The structure on the S3 stage mimics the directory structure.

## Delta format creation
```python
import pandas as pd
from deltalake.writer import write_deltalake

df = pd.DataFrame({
    'id': [1, 2, 3],
    'name': ['Alice', 'Bob', 'Carol'],
    'amount': [100.0, 200.0, 150.0],
})

write_deltalake('/tmp/sample_delta_table', df)
```

## Iceberg format creation

```python
  from pyiceberg.catalog import load_catalog
  from pyiceberg.schema import Schema
  from pyiceberg.types import NestedField, LongType, StringType

  catalog = load_catalog("local", **{"type": "sql", "uri": "sqlite:////tmp/iceberg.db", "warehouse": "/tmp/iceberg_warehouse"})

  catalog.create_namespace("ns")

  schema = Schema(
      NestedField(1, "id", LongType(), required=False),
      NestedField(2, "name", StringType(), required=False),
  )

  table = catalog.create_table("ns.iceberg_test_table", schema=schema)
  print(table.metadata_location)
```
