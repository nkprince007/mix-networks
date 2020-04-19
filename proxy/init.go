package main

import (
	"time"

	"github.com/nkprince007/mix-networks/mixes"
)

func getThresholdMix() mixes.Mix {
	mix := &mixes.ThresholdMix{Size: 4}
	mix.Init()
	return mix
}

func getTimedMix() mixes.Mix {
	mixTimeBufferSize := 5000 * time.Millisecond
	mix := &mixes.TimedMix{TimeBufferMillis: mixTimeBufferSize}
	mix.Init()
	return mix
}

func getCottrellMix() mixes.Mix {
	mixTimeBufferSize := 5000 * time.Millisecond
	mix := &mixes.CottrellMix{
		TimeBufferMillis: mixTimeBufferSize,
		MinimumPoolSize:  3,
		Threshold:        5,
		Fraction:         float32(0.5),
	}
	mix.Init()
	return mix
}

func getRGBMix() mixes.Mix {
	mix := &mixes.RgbMix{
		PeriodMillis: 5000 * time.Millisecond,
	}
	mix.Init()
	return mix
}
