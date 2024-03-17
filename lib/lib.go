package lib

import (
	"context"
	"strings"
	"text/template"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxClient struct {
	qapi        api.QueryAPI
	wapi        api.WriteAPIBlocking
	bucket      string
	measurement string
}

func (c *InfluxClient) Record(id string, point int, when time.Time) {
	p := influxdb2.NewPointWithMeasurement(c.measurement).
		AddTag("id", id).
		AddField("point", point).
		SetTime(when)

	if err := c.wapi.WritePoint(context.Background(), p); err != nil {
		panic(err)
	}
}

type RankRecord struct {
	Id    string
	Point int
}

func Rank(qapi api.QueryAPI, rng string) []RankRecord {
	var buf strings.Builder
	var tmplName, tmplPath, rQ string

	switch rng {
	case "weekly":
		tmplName = "rank.boundary.flux"
		tmplPath = "./lib/queries/rank.boundary.flux"
		rQ = "boundaries.week()"
	case "monthly":
		tmplName = "rank.boundary.flux"
		tmplPath = "./lib/queries/rank.boundary.flux"
		rQ = "boundaries.month()"
	case "all":
		tmplName = "rank.flux"
		tmplPath = "./lib/queries/rank.flux"
		rQ = "0"
	default:
		panic("unknown range")
	}

	if tmpl, err := template.New(tmplName).ParseFiles(
		tmplPath,
	); err != nil {
		panic(err)
	} else if err := tmpl.Execute(&buf, struct {
		Bucket      string
		Range       string
		Measurement string
		Limit       int
	}{
		"hello",
		rQ,
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
	var buf strings.Builder

	if tmpl, err := template.New("garden.flux").ParseFiles(
		"./lib/queries/garden.flux",
	); err != nil {
		panic(err)
	} else if err := tmpl.Execute(&buf, struct {
		Bucket      string
		Start       string
		Measurement string
		Id          string
		Window      string
	}{
		"hello",
		"-30d",
		"chat",
		"662201438621138954",
		"1d",
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

type Histogram interface {
	Process(data []int, height int) func(int) int
}
