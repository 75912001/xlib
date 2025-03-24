package main

import (
	"context"
	"fmt"
	xconstants "github.com/75912001/xlib/constants"
	xutil "github.com/75912001/xlib/control"
	xtimer "github.com/75912001/xlib/timer"
	"time"
)

func cbSecond(arg ...interface{}) error {
	fmt.Printf("cbSecond:%v\n", arg...)
	return nil
}

func cbMillisecond(arg ...interface{}) error {
	fmt.Printf("cbMillisecond:%v\n", arg...)
	return nil
}

type addSecondSignal struct {
}
type addMillisecondSignal struct {
}

func exampleTimer() {
	if false {
		return
	}
	var timer xtimer.ITimer
	timer = xtimer.NewTimer()
	busChannel := make(chan interface{}, xconstants.BusChannelCapacityDefault)
	err := timer.Start(context.Background(),
		xtimer.NewOptions().
			WithOutgoingTimerOutChan(busChannel),
	)
	if err != nil {
		panic(err)
	}

	busChannel <- addSecondSignal{}
	busChannel <- addMillisecondSignal{}
	for {
		select {
		case v := <-busChannel:
			switch t := v.(type) {
			case addSecondSignal:
				for i := 0; i < 10; i++ {
					defaultCallBack := xutil.NewCallBack(cbSecond, uint64(i))
					second := timer.AddSecond(defaultCallBack, time.Now().Unix()+int64(i))
					switch i {
					case 3, 7, 9:
						timer.DelSecond(second)
					default:
					}
				}
			case addMillisecondSignal:
				for i := 0; i < 10000; i += 1000 {
					defaultCallBack := xutil.NewCallBack(cbMillisecond, uint64(i))
					millisecond := timer.AddMillisecond(defaultCallBack, time.Now().UnixMilli()+int64(i))
					switch i {
					case 3000, 7000, 9000:
						timer.DelMillisecond(millisecond)
					default:
					}
				}
			case *xtimer.EventTimerSecond:
				_ = t.ICallBack.Execute()
			case *xtimer.EventTimerMillisecond:
				_ = t.ICallBack.Execute()
			}
		}
	}
	return
}
