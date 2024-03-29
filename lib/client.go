package lib

import (
	"context"
	"errors"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	influxquery "github.com/influxdata/influxdb-client-go/v2/api/query"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/samber/lo"
)

type InfluxClient struct {
	client      influxdb2.Client
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
		client:      client,
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

func (c *InfluxClient) Record(id string, point int, when time.Time) error {
	p := influxdb2.NewPointWithMeasurement(c.measurement).
		AddTag("id", id).
		AddField("point", point).
		SetTime(when)

	return c.wapi.WritePoint(context.Background(), p)
}

type RankRecord struct {
	Id    string
	Point int
}

func (c *InfluxClient) Rank(rng string) ([]RankRecord, error) {
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
		return nil, errors.New("unknown range")
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
		return nil, err
	}

	res, err := c.queryInner(query)
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, errors.New("unexpected number of tables")
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
	}), nil
}

func (c *InfluxClient) Garden(id string) ([]int, error) {
	query, err := useTemplate("./lib/queries/garden.flux", struct {
		Bucket      string
		Start       string
		Measurement string
		Id          string
		Window      string
	}{
		c.bucket,
		"-29d",
		c.measurement,
		id,
		"1d",
	})
	if err != nil {
		return nil, err
	}

	res, err := c.queryInner(query)
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, err
	}
	if len(res[0]) != 30 {
		return nil, errors.New("invalid result size")
	}

	return lo.Map(res[0], func(r *influxquery.FluxRecord, _ int) int {
		if v, ok := r.Value().(int64); !ok {
			return 0
		} else {
			return int(int(v))
		}
	}), nil
}

type InfluxStatus struct {
	Ready  domain.Ready
	Health domain.HealthCheck
}

func (c *InfluxClient) Status() (InfluxStatus, error) {
	ctx := context.Background()
	ready, err := c.client.Ready(ctx)
	if err != nil {
		return InfluxStatus{}, err
	}

	health, err := c.client.Health(ctx)
	if err != nil {
		return InfluxStatus{}, err
	}

	return InfluxStatus{
		Ready:  *ready,
		Health: *health,
	}, nil
}
