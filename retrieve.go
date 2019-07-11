package psdll

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gocloud.dev/blob"
)

// ReadFromURL returns dead-letter-log files from present URL.
func ReadFromURL(ctx context.Context, uri string) (map[string]DeadLetterLog, error) {
	bucket, files, err := retrieveFiles(ctx, uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrive %q", uri)
	}
	defer func() {
		if err := bucket.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	logs := make(map[string]DeadLetterLog)
	for _, file := range files {
		b, err := bucket.ReadAll(ctx, file)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read file %q", file)
		}
		var log DeadLetterLog
		if err := json.Unmarshal(b, &log); err != nil {
			return nil, errors.Wrapf(err, "%q is not a JSON file", file)
		}
		logs[file] = log
	}
	return logs, nil
}

func retrieveFiles(ctx context.Context, uri string) (*blob.Bucket, []string, error) {
	bucket, err := blob.OpenBucket(ctx, uri)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	key := strings.TrimLeft(u.Path, "/")
	if len(key) > 0 {
		if exists, err := bucket.Exists(ctx, key); err != nil {
			return nil, nil, errors.WithStack(err)
		} else if exists {
			return bucket, []string{key}, nil
		}
	}
	var files []string
	i := bucket.List(&blob.ListOptions{
		Prefix: key,
	})
	for {
		obj, err := i.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		if isMatchFileOrDir(key, obj.Key) {
			files = append(files, obj.Key)
		}
	}
	return bucket, files, nil
}

func isMatchFileOrDir(query, actual string) bool {
	a := strings.Split(actual, "/")
	for i, s := range strings.Split(query, "/") {
		if s != "" && s != a[i] {
			return false
		}
	}
	return true
}
