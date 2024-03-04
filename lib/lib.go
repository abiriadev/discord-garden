package main

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func Record(wapi api.WriteAPIBlocking, id string, content string) {
	p := influxdb2.NewPointWithMeasurement("chat").
		AddTag("id", id).
		AddField("content", content).
		SetTime(time.Now())

	if err := wapi.WritePoint(context.Background(), p); err != nil {
		panic(err)
	}
}

func InitClient(addr string, token string, org string, bucket string) api.WriteAPIBlocking {
	client := influxdb2.NewClient(
		addr,
		token,
	)

	return client.WriteAPIBlocking(org, bucket)
}
