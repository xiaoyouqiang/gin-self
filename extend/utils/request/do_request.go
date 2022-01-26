package request

import (
	"context"
	"errors"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
)

func HttpGetByCtx(ctx context.Context,url string) (string,error)  {
	resp,err := ctxhttp.Get(ctx,http.DefaultClient,url)

	if resp != nil {
		defer resp.Body.Close()
	}

	if resp == nil || err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}
	if len(b) == 0 {
		return "",errors.New("empty data")
	}

	return string(b),nil
}
