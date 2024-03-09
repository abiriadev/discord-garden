package lib

import (
	"context"
	"strings"
	"text/template"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func Record(wapi api.WriteAPIBlocking, id string, point int, when time.Time) {
	p := influxdb2.NewPointWithMeasurement("chat").
		AddTag("id", id).
		AddField("point", point).
		SetTime(when)

	if err := wapi.WritePoint(context.Background(), p); err != nil {
		panic(err)
	}
}

type RankRecord struct {
	Id    string
	Point int
}

func Rank(qapi api.QueryAPI) []RankRecord {
	var buf strings.Builder

	if tmpl, err := template.New("rank").Parse(
		`from(bucket: {{.Bucket}})
			|> range(start: 0)
			|> filter(fn: (r) => r["_measurement"] == {{.Measurement}})
			|> group(columns: ["id"])
			|> sum(column: "_value")
			|> group()
			|> sort(columns: ["_value"], desc: true)
			|> limit(n: {{.Limit}})`,
	); err != nil {
		panic(err)
	} else if err := tmpl.Execute(&buf, struct {
		Bucket      string
		Measurement string
		Limit       int
	}{
		"hello",
		"chat",
		10,
	}); err != nil {
		panic(err)
	}

	res, err := qapi.Query(
		context.Background(),
		buf.String(),
	)
	if err != nil {
		panic(err)
	}

	rankMap := []RankRecord{}

	for res.Next() {
		id, ok := res.Record().ValueByKey("id").(string)
		if !ok {
			id = "anon"
		}

		rankMap = append(rankMap, RankRecord{
			Id:    id,
			Point: int(res.Record().Value().(int64)),
		})
	}

	if res.Err() != nil {
		panic(res.Err())
	}

	return rankMap
}

func Garden(qapi api.QueryAPI) []int {
	res, err := qapi.Query(context.Background(),
		`from(bucket: "hello")
			|> range(start: -1h)
			|> filter(fn: (r) => r["_measurement"] == "chat" and r.id == "662201438621138954")
			|> aggregateWindow(
				every: 1m,
				fn: (column, tables=<-) =>
					tables |> sum(column: "_value"),
				createEmpty: true
			)`,
	)
	if err != nil {
		panic(err)
	}

	gardenMap := []int{}

	for res.Next() {
		v, _ := res.Record().Value().(int64)

		gardenMap = append(gardenMap, int(v))
	}

	return gardenMap
}

func InitClient(
	addr string,
	token string,
	org string,
	bucket string,
) (api.QueryAPI, api.WriteAPIBlocking) {
	client := influxdb2.NewClient(
		addr,
		token,
	)

	return client.QueryAPI(org), client.WriteAPIBlocking(org, bucket)
}
