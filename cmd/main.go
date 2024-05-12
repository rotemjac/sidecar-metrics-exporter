package main

import (
  "os"
  "time"
  "strconv"
  "net/http"
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
    HEADER_NAME string
    HEADER_VALUE string
    NEW_METRICS_PORT string
)


func main(){

	// #############  ENV vars ############# //
	logLevel = os.Getenv("LOG_LEVEL")
	SCRAPING_ENDPOINT = os.Getenv("SCRAPING_ENDPOINT")
    SCRAPE_INTERVAL = os.Getenv("SCRAPE_INTERVAL")
    HEADER_NAME = os.Getenv("HEADER_NAME")
    HEADER_VALUE = os.Getenv("HEADER_VALUE")
    NEW_METRICS_PORT = os.Getenv("NEW_METRICS_PORT")

	// ############# Logging ############# //
	// Set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if logLevel == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	logger = zerolog.New(os.Stderr).With().Logger()
	logger.Debug().Msg("Main module initialized..")


    // ############# Main ############# //

    // Wait 30 seconds for main container and then start polling
    time.Sleep(30 * time.Second)
	logger.Debug().Msg("Waiting for main container to load..")
	for !mainContainerReady() {
        time.Sleep(5 * time.Second)
    }

    // Expose metrics endpoint
    go func(){
        http.Handle("/metrics", metrics.GetPromHttp())
        http.Handle("/readiness", probeHandler())
        port := ":"+NEW_METRICS_PORT
        http.ListenAndServe(port, nil)
    }()

    // Fetch metric list
    // Notice: Metrics are in 5 groups so the ENV var will not be a long string and because we don't want to add a configMap to the original Helm chart
    metricsList := mergeToOneList(
        os.Getenv("METRICS_GROUP1"),
        os.Getenv("METRICS_GROUP2"),
        os.Getenv("METRICS_GROUP3"),
        os.Getenv("METRICS_GROUP4"),
        os.Getenv("METRICS_GROUP5"),
    )
   	logger.Debug().Msgf("MetricsList: %s ", metricsList)
    metrics.RegisterAllMetrics(metricsList)

    // Start a cron job to poll main container every SCRAPE_INTERVAL seconds
    s := gocron.NewScheduler(time.UTC)
    interval,_ := strconv.Atoi(SCRAPE_INTERVAL)
    s.Every(interval).Seconds().Do(func() {
        fetchMetrics(SCRAPING_ENDPOINT, metricsList, HEADER_NAME, HEADER_VALUE)
    })
    s.StartBlocking()
}

func mainContainerReady() bool{
    res,err := transport.CallMainContainerEndpoint(SCRAPING_ENDPOINT, HEADER_NAME, HEADER_VALUE)
    if (res == "" || err != nil) {
      return false
    } else {
      return true
    }
}

func probeHandler() http.HandlerFunc{
    return func(w http.ResponseWriter, req *http.Request) {
        res,_ := transport.CallMainContainerEndpoint(SCRAPING_ENDPOINT, HEADER_NAME, HEADER_VALUE)
        if res == "" {
            w.WriteHeader(http.StatusOK)
        } else {
            w.WriteHeader(http.StatusBadRequest)
        }
    }
}


 func fetchMetrics(url string, metricsList []string, headerName string, headerValue string){
    text,_ := transport.CallMainContainerEndpoint(url, headerName, headerValue)

    // Split the text into lines
    lines := strings.Split(text, "\n")
    logger.Debug().Msgf("len(lines): %d", len(lines))

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
                logger.Debug().Msgf("%s %f", gaugeName, gaugeValue)
                metrics.SetGauge(gaugeName, gaugeValue)
             }
         }
    }
}




