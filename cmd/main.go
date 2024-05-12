package main

import (
  "os"
  "time"
  "strconv"
  "io/ioutil"
  "strings"
  "github.com/rs/zerolog"
  "github.com/go-co-op/gocron"
  "github.com/rotemjac/sidecar-metrics-exporter/internal/transport"
  "github.com/rotemjac/sidecar-metrics-exporter/internal/metrics"
)

var (
    logger zerolog.Logger
    logLevel string
    SCRAPING_ENDPOINT string
    SCRAPE_INTERVAL string
    TRINO_HEADER_NAME string
    TRINO_HEADER_VALUE string
    METRICS string
)


func main(){

	// #############  ENV vars ############# //
	logLevel = os.Getenv("LOG_LEVEL")
	SCRAPING_ENDPOINT = os.Getenv("SCRAPING_ENDPOINT")
    SCRAPE_INTERVAL = os.Getenv("SCRAPE_INTERVAL")
    TRINO_HEADER_NAME = os.Getenv("TRINO_HEADER_NAME")
    TRINO_HEADER_VALUE = os.Getenv("TRINO_HEADER_VALUE")
    NEW_METRICS_PORT = os.Getenv("NEW_METRICS_PORT")

	// ############# Logging ############# //
	// Set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if logLevel == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	logger = zerolog.New(os.Stderr).With().Logger()
	logger.Debug().Msg(runtime.Caller(0), " initialized")


    // ############# Main ############# //
    // Expose metrics

    go func(){
        http.Handle("/metrics", promhttp.Handler())
        port := ":"+NEW_METRICS_PORT
        http.ListenAndServe(port, nil)
    }()

    // Fetch metric list
    // Notice: Metrics are in 5 groups so the ENV var will not be a long string and because we don't want to add a configMap to the original Helm chart
    metricsList := mergeLists(
        os.Getenv("METRICS_GROUP1"),
        os.Getenv("METRICS_GROUP2"),
        os.Getenv("METRICS_GROUP3"),
        os.Getenv("METRICS_GROUP4"),
        os.Getenv("METRICS_GROUP5")
    )
   	logger.Debug().Msg("MetricsList: ", metricsList)
    metrics.RegisterAllMetrics(metricsList)

    // Start a cron job to poll main container every SCRAPE_INTERVAL seconds
    s := gocron.NewScheduler(time.UTC)
    interval,_ := strconv.Atoi(SCRAPE_INTERVAL)
    s.Every(interval).Seconds().Do(func() {
        fetchMetrics(SCRAPING_ENDPOINT, metricsList, TRINO_HEADER_NAME, TRINO_HEADER_VALUE)
    })
    s.StartBlocking()
}


 func fetchMetrics(url string, metricsList []string, headerName string, headerValue string){
    text := communication.GetMainContainerMetrics(headerName, headerValue)


    // Split the text into lines
    lines := strings.Split(text, "\n")
    logger.Debug().Msg("len(lines):", len(lines))

    // Iterate over lines and extract specific strings
    for _, line := range lines {
         if strings.Contains(line, "# TYPE") {
            continue;
         }
         for _,metric := range metricsList {
             if strings.Contains(line, metric) {
                arr := strings.Split(line, " ")
                valFloat, _ := strconv.ParseFloat(arr[1], 64)
                gaugeName := arr[0]
                gaugeValue := valFloat
                logger.Debug().Msg(gaugeName, gaugeValue)
                metrics.SetGauge(gaugeName, gaugeValue)
             }
         }
    }
}




