# Discord Garden

## Run

```sh
$ DISCORD_TOKEN=<token> INFLUX_TOKEN=<influx token> INFLUX_PASSWORD=<password> docker compose up -d
```

## Stack

-   [https://github.com/bwmarrin/discordgo](discordgo)
-   [https://github.com/influxdata/influxdb](InfluxDB)

## Plan

-   [x] Show `weekly`, `monthly`, `total` ranking per server
-   [x] Show garden of each user
-   [ ] Support various garden histogram functions
    -   [x] Binary-mean
-   [x] Show basic bot and database status
-   [x] Handle error and show basic error information
-   [ ] Spam record deletion feature
-   [ ] Admin pandel
-   [ ] Support customizing timezones per user
-   [ ] Support multiple servers
-   [ ] Migrate to Rust

## License

[![GitHub](https://img.shields.io/github/license/abiriadev/pia?color=39d353&style=for-the-badge)](./LICENSE)
