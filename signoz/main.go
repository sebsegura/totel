package main

import (
	"context"
	"github.com/Bancar/goala/ulog"
	"totel/signoz/metrics"
)

func main() {
	ctx := context.Background()
	ulog.Info("starting meter provider")
	mp := metrics.InitMeter()
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			ulog.With(ulog.Str("error", err.Error())).Error("cannot shutdown meter provider")
			panic(err)
		}
	}()
	ulog.Info("started meter provider")

	if err := metrics.IncrementCounter(ctx, mp, "run.counter"); err != nil {
		ulog.With(ulog.Str("error", err.Error())).Error("cannot increment counter")
		panic(err)
	}
	ulog.Info("metric emitted")
}
