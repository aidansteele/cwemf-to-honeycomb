package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigStateFromEnv})
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}

	p := &publisher{
		s3man:   s3manager.NewDownloader(sess),
		dataset: os.Getenv("DATASET"),
		teamKey: os.Getenv("TEAM_KEY"),
	}

	lambda.Start(p.handle)
}

type publisher struct {
	s3man   *s3manager.Downloader
	dataset string
	teamKey string
}

func (p *publisher) handle(ctx context.Context, input *events.S3Event) error {
	j, _ := json.Marshal(input)
	fmt.Println(string(j))

	bucket := input.Records[0].S3.Bucket.Name
	key := input.Records[0].S3.Object.URLDecodedKey

	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := p.s3man.DownloadWithContext(ctx, buf, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		return errors.WithStack(err)
	}

	gzr, err := gzip.NewReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return errors.WithStack(err)
	}

	post := &bytes.Buffer{}
	gzw := gzip.NewWriter(post)
	gzw.Write([]byte("["))
	needsComma := false

	scan := bufio.NewScanner(gzr)
	for scan.Scan() {
		if needsComma {
			gzw.Write([]byte(","))
		}

		gzw.Write(scan.Bytes())
		needsComma = true
	}

	gzw.Write([]byte("]"))
	gzw.Close()

	u := fmt.Sprintf("https://api.honeycomb.io/1/batch/%s", url.PathEscape(p.dataset))
	req, err := http.NewRequestWithContext(ctx, "POST", u, post)
	req.Header.Set("X-Honeycomb-Team", p.teamKey)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{Transport: &logtransport{}}
	resp, err := c.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}

	if resp.StatusCode >= 300 {
		return errors.Errorf("bad status code %s", resp.Status)
	}

	return nil
}

type logtransport struct {

}

func (l *logtransport) RoundTrip(request *http.Request) (resp *http.Response, err error) {
	dump, _ := httputil.DumpRequestOut(request, true)
	fmt.Println(string(dump))

	resp, err = http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return
	}

	dump, _ = httputil.DumpResponse(resp, true)
	fmt.Println(string(dump))

	return
}
