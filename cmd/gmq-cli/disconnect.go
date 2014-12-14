package main

// disconnectWithLock locks the Client and disconnects the Network Connection.
func disconnectWithLock(ctx *context) error {
	ctx.climu.Lock()
	defer ctx.climu.Unlock()

	return disconnect(ctx)
}

// disconnect disconnects the Network Connection.
func disconnect(ctx *context) error {
	// Disconnect the Network Connection.
	if err := ctx.cli.Disconnect(); err != nil {
		return err
	}

	// Send an end signal to the sending goroutine.
	go func() {
		ctx.sendEndc <- struct{}{}
	}()

	// Send an end signal to the receiving goroutine.
	go func() {
		ctx.recvEndc <- struct{}{}
	}()

	// Send an end signal to the CONNACK monitoring goroutine.
	go func() {
		ctx.connackEndc <- struct{}{}
	}()

	// Wait for receiving an ended signal from each goroutine.
	<-ctx.sendEndedc
	<-ctx.readEndedc
	<-ctx.recvEndedc
	<-ctx.connackEndedc

	// Clear the Network Connection.
	ctx.cli.ClearConnection()

	return nil
}
