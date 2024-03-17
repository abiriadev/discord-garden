package lib

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	influxquery "github.com/influxdata/influxdb-client-go/v2/api/query"
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

func (c *InfluxClient) queryInner(query string) ([][]*influxquery.FluxRecord, error) {
	tables := make([][]*influxquery.FluxRecord, 0)

	res, err := c.qapi.Query(
		context.Background(),
		query,
	)
	if err != nil {
		return tables, err
	}

	for res.Next() {
		if res.TableChanged() {
			tables = append(tables, make([]*influxquery.FluxRecord, 0))
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

	res, err := c.queryInner(query)
	if err != nil {
		panic(err)
	}
	if len(res) != 1 {
		panic("unexpected number of tables")
	}

	return lo.Map(res[0], func(r *influxquery.FluxRecord, _ int) RankRecord {
		id, ok := r.ValueByKey("id").(string)
		if !ok {
			id = "anon"
		}

		return RankRecord{
			Id:    id,
			Point: int(r.Value().(int64)),
		}
	})
}

func (c *InfluxClient) Garden() []int {
	query, err := useTemplate("./lib/queries/garden.flux", struct {
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
	})
	if err != nil {
		panic(err)
	}

	res, err := c.queryInner(query)
	if err != nil {
		panic(err)
	}
	if len(res) != 1 {
		panic("unexpected number of tables")
	}

	return lo.Map(res[0], func(r *influxquery.FluxRecord, _ int) int {
		return int(r.Value().(int64))
	})
}
