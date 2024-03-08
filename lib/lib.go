package lib

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func Record(wapi api.WriteAPIBlocking, id string, point int, when time.Time) {
	p := influxdb2.NewPointWithMeasurement("chat").
		AddTag("id", id).
		AddField("content", point).
		SetTime(when)

	if err := wapi.WritePoint(context.Background(), p); err != nil {
		panic(err)
	}
}

func Rank(qapi api.QueryAPI) map[string]int {
	res, err := qapi.Query(
		context.Background(),
		`from(bucket: "hello")
			|> range(start: -6d)
			|> filter(fn: (r) => r["_measurement"] == "chat")
			|> group(columns: ["id"])
			|> count()
			|> group()
			|> sort(columns: ["_value"], desc: true)`,
	)
	if err != nil {
		panic(err)
	}

	rankMap := map[string]int{}

	for res.Next() {
		rankMap[res.Record().ValueByKey("id").(string)] = res.Record().Value().(int)
	}

	if res.Err() != nil {
		panic(res.Err())
	}

	return rankMap
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
