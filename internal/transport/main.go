package transport

import (
  "os"
  "net/http"
  "time"
  "io/ioutil"
  //"runtime"
  "github.com/rs/zerolog"
)

const (
    DISABLE_KEEP_ALIVES bool = false
    MAX_IDLE_CONNS int = 1
    //IDLE_CONN_TIMEOUT int = 0
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
	//logger.Debug().Msg(runtime.Caller(0), " initialized")
	logger.Debug().Msg("Metrics module initialized..")


    // ############# Init ############# //
    tr = &http.Transport{
        DisableKeepAlives: DISABLE_KEEP_ALIVES,
        MaxIdleConns: MAX_IDLE_CONNS,
        IdleConnTimeout: 0,
    }
}

func CallMainContainerEndpoint(url string, headerName string, headerValue string) (string, error){
    request, err := getNewRequest(url, headerName, headerValue)
    if err != nil {
      	logger.Error().Msgf("Error in getNewRequest: %s" ,err)
    }
    res, err := tr.RoundTrip(request)
    if err != nil {
      	logger.Error().Msgf("Error in transport RoundTrip: %s" ,err)
    }
	logger.Debug().Msg("###########################")
    logger.Debug().Msg(time.Now().String())
    body, err := ioutil.ReadAll(res.Body)
    text := string(body)
    return text, err
}

func getNewRequest(url string, headerName string, headerValue string) (*http.Request, error){
    req, _ := http.NewRequest("GET", url, nil)
    if (headerName != "") {
        req.Header.Set(headerName,headerValue)
    }
    req.Close = false
    return req, nil
}
