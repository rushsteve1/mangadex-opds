package shared

import "net/url"

var GlobalOptions Options

type Options struct {
	Language  string
	Query     url.Values
	DataSaver bool
	MDUploads bool
	DevApi    bool
}

func ReadOptionsFromEnv() (o Options, err error) {
	return o, err
}

func TestOptions() Options {
	return Options{
		Language:  "en",
		DataSaver: true,
		DevApi:    true,
	}
}
