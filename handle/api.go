package handle

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
)

// so that fields and methods are easily accessible from one type.
type Options struct {
	PrometheusAddress string

	SecureEnabled bool
	TlsCert       string
	TlsCertKey    string
	TlsCA         string
}

var (
	minTime          = time.Unix(math.MinInt64/1000+62135596801, 0).UTC()
	maxTime          = time.Unix(math.MaxInt64/1000-62135596801, 999999999).UTC()
	minTimeFormatted = minTime.Format(time.RFC3339Nano)
	maxTimeFormatted = maxTime.Format(time.RFC3339Nano)
)

func parseQuery(query string) map[string]string {
	query = strings.TrimRight(strings.TrimLeft(query, "{"), "}")
	querys := strings.Split(query, ",")
	labels := make(map[string]string)
	for i := range querys {
		labelStr := querys[i]
		l := strings.Split(labelStr, "=")
		labels[l[0]] = l[1]
	}
	return labels

}

func parseTime(s string) (time.Time, error) {
	if t, err := strconv.ParseFloat(s, 64); err == nil {
		s, ns := math.Modf(t)
		ns = math.Round(ns*1000) / 1000
		return time.Unix(int64(s), int64(ns*float64(time.Second))), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}

	switch s {
	case minTimeFormatted:
		return minTime, nil
	case maxTimeFormatted:
		return maxTime, nil
	}
	return time.Time{}, errors.Errorf("cannot parse %q to a valid timestamp", s)
}

func parseDuration(s string) (time.Duration, error) {
	if d, err := strconv.ParseFloat(s, 64); err == nil {
		ts := d * float64(time.Second)
		if ts > float64(math.MaxInt64) || ts < float64(math.MinInt64) {
			return 0, errors.Errorf("cannot parse %q to a valid duration. It overflows int64", s)
		}
		return time.Duration(ts), nil
	}
	if d, err := model.ParseDuration(s); err == nil {
		return time.Duration(d), nil
	}
	return 0, errors.Errorf("cannot parse %q to a valid duration", s)
}
