package options

import "github.com/BurntSushi/toml"

type (
	HandleGet func(key, defVal string) (value string)
)

var (
	DefaultVal = "NULL"
)

type Options struct {
	fpath string
	hget  HandleGet
}

func NewOptions(fpath string, hget HandleGet) *Options {
	return &Options{
		fpath: fpath,
		hget:  hget,
	}
}

func (o *Options) Parsing(obj interface{}) error {
	if _, err := toml.DecodeFile(o.fpath, obj); err != nil {
		return err
	}

	if o.hget != nil {
		return resolve(obj, o.hget)
	}
	return nil
}

func (o *Options) Sync(obj interface{}) (sync chan<- struct{}) {
	ch := make(chan struct{})
	go func() {
		for range ch {
			resolve(obj, o.hget)
		}
	}()

	return ch
}
