package main

import (
	"context"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/atato/api/internal/util"
	api "github.com/atato/api/proto"
)

const (
	OpTypeSet  = "Set"
	OpTypeDump = "Dump"
	OpTypeIncr = "Incr"
)

var (
	address string
	conn    *grpc.ClientConn
	err     error
	key     string
	value   string
	opType  string
)

func init() {
	const (
		addressUsage   = "Server connection address"
		addressDefault = "127.0.0.1:12345"

		keyUsage   = "Key"
		keyDefault = ""

		valueUsage   = "Value"
		valueDefault = ""

		opTypeUsage = "Operation type: " +
			OpTypeSet +
			"|" + OpTypeDump +
			"|" + OpTypeIncr
		opTypeDefault = ""
	)

	flag.StringVar(&address, "addr", addressDefault, addressUsage)
	flag.StringVar(&key, "k", keyDefault, keyUsage)
	flag.StringVar(&value, "v", valueDefault, valueUsage)
	flag.StringVar(&opType, "o", opType, opTypeUsage)
}

func checkArgs() {
	if key == "" {
		panic("No key")
	}
}

func logArgs() {
	log.Info("App mockredis-cli start...")
	log.Info("Server connection address: ", address)
}

func opToStr() (string, error) {
	switch opType {
	case OpTypeSet:
		return "OpTypeSet", nil
	case OpTypeDump:
		return "OpTypeDump", nil
	case OpTypeIncr:
		return "OpTypeIncr", nil
	default:
		return "", fmt.Errorf("nknown operation")
	}
}

func main() {
	defer util.RecoverAtStartup()

	flag.Parse()
	checkArgs()
	logArgs()

	conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		util.PanicOnError(err)
	}
	defer conn.Close()

	c := api.NewCacheServiceClient(conn)

	op, err := opToStr()
	util.PanicOnError(err)
	log.Printf("Operation in progress '%s'", op)

	switch opType {
	case OpTypeSet:
		err = set(c)
	case OpTypeDump:
		err = dump(c)
	case OpTypeIncr:
		err = incr(c)
	}
	util.PanicOnError(err)
}

func set(c api.CacheServiceClient) error {
	log.Info("Call method SET")

	// Add key
	i := &api.Item{
		Key:        key,
		Value:      value,
		Expiration: "1m",
	}

	resp, err := c.Set(context.Background(), i)
	if err != nil {
		return err
	}
	log.Info("Response from server: ", resp)

	return nil
}

func dump(c api.CacheServiceClient) error {
	log.Info("Call method DUMP")

	resp, err := c.Dump(context.Background(), &api.GetKey{
		Key: key,
	})
	if err != nil {
		return err
	}
	log.Info("Response from server: ", resp)

	return nil
}

func incr(c api.CacheServiceClient) error {
	log.Info("Call method INCR")

	resp, err := c.Incr(context.Background(), &api.GetKey{
		Key: key,
	})
	if err != nil {
		return err
	}
	log.Info("Response from server: ", resp)

	return nil
}
