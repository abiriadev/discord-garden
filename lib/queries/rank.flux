from(bucket: "{{.Bucket}}")
	|> range(start: {{.Range}})
	|> filter(fn: (r) => r["_measurement"] == "{{.Measurement}}")
	|> group(columns: ["id"])
	|> sum(column: "_value")
	|> group()
	|> sort(columns: ["_value"], desc: true)
	|> limit(n: {{.Limit}})
