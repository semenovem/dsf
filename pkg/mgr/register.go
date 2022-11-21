package mgr

func (a *mgr) RegisterStarted(fn func()) {
  if a.wait {
    a.log.Panic(ErrHandler)
  }
  a.fnsStarted = append(a.fnsStarted, fn)
}

func (a *mgr) RegisterFailed(fn func()) {
  if a.wait {
    a.log.Panic(ErrHandler)
  }
  a.fnsFailed = append(a.fnsFailed, fn)
}

func (a *mgr) RegisterBeforeStarted(fn func() error) {
  if a.wait {
    a.log.Panic(ErrHandler)
  }
  a.fnsAfterTasks = append(a.fnsAfterTasks, fn)
}

func (a *mgr) fireStarted() {
  for _, fn := range a.fnsStarted {
    go fn()
  }
  a.fnsStarted = nil
}

func (a *mgr) fireFailed() {
  for _, fn := range a.fnsFailed {
    go fn()
  }
  a.fnsFailed = nil
}
