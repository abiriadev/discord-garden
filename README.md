<h1 align="center">🌱 Discord Garden</h1>
<p align="center">Show your discord activities as a GitHub-calendar-like garden</p>

## Run

```sh
$ DISCORD_TOKEN=<token> INFLUX_TOKEN=<influx token> INFLUX_PASSWORD=<password> docker compose up -d
```

## Stack

-   [discordgo](https://github.com/bwmarrin/discordgo)
-   [InfluxDB](https://github.com/influxdata/influxdb)

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
