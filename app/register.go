package app

func (a *apptool) RegisterStarted(fn func()) {
	if a.wait {
		a.Log.Panic("Handler cannot be added after Wait")
	}
	a.fnsStarted = append(a.fnsStarted, fn)
}

func (a *apptool) RegisterFailed(fn func()) {
	if a.wait {
		a.Log.Panic("Handler cannot be added after Wait")
	}
	a.fnsFailed = append(a.fnsFailed, fn)
}

func (a *apptool) fireStarted() {
	for _, fn := range a.fnsStarted {
		go fn()
	}
	a.fnsStarted = nil
}

func (a *apptool) fireFailed() {
	for _, fn := range a.fnsFailed {
		go fn()
	}
	a.fnsFailed = nil
}
