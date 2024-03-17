import "experimental/date/boundaries"
import "timezone"

option location = timezone.location(name: "Asia/Seoul")

b = {{.Range}}

from(bucket: "{{.Bucket}}")
	|> range(start: b.start, stop: b.stop)
	|> filter(fn: (r) => r["_measurement"] == "{{.Measurement}}")
	|> group(columns: ["id"])
	|> sum(column: "_value")
	|> group()
	|> sort(columns: ["_value"], desc: true)
	|> limit(n: {{.Limit}})
