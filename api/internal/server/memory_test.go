package server

import (
	"context"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	api "github.com/atato/api/proto"
)

const (
	bufSize = 1024 * 1024
	expire  = 10
	cleanup = 1
)

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	api.RegisterCacheServiceServer(s, NewCacheService(time.Duration(expire)*time.Minute, time.Duration(cleanup)*time.Second))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSet(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	c := api.NewCacheServiceClient(conn)
	keyVal1 := &api.Item{
		Key:        "company",
		Value:      "atato",
		Expiration: "1m",
	}

	resp, err := c.Set(context.Background(), keyVal1)
	if err != nil {
		t.Fatalf("Adding key failed: %v", err)
	}
	if resp.Key != "company" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Key, "company")
	}
	if resp.Value != "atato" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Value, "atato")
	}

	// Checking for race condition
	for i := 0; i < 100; i++ {
		go c.Set(context.Background(), &api.Item{
			Key:        strconv.Itoa(i),
			Value:      "Value of i is ",
			Expiration: strconv.Itoa(i),
		})
	}

}

func TestDump(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	c := api.NewCacheServiceClient(conn)

	keyGet := &api.GetKey{
		Key: "company",
	}
	resp, err := c.Dump(context.Background(), keyGet)
	if err != nil {
		t.Fatalf("Getting key failed: %v", err)
	}
	if resp.Key != "company" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Key, "company")
	}
	if resp.Value != "atato" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Value, "atato")
	}
}

func TestIncrWithValidKey(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	c := api.NewCacheServiceClient(conn)
	keyVal1 := &api.Item{
		Key:        "var",
		Value:      "110",
		Expiration: "1m",
	}

	c.Set(context.Background(), keyVal1)

	keyGet := &api.GetKey{
		Key: "var",
	}

	_, err = c.Incr(context.Background(), keyGet)
	if err != nil {
		t.Fatalf("Increment key failed: %v", err)
	}

	resp, err := c.Dump(context.Background(), keyGet)
	if err != nil {
		t.Fatalf("Getting key failed: %v", err)
	}

	if resp.Key != "var" {
		t.Errorf("handler returned unexpected key: got %v want %v",
			resp.Key, "var")
	}

	if resp.Value != "111" {
		t.Errorf("handler returned unexpected value: got %v want %v",
			resp.Value, "111")
	}
}

func TestIncrWithKeyIsNotExist(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	c := api.NewCacheServiceClient(conn)

	keyGet := &api.GetKey{
		Key: "car",
	}

	_, err = c.Incr(context.Background(), keyGet)
	if err != nil {
		t.Fatalf("Increment key failed: %v", err)
	}

	resp, err := c.Dump(context.Background(), keyGet)
	if err != nil {
		t.Fatalf("Getting key failed: %v", err)
	}

	if resp.Key != "car" {
		t.Errorf("handler returned unexpected key: got %v want %v",
			resp.Key, "car")
	}
	if resp.Value != "0" {
		t.Errorf("handler returned unexpected value: got %v want %v",
			resp.Value, "0")
	}
}
