package lib

import (
	"context"
	"strings"
	"text/template"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/query"
	"github.com/samber/lo"
)

type InfluxClient struct {
	qapi        api.QueryAPI
	wapi        api.WriteAPIBlocking
	bucket      string
	measurement string
}

type InfluxClientConfig struct {
	Host        string
	Token       string
	Org         string
	Bucket      string
	Measurement string
}

func NewInfluxClient(
	config InfluxClientConfig,
) InfluxClient {
	client := influxdb2.NewClient(
		config.Host,
		config.Token,
	)

	return InfluxClient{
		qapi:        client.QueryAPI(config.Org),
		wapi:        client.WriteAPIBlocking(config.Org, config.Bucket),
		bucket:      config.Bucket,
		measurement: config.Measurement,
	}
}

type Fr = *query.FluxRecord

func (c *InfluxClient) queryInner(query string) ([][]*query.FluxRecord, error) {
	tables := make([][]Fr, 0)

	res, err := c.qapi.Query(
		context.Background(),
		query,
	)
	if err != nil {
		return tables, err
	}

	for res.Next() {
		if res.TableChanged() {
			tables = append(tables, make([]Fr, 0))
		}
		tables[len(tables)-1] = append(tables[len(tables)-1], res.Record())
	}
	if res.Err() != nil {
		return tables, res.Err()
	}

	return tables, nil
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

func (c *InfluxClient) Rank(rng string) []RankRecord {
	var tmplPath, rQ string

	switch rng {
	case "weekly":
		tmplPath = "./lib/queries/rank.boundary.flux"
		rQ = "boundaries.week()"
	case "monthly":
		tmplPath = "./lib/queries/rank.boundary.flux"
		rQ = "boundaries.month()"
	case "all":
		tmplPath = "./lib/queries/rank.flux"
		rQ = "0"
	default:
		panic("unknown range")
	}

	query, err := useTemplate(tmplPath, struct {
		Bucket      string
		Range       string
		Measurement string
		Limit       int
	}{
		c.bucket,
		rQ,
		c.measurement,
		10,
	})
	if err != nil {
		panic(err)
	}

	res, err := c.qapi.Query(
		context.Background(),
		query,
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

func (c *InfluxClient) Garden() []int {
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
		c.bucket,
		"-30d",
		c.measurement,
		"662201438621138954",
		"1d",
	}); err != nil {
		panic(err)
	}

	res, err := c.queryInner(buf.String())
	if err != nil {
		panic(err)
	}
	if len(res) != 1 {
		panic("unexpected number of tables")
	}

	return lo.Map(res[0], func(r *query.FluxRecord, _ int) int {
		return int(r.Value().(int64))
	})
}
