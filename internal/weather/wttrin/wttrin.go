package wttrin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/noilpa/ctxlog"
	"github.com/sirupsen/logrus"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	cli    httpClient
	host   *url.URL
	format string
}

func New(c Config) *Client {
	h, err := url.Parse(c.Host)
	if err != nil {
		panic(err)
	}
	return &Client{
		host:   h,
		format: c.Format,
		cli: &http.Client{
			Transport: &http.Transport{
				IdleConnTimeout: 5 * time.Second,
			},
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) IsRainy(ctx context.Context, location string, threshold int) (bool, error) {
	c.host.Path = location
	v := url.Values{
		"format": {c.format},
	}
	c.host.RawQuery = v.Encode()
	req, err := http.NewRequest(http.MethodGet, c.host.String(), http.NoBody)
	if err != nil {
		return false, err
	}

	res, err := c.cli.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	var r WttrinResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return false, err
	}

	if len(r.Weather) == 0 {
		return false, errors.New("wttr.in unexpected weather length in response")
	}

	ctxlog.From(ctx).WithFields(logrus.Fields{
		"area_name": r.NearestArea[0].Areaname[0].Value,
		"country":   r.NearestArea[0].Country[0].Value,
		"region":    r.NearestArea[0].Region[0].Value,
	}).Info("nearest area")

	errCnt := 0
	for _, predict := range r.Weather[0].Hourly {
		v, err := strconv.Atoi(predict.Chanceofrain)
		if err != nil {
			errCnt++
			ctxlog.From(ctx).WithError(err).Error("wttr.in chance of rain has invalid format")
			continue
		}
		if v >= threshold {
			return true, nil
		}
	}

	if errCnt == len(r.Weather[0].Hourly) {
		return false, errors.New("wttr.in could not get chance of rain")
	}

	return false, nil
}
