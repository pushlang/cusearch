package google

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pushlang/cusearch/v2/contxt"
)

type Results []Result

type Result struct {
	Title string 
	URL string
}

func Search(ctx context.Context, query string) (Results, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/customsearch/v1?key=gIzaSyB7JRboB996gA-goCGlDB8DfTIaIVD1ZRo", nil)
	
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", query)

	if cx, ok := contxt.FromContext(ctx); ok {
		q.Set("cx", cx)
	}
	
	req.URL.RawQuery = q.Encode()
	
	var results Results
	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data struct {
			Results []struct {
				TitleNoFormatting string `json:"title"`
				URL               string `json:"link"`
			} `json:"items"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		for _, res := range data.Results {
			results = append(results, Result{Title: res.TitleNoFormatting, URL: res.URL})
		}
		return nil
	})

	return results, err
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	c := make(chan error, 1)
	req = req.WithContext(ctx)
	go func() { c <- f(http.DefaultClient.Do(req)) }()
	select {
	case <-ctx.Done():
		<-c
		return ctx.Err()
	case err := <-c:
		return err
	}
}
