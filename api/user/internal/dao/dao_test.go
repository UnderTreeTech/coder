package dao

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/UnderTreeTech/waterdrop/pkg/log"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"
)

var d Dao
var ctx = context.Background()
var retry = 12

func TestMain(m *testing.M) {
	flag.Set("conf", "../../tests/application.toml")
	flag.Parse()

	conf.Init()
	defer log.New(nil).Sync()

	d = New()
	defer d.Close()

	// wait for docker services up to health
	// panic if don't come to health after 2min
	for i := 0; i < retry; i++ {
		if err := d.Ping(context.Background()); err == nil {
			break

		}

		if i > retry {
			panic("connect to database fail")
		}

		time.Sleep(time.Second * 10)
	}

	ret := m.Run()
	os.Exit(ret)
}
