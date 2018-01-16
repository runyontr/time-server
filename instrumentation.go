package main

import (
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	t "time"
	"github.com/prometheus/client_golang/prometheus"
	"fmt"
)


var requestCount metrics.Counter
var requestLatency metrics.Histogram

func init(){
	//make the counters and metrics
	fieldKeys := []string{"method", "error"}

	requestCount = kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "runyontr",
		Subsystem: "time_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency = kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "runyontr",
		Subsystem: "time_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
}


//instrumentationTimeService implements the AppInfoService interface.  It provides metrics on calls to the Next service
type instrumentationTimeService struct{

	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Next TimeService
}



func NewInstrumentationTimeService(svc TimeService) TimeService{

	return &instrumentationTimeService{
		Next: svc,
		requestCount: requestCount,
		requestLatency: requestLatency,
	}
}


//GetAppInfo returns the app info of the running application
func (s *instrumentationTimeService) GetTime() (info time, err error) {
	defer func(startTime t.Time){
		requestCount.With( "method","GetTime","error",fmt.Sprintf("%v",err)).Add(1)
		requestLatency.With( "method","GetTime","error",fmt.Sprintf("%v",err)).Observe(float64(t.Since(startTime)))
	}(t.Now())
	info, err = s.Next.GetTime()
	return
}

