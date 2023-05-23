package nordpool

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
