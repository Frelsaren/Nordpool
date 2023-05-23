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

type PriceOptions struct {
	Area     string
	Currency string
	Date     string
	From     time.Time
	To       time.Time
}

type NordpoolRange struct {
	Hour  int
	Day   int
	Week  int
	Month int
}

type NordpoolOptions struct {
	URL           string
	MaxRange      NordpoolRange
	MaxRangeValue int
	Date          time.Time
	PriceOptions
}

type DateType struct{}

type Result struct {
	Date  time.Time `json:"date"`
	Value int       `json:"value"`
	Area  string    `json:"area"`
}

type Prices struct{}

func (p Prices) At(opts PriceOptions) (Result, error) {
	var date time.Time
	if opts.Date != "" {
		date, _ = time.Parse(time.RFC3339, opts.Date)
	} else {
		date = time.Now()
	}
	location, _ := time.LoadLocation("Europe/Oslo")
	date = date.In(location)

	results, err := p.getValues(NordpoolOptions{
		URL:          config.priceUrlHourly,
		MaxRange:     NordpoolRange{Hour: 1},
		Date:         date,
		PriceOptions: opts,
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

func (p Prices) Hourly(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		URL:      config.priceUrlHourly,
		MaxRange: NordpoolRange{Day: 1},
		Date:     date,
	}
	if opts.Date != "" {
		if t, err := time.Parse(time.RFC3339, opts.Date); err == nil {
			nordpoolOpts.Date = t
		} else {
			return nil, err
		}
	}
	if opts.Currency != "" {
		nordpoolOpts.Currency = opts.Currency
	}
	if opts.Area != "" {
		nordpoolOpts.Area = opts.Area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) Daily(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		URL:      config.priceUrlDaily,
		MaxRange: NordpoolRange{Day: 31},
		Date:     date,
	}
	if opts.Date != "" {
		t, err := time.Parse(time.RFC3339, opts.Date)
		if err != nil {
			return nil, err
		}
		nordpoolOpts.Date = t
	}
	if opts.Currency != "" {
		nordpoolOpts.Currency = opts.Currency
	}
	if opts.Area != "" {
		nordpoolOpts.Area = opts.Area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) Weekly(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		URL:      config.priceUrlWeekly,
		MaxRange: NordpoolRange{Week: 24},
		Date:     date,
	}
	if opts.Date != "" {
		t, err := time.Parse(time.RFC3339, opts.Date)
		if err != nil {
			return nil, err
		}
		nordpoolOpts.Date = t
	}
	if opts.Currency != "" {
		nordpoolOpts.Currency = opts.Currency
	}
	if opts.Area != "" {
		nordpoolOpts.Area = opts.Area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) Monthly(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		URL:      config.priceUrlMonthly,
		MaxRange: NordpoolRange{Month: 53},
		Date:     date,
	}
	if opts.Date != "" {
		t, err := time.Parse(time.RFC3339, opts.Date)
		if err != nil {
			return nil, err
		}
		nordpoolOpts.Date = t
	}
	if opts.Currency != "" {
		nordpoolOpts.Currency = opts.Currency
	}
	if opts.Area != "" {
		nordpoolOpts.Area = opts.Area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) Yearly(opts PriceOptions) ([]Result, error) {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return nil, err
	}
	date := time.Now().In(location)
	nordpoolOpts := NordpoolOptions{
		URL:  config.priceUrlYearly,
		Date: date,
	}
	if opts.Currency != "" {
		nordpoolOpts.Currency = opts.Currency
	}
	if opts.Area != "" {
		nordpoolOpts.Area = opts.Area
	}
	return p.getValues(nordpoolOpts)
}

func (p Prices) getValues(opts NordpoolOptions) ([]Result, error) {
	var fromTime time.Time
	if !opts.From.IsZero() {
		var err error
		fromTime, err = time.Parse(time.RFC3339, opts.From.Format(time.RFC3339))
		if err != nil {

			return nil, err
		}
		if err != nil {
			return nil, err
		}

	}
	var toTime time.Time
	if opts.To.IsZero() {
		toTime, _ = time.Parse(time.RFC3339, opts.To.Format(time.RFC3339))
	}
	var MaxRangeKey string
	var MaxRangeValue int
	if opts.MaxRange != (NordpoolRange{}) {
		if opts.MaxRange.Day != 0 {
			MaxRangeKey = "day"
			MaxRangeValue = opts.MaxRange.Day
		} else if opts.MaxRange.Hour != 0 {
			MaxRangeKey = "hour"
			MaxRangeValue = opts.MaxRange.Hour
		} else if opts.MaxRange.Month != 0 {
			MaxRangeKey = "month"
			MaxRangeValue = opts.MaxRange.Month
		} else if opts.MaxRange.Week != 0 {
			MaxRangeKey = "week"
			MaxRangeValue = opts.MaxRange.Week
		}
	}
	if !fromTime.IsZero() && !toTime.IsZero() && MaxRangeKey != "" && MaxRangeValue != 0 {
		minFromTime := toTime.Add(-time.Duration(MaxRangeValue) * time.Hour)
		if fromTime.Before(minFromTime) {
			fmt.Println("Time span too long. Setting start time to", minFromTime.Format(time.RFC3339))
			fromTime = minFromTime
		}
	}

	currency := opts.Currency
	if currency == "" {
		currency = "EUR"
	}
	location, err := time.LoadLocation(config.timezone)
	if err != nil {
		return nil, err
	}
	date, err := time.ParseInLocation(time.RFC3339, opts.Date.Format(time.RFC3339), location)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?currency=,%s,%s,%s&endDate=%s", opts.URL, currency, currency, currency, date.Format("02-01-2006"))

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
				if opts.Area == "" || opts.Area == area {
					values = append(values, Result{Area: area, Date: date, Value: int(value)})
				}
			}
		}
		return values, nil
	}

	return []Result{}, nil
}
