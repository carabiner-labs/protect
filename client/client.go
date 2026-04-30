// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

// Package client provides a thin gRPC client wrapper around the
// Edera Protect ControlService defined in the open source
// proto definitions https://github.com/edera-dev/protos.
package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	controlv1 "github.com/carabiner-labs/protect/gen/protect/control/v1"
)

// Client is a connected ControlService client. It owns the underlying
// gRPC connection; call Close when done.
type Client struct {
	Control controlv1.ControlServiceClient

	conn *grpc.ClientConn
}

// Option configures a Client at construction time.
type Option func(*options)

type options struct {
	tlsConfig   *tls.Config
	insecure    bool
	dialOptions []grpc.DialOption
}

// WithTLS uses the given TLS configuration for transport credentials.
// If both WithTLS and WithInsecure are set, WithTLS wins.
func WithTLS(cfg *tls.Config) Option {
	return func(o *options) { o.tlsConfig = cfg }
}

// WithInsecure dials without transport security. Intended for local
// development against a daemon on a trusted network or unix socket.
func WithInsecure() Option {
	return func(o *options) { o.insecure = true }
}

// WithDialOption appends arbitrary grpc.DialOption values. Useful for
// interceptors, keepalives, or custom resolvers.
func WithDialOption(opts ...grpc.DialOption) Option {
	return func(o *options) { o.dialOptions = append(o.dialOptions, opts...) }
}

// Dial connects to a Protect ControlService at target. By default it
// uses TLS with the system root CAs; pass WithInsecure or WithTLS to
// override.
func Dial(ctx context.Context, target string, opts ...Option) (*Client, error) {
	if target == "" {
		return nil, errors.New("protect/client: target must not be empty")
	}

	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	dialOpts := make([]grpc.DialOption, 0, len(o.dialOptions)+1)
	switch {
	case o.tlsConfig != nil:
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(o.tlsConfig)))
	case o.insecure:
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	default:
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})))
	}
	dialOpts = append(dialOpts, o.dialOptions...)

	conn, err := grpc.NewClient(target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("protect/client: dial %s: %w", target, err)
	}

	return &Client{
		Control: controlv1.NewControlServiceClient(conn),
		conn:    conn,
	}, nil
}

// Conn returns the underlying gRPC connection. Useful for advanced
// callers that want to construct additional service clients sharing
// the same channel.
func (c *Client) Conn() *grpc.ClientConn { return c.conn }

// Close releases the underlying gRPC connection.
func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
