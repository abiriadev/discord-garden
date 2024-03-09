package lib

import (
	"context"
	"fmt"
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
	res, err := qapi.Query(
		context.Background(),
		`from(bucket: "hello")
			|> range(start: -6d)
			|> filter(fn: (r) => r["_measurement"] == "chat")
			|> group(columns: ["id"])
			|> count()
			|> group()
			|> sort(columns: ["_value"], desc: true)
			|> limit(n: 10)`,
	)
	if err != nil {
		panic(err)
	}

	rankMap := []RankRecord{}

	for res.Next() {
		var id string

		switch v := res.Record().ValueByKey("id").(type) {
		case string:
			id = v
		case nil:
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

func Garden(qapi api.QueryAPI) {
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

	for res.Next() {
		if res.TableChanged() {
			fmt.Printf("table: %s\n", res.TableMetadata().String())
		}

		// switch v := res.Record().ValueByKey("id").(type) {
		fmt.Printf("value: %v\n", res.Record().Value())
		// rankMap = append(rankMap, RankRecord{
		// 	Id:    id,
		// 	Point: int(res.Record().Value().(int64)),
		// })
	}
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
