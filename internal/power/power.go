package power

import (
	"fmt"
	"log/slog"

	"github.com/godbus/dbus/v5"
)

const (
	dbusDest      = "org.freedesktop.login1"
	dbusPath      = "/org/freedesktop/login1"
	dbusInterface = "org.freedesktop.login1.Manager"
)

func Shutdown() error {
	slog.Info("requesting system shutdown")
	return callLogind("PowerOff")
}

func Hibernate() error {
	slog.Info("requesting system hibernation")
	return callLogind("Hibernate")
}

func callLogind(method string) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return fmt.Errorf("failed to connect to system bus: %w", err)
	}
	defer conn.Close()

	obj := conn.Object(dbusDest, dbus.ObjectPath(dbusPath))
	// The boolean argument is for "interactive" (polkit dialog)
	call := obj.Call(dbusInterface+"."+method, 0, true)
	if call.Err != nil {
		return fmt.Errorf("dbus call to %s failed: %w", method, call.Err)
	}
	return nil
}
