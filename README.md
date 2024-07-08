# EventCenter

事件分发、响应工具。非线程安全。

## Feature & Example usage

监听指定事件
```golang
	type HelloCommand struct {
		Uid uint64
		Ctx string
	}

	On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hi\n", event.Uid, event.Ctx)
	})

	Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"}) //reply to 12345, "hahaha", Hi
```

链式写法
```golang
	On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hi\n", event.Uid, event.Ctx)
	}).On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hello\n", event.Uid, event.Ctx)
	})
```

```golang
	offReply := func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", going to off\n", event.Uid, event.Ctx)
	}
	On(offReply)
	Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"})
	Off(offReply) //取消对HelloCommand的监听
	Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"})
```

OnMonitor监控响应事件运行时常
```golang
	offReply2 := func(event WorldCommand) {
		time.Sleep(time.Second)
		t.Logf("reply to %d, \"%s\", it's slow\n", event.Uid, event.Ctx)
	}
	//指定响应函数如果运行超过50毫秒，则运行指定的回调函数
	OnMonitor(offReply2, time.Millisecond*50, func(event WorldCommand, elapse time.Duration) {
		t.Logf("    warning! cmd:%+v execution time:%d is too long", event, elapse.Milliseconds())
	})
	OnMonitor(func(event WorldCommand) {
		t.Logf("reply to %d, \"%s\", it's not slow\n", event.Uid, event.Ctx)
	}, time.Millisecond*50, func(event WorldCommand, elapse time.Duration) {
		t.Logf("    warning! cmd:%+v execution time:%d is too long", event, elapse)
	})
	Fire(WorldCommand{Uid: 5555, Ctx: "hohoho"})

	t.Logf("===========================================\n")
	OffMonitor(offReply2)
	Fire(WorldCommand{Uid: 5555, Ctx: "hohoho"})
```
