package transport

import (
  "net/http"
  "runtime"
  "github.com/rs/zerolog"
)

const (
    DISABLE_KEEP_ALIVES bool = false
    MAX_IDLE_CONNS int = 1
    IDLE_CONN_TIMEOUT int = 0
)

var (
    logger zerolog.Logger
    logLevel string
    tr *http.Transport
)


func init() {

	// #############  ENV vars ############# //
	logLevel = os.Getenv("LOG_LEVEL")

	// ############# Logging ############# //
	// Set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if logLevel == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	logger = zerolog.New(os.Stderr).With().Logger()
	logger.Debug().Msg(runtime.Caller(0), " initialized")

    // ############# Init ############# //
    tr = &http.Transport{
        DisableKeepAlives: DISABLE_KEEP_ALIVES,
        MaxIdleConns: MAX_IDLE_CONNS,
        IdleConnTimeout: IDLE_CONN_TIMEOUT
    }
}

func GetMainContainerMetrics(headerName string, headerValue string) (string, error){
    request := getNewRequest(headerName, headerValue)
    res, err := tr.RoundTrip(request)
    if err != nil {
      	logger.Error().Msg("Error in transport RoundTrip: " ,err)
    }
	logger.Debug().Msg("###########################")
    logger.Debug().Msg(time.Now())
    body, err := ioutil.ReadAll(res.Body)
    text := string(body)
    return text, err
}

func getNewRequest(headerName string, headerValue string) (*Request, error){
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set(headerName,headerValue)
    req.Close = false
    return req,_
}
