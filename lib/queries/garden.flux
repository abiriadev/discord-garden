import "experimental/date/boundaries"
import "timezone"

option location = timezone.location(name: "Asia/Seoul")

from(bucket: "{{.Bucket}}")
	|> range(start: {{.Start}})
	|> filter(fn: (r) => r["_measurement"] == "{{.Measurement}}" and r.id == "{{.Id}}")
	|> aggregateWindow(
		every: {{.Window}},
		fn: (column, tables=<-) =>
			tables |> sum(column: "_value"),
		createEmpty: true
	)
