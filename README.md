# EventCenter

事件分发、响应工具。非线程安全。

## Feature & Example usage

监听指定事件
```golang
	import (
		"github.com/Lei2050/EventCenter"
	)

	type HelloCommand struct {
		Uid uint64
		Ctx string
	}
	
	EventCenter.On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hi\n", event.Uid, event.Ctx)
	})
	
	EventCenter.Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"}) //reply to 12345, "hahaha", Hi
```

链式写法
```golang
	EventCenter.On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hi\n", event.Uid, event.Ctx)
	}).On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hello\n", event.Uid, event.Ctx)
	})
```

```golang
	offReply := func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", going to off\n", event.Uid, event.Ctx)
	}
	EventCenter.On(offReply)
	EventCenter.Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"})
	EventCenter.Off(offReply) //取消对HelloCommand的监听
	EventCenter.Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"})
```

OnMonitor监控响应事件运行时常
```golang
	offReply2 := func(event WorldCommand) {
		time.Sleep(time.Second)
		t.Logf("reply to %d, \"%s\", it's slow\n", event.Uid, event.Ctx)
	}
	//指定响应函数如果运行超过50毫秒，则运行指定的回调函数
	EventCenter.OnMonitor(offReply2, time.Millisecond*50, func(event WorldCommand, elapse time.Duration) {
		t.Logf("    warning! cmd:%+v execution time:%d is too long", event, elapse.Milliseconds())
	})
	EventCenter.OnMonitor(func(event WorldCommand) {
		t.Logf("reply to %d, \"%s\", it's not slow\n", event.Uid, event.Ctx)
	}, time.Millisecond*50, func(event WorldCommand, elapse time.Duration) {
		t.Logf("    warning! cmd:%+v execution time:%d is too long", event, elapse)
	})
	EventCenter.Fire(WorldCommand{Uid: 5555, Ctx: "hohoho"})

	t.Logf("===========================================\n")
	EventCenter.OffMonitor(offReply2)
	EventCenter.Fire(WorldCommand{Uid: 5555, Ctx: "hohoho"})
```
