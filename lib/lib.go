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
