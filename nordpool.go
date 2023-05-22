package nordpool

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const baseUrl = "https://www.nordpoolgroup.com/api"

var config = struct {
	baseUrl         string
	priceUrlHourly  string
	priceUrlDaily   string
	priceUrlWeekly  string
	priceUrlMonthly string
	priceUrlYearly  string
	timezone        string
}{
	baseUrl:         baseUrl,
	priceUrlHourly:  baseUrl + "/marketdata/page/10",
	priceUrlDaily:   baseUrl + "/marketdata/page/11",
	priceUrlWeekly:  baseUrl + "/marketdata/page/12",
	priceUrlMonthly: baseUrl + "/marketdata/page/13",
	priceUrlYearly:  baseUrl + "/marketdata/page/14",
	timezone:        "Europe/Oslo",
}

type PriceOptions struct {
	area     string
	currency string
	date     string
	from     time.Time
	to       time.Time
}
type NordpoolRange struct {
	hour  int
	day   int
	week  int
	month int
}
type NordpoolOptions struct {
	url           string
	maxRange      NordpoolRange
	maxRangeValue int
	date          time.Time
	PriceOptions  // Omit<PriceOptions, "date">
}

type DateType struct{}

type Result struct {
	Date  time.Time `json:"date"`
	Value int       `json:"value"`
	Area  string    `json:"area"`
}

type Prices struct{}

func (p Prices) at(opts PriceOptions) (Result, error) {
	var date time.Time
	if opts.date != "" {
		date, _ = time.Parse(time.RFC3339, opts.date)
	} else {
		date = time.Now()
	}
	location, _ := time.LoadLocation("Europe/Oslo")
	date = date.In(location)

	results, err := p.getValues(NordpoolOptions{
		url:      config.priceUrlHourly,
		maxRange: NordpoolRange{hour: 1},
		date:     date,
	})
	if err != nil {
		return Result{}, err
	}
	if len(results) > 0 {
		for _, result := range results {
			if result.Date.Day() == date.Day() &&
				result.Date.Hour() == date.Hour() &&
				result.Date.Year() == date.Year() {
				return result, nil
			}
		}
	}

	return Result{}, fmt.Errorf("no results found for %s", date.Format(time.RFC3339))
}

