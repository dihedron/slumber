package gui

import (
	"context"
	"fmt"
	"time"

	"github.com/dihedron/slumber/command/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client api.VMControlClient
	userid string
}

func NewClient(address, userid string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}

	return &Client{
		conn:   conn,
		client: api.NewVMControlClient(conn),
		userid: userid,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Status(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.client.Status(ctx, &api.StatusRequest{UserId: c.userid})
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}

func (c *Client) Start(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := c.client.Start(ctx, &api.StartRequest{UserId: c.userid})
	return err
}

func (c *Client) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := c.client.Stop(ctx, &api.StopRequest{UserId: c.userid})
	return err
}

func (c *Client) Pause(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := c.client.Pause(ctx, &api.PauseRequest{UserId: c.userid})
	return err
}

func (c *Client) Unpause(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := c.client.Unpause(ctx, &api.UnpauseRequest{UserId: c.userid})
	return err
}
