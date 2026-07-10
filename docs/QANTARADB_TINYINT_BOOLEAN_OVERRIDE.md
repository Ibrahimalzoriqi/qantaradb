# Qualified Tinyint Boolean Override

`tinyint(1)` is often a MySQL boolean, but legacy status fields can use numeric
states beyond `0` and `1`. Mapping such a column to PostgreSQL `boolean` loses
meaning and causes QantaraDB parity failure.

Keep the backward-compatible default and opt out only with a qualified table
and column name:

```yaml
mapper:
  tinyint1_as_bool: true
  tinyint1_boolean_denylist:
    - offline_payments.status
```

The denylist is exact and case-insensitive. It only affects the named
`table.column`; other `tinyint(1)` columns retain the configured boolean
mapping. An opted-out column maps to PostgreSQL `smallint`, so the loader does
not coerce its numeric values to `true` or `false`.

Local QantaraDB proof against the FoodTech MySQL reference completed schema and
data copy, foreign-key validation, and sequence reset with no blocking reasons.
`offline_payments` became a normalized-equivalent validation result after the
override. No source rows or credentials are recorded in this document.
