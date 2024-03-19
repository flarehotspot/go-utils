/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package sdkconnmgr

import "context"

// ClientDevice represents a client device connected to the network.
type ClientDevice interface {
	// Returns the database id of the device.
	Id() int64

	// Returns the hostname of the device.
	Hostname() string

	// Returns the IP address of the device.
	IpAddr() string

	// Returns the MAC address of the device.
	MacAddr() string

	// Updates the client device.
	Update(ctx context.Context, mac string, ip string, hostname string) error

	// Emits a socket event to a client device.
	// The event will be propagated to the client's browser via server-sent events.
	Emit(event string, data any)

	// Subscribes to a socket event.
    // It returns a channel that will receive data when the event is emitted.
	Subscribe(event string) <-chan []byte
}
