package EventCenter

import (
	"testing"
	"time"
)

type HelloCommand struct {
	Uid uint64
	Ctx string
}

func TestEventCenter(t *testing.T) {
	helloEc := EventCenter[HelloCommand]{}
	helloEc.On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hi\n", event.Uid, event.Ctx)
	}).On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hello\n", event.Uid, event.Ctx)
	})
	offReply := func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Fuck you\n", event.Uid, event.Ctx)
	}
	helloEc.On(offReply)
	helloEc.Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"})
	t.Logf("===========================================\n")
	helloEc.Off(offReply)
	helloEc.Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"})

	t.Logf("===========================================\n")
	offReply = func(event HelloCommand) {
		time.Sleep(time.Second)
		t.Logf("reply to %d, \"%s\", it's slow\n", event.Uid, event.Ctx)
	}
	helloEc.OnMonitor(offReply, time.Millisecond*50, func(event HelloCommand, elapse time.Duration) {
		t.Logf("    warning! cmd:%+v execution time:%d is too long", event, elapse.Milliseconds())
	})
	helloEc.OnMonitor(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", it's not slow\n", event.Uid, event.Ctx)
	}, time.Millisecond*50, func(event HelloCommand, elapse time.Duration) {
		t.Logf("    warning! cmd:%+v execution time:%d is too long", event, elapse)
	})
	helloEc.Fire(HelloCommand{Uid: 5555, Ctx: "hohoho"})

	t.Logf("===========================================\n")
	helloEc.OffMonitor(offReply)
	helloEc.Fire(HelloCommand{Uid: 5555, Ctx: "hohoho"})
	t.Logf("===========================================\n")
	helloEc.OffMonitor(offReply)
	helloEc.Fire(HelloCommand{Uid: 5555, Ctx: "hohoho"})

	t.Logf("===========================================\n")
	helloEc.Off(func(event HelloCommand) {})
	helloEc.OffMonitor(func(event HelloCommand) {})
	helloEc.Fire(HelloCommand{Uid: 8888, Ctx: "world"})
}

type WorldCommand struct {
	Uid uint64
	Ctx string
}

func TestEventCenterMgr(t *testing.T) {
	On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hi\n", event.Uid, event.Ctx)
	}).On(func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Hello\n", event.Uid, event.Ctx)
	})
	offReply := func(event HelloCommand) {
		t.Logf("reply to %d, \"%s\", Fuck you\n", event.Uid, event.Ctx)
	}
	On(offReply)
	Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"})
	t.Logf("===========================================\n")
	Off(offReply)
	Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"})

	t.Logf("===========================================\n")
	offReply2 := func(event WorldCommand) {
		time.Sleep(time.Second)
		t.Logf("reply to %d, \"%s\", it's slow\n", event.Uid, event.Ctx)
	}
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
	t.Logf("===========================================\n")
	OffMonitor(offReply2)
	Fire(WorldCommand{Uid: 5555, Ctx: "hohoho"})

	t.Logf("===========================================\n")
	Off(func(event WorldCommand) {})
	OffMonitor(func(event WorldCommand) {})
	Fire(WorldCommand{Uid: 8888, Ctx: "world"})

	Fire(HelloCommand{Uid: 12345, Ctx: "hahaha"})
	Fire(WorldCommand{Uid: 8888, Ctx: "world"})
}
