# Nordpool Elspot API Go client

This is an unofficial Nordpool Elsport API client for Go
It is heavily inspired by https://github.com/samuelmr/nordpool-node

## Installation

`go get github.com/Frelsaren/nordpool`

## Usage

To use the nordpool module in your Go code, you can import it using the appropriate import path:

`import "github.com/yourusername/nordpool"`

The nordpool module provides a GetHourlyPrices function that allows you to retrieve hourly electricity prices for a specific region and date:

```
func main() {
	prices := nordpool.Prices{}
	options := nordpool.PriceOptions{
		Area:     "Kr.sand",
		Currency: "NOK",
	}
	data, err := prices.Hourly(options)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
```

This code retrieves the hourly electricity prices for the Kr.sand region in Norway on January 1, 2022, and prints them to the console.

Contributing
Contributions to the nordpool module are welcome! If you find a bug or have a feature request, please open an issue on the GitHub repository. If you would like to contribute code, please fork the repository and submit a pull request.

License
The nordpool module is licensed under the MIT License. See the LICENSE file for more information.
