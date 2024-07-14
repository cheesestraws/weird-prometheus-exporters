package wiltshirebins

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const BinURL = "https://ilforms.wiltshire.gov.uk//wastecollectiondays/collectionlist"

type Client struct {
	httpClient *http.Client
}

var DefaultClient *Client = &Client{
	httpClient: http.DefaultClient,
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) Get(ctx context.Context, month int, year int, postcode string, uprn string) (Calendar, error) {
	var empty Calendar

	// Make the body
	vs := url.Values{}

	vs.Add("Month", strconv.Itoa(month))
	vs.Add("Year", strconv.Itoa(year))
	vs.Add("Postcode", postcode)
	vs.Add("Uprn", uprn)

	bodyString := vs.Encode()
	body := strings.NewReader(bodyString)

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", BinURL, body)
	if err != nil {
		empty.errorFill("request_creation_error")
		return empty, err
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	// Send it out into the cold world
	resp, err := c.httpClient.Do(req)
	if err != nil {
		empty.errorFill("error_response")
		return empty, err
	}
	defer resp.Body.Close()

	// Parse it
	days, err := parse(resp.Body)
	if err != nil {
		days.errorFill("parse_error")
		return days, err
	}

	days.sanityCheck()

	return days, nil
}

func (c *Client) GetForDate(ctx context.Context, date time.Time, postcode string, uprn string) (Collections, error) {
	month := int(date.Month())
	year := int(date.Year())
	day := date.Day()

	calendar, err := c.Get(ctx, month, year, postcode, uprn)

	return calendar[day-1], err
}