func (p Prices) hourly(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		// handle error
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		url:      config.priceUrlHourly,
		maxRange: NordpoolRange{day: 1},
		date:     date,
	}
	if opts.date != "" {
		t, err := time.Parse(time.RFC3339, opts.date)
		if err != nil {
			return nil, err
		}
		nordpoolOpts.date = t
	}
	if opts.currency != "" {
		nordpoolOpts.currency = opts.currency
	}
	if opts.area != "" {
		nordpoolOpts.area = opts.area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) daily(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		url:      config.priceUrlDaily,
		maxRange: NordpoolRange{day: 31},
		date:     date,
	}
	if opts.date != "" {
		t, err := time.Parse(time.RFC3339, opts.date)
		if err != nil {
			return nil, err
		}
		nordpoolOpts.date = t
	}
	if opts.currency != "" {
		nordpoolOpts.currency = opts.currency
	}
	if opts.area != "" {
		nordpoolOpts.area = opts.area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) weekly(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		url:      config.priceUrlWeekly,
		maxRange: NordpoolRange{week: 24},
		date:     date,
	}
	if opts.date != "" {
		t, err := time.Parse(time.RFC3339, opts.date)
		if err != nil {
			return nil, err
		}
		nordpoolOpts.date = t
	}
	if opts.currency != "" {
		nordpoolOpts.currency = opts.currency
	}
	if opts.area != "" {
		nordpoolOpts.area = opts.area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) monthly(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		url:      config.priceUrlMonthly,
		maxRange: NordpoolRange{month: 53},
		date:     date,
	}
	if opts.date != "" {
		t, err := time.Parse(time.RFC3339, opts.date)
		if err != nil {
			return nil, err
		}
		nordpoolOpts.date = t
	}
	if opts.currency != "" {
		nordpoolOpts.currency = opts.currency
	}
	if opts.area != "" {
		nordpoolOpts.area = opts.area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) yearly(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		url:  config.priceUrlYearly,
		date: date,
	}
	if opts.currency != "" {
		nordpoolOpts.currency = opts.currency
	}
	if opts.area != "" {
		nordpoolOpts.area = opts.area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) getValues(opts NordpoolOptions) ([]Result, error) {
	var fromTime time.Time
	if !opts.from.IsZero() {
		var err error
		fromTime, err = time.Parse(time.RFC3339, opts.from.Format(time.RFC3339))
		if err != nil {

			return nil, err
		}
		if err != nil {
			return nil, err
		}

	}
	var toTime time.Time
	if opts.to.IsZero() {
		toTime, _ = time.Parse(time.RFC3339, opts.to.Format(time.RFC3339))
	}
	var maxRangeKey string
	var maxRangeValue int
	if opts.maxRange != (NordpoolRange{}) {
		if opts.maxRange.day != 0 {
			maxRangeKey = "day"
			maxRangeValue = opts.maxRange.day
		} else if opts.maxRange.hour != 0 {
			maxRangeKey = "hour"
			maxRangeValue = opts.maxRange.hour
		} else if opts.maxRange.month != 0 {
			maxRangeKey = "month"
			maxRangeValue = opts.maxRange.month
		} else if opts.maxRange.week != 0 {
			maxRangeKey = "week"
			maxRangeValue = opts.maxRange.week
		}
	}
	if !fromTime.IsZero() && !toTime.IsZero() && maxRangeKey != "" && maxRangeValue != 0 {
		minFromTime := toTime.Add(-time.Duration(maxRangeValue) * time.Hour)
		if fromTime.Before(minFromTime) {
			fmt.Println("Time span too long. Setting start time to", minFromTime.Format(time.RFC3339))
			fromTime = minFromTime
		}
	}

	currency := opts.currency
	if currency == "" {
		currency = "EUR"
	}
	location, err := time.LoadLocation(config.timezone)
	if err != nil {
		return nil, err
	}
	date, err := time.ParseInLocation(time.RFC3339, opts.date.Format(time.RFC3339), location)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?currency=,%s,%s,%s&endDate=%s", opts.url, currency, currency, currency, date.Format("02-01-2006"))

	resp, err := http.Get(url)

	if err != nil {
		return []Result{}, err
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return nil, err
	}
	data = data["data"].(map[string]interface{})
	// if data is not null and data.Rows is not null and data.Rows is longer than 0
	if data != nil && data["Rows"] != nil && len(data["Rows"].([]interface{})) > 0 {
		values := []Result{}
		for _, row := range data["Rows"].([]interface{}) {
			rowMap := row.(map[string]interface{})
			if rowMap["IsExtraRow"].(bool) {
				continue
			}
			date, err := time.ParseInLocation("2006-01-02T15:04:05", rowMap["StartTime"].(string), location)
			if err != nil || date.IsZero() {
				continue
			} else if (date.Unix() < fromTime.Unix()) || (toTime != time.Time{} && date.Unix() >= toTime.Unix()) {
				continue
			}
			for _, column := range rowMap["Columns"].([]interface{}) {
				columnMap := column.(map[string]interface{})
				valueStr := columnMap["Value"].(string)
				valueStr = regexp.MustCompile(`[^\d.]`).ReplaceAllString(valueStr, "")
				value, err := strconv.ParseFloat(valueStr, 64)
				if err != nil || math.IsNaN(value) {
					continue
				}
				area := columnMap["Name"].(string)
				if opts.area == "" || opts.area == area {
					values = append(values, Result{Area: area, Date: date, Value: int(value)})
				}
			}
		}
		return values, nil
	}

	return []Result{}, nil
}
