package commands

import (
	"context"
	"fmt"
	"github.com/claco/claco.pinochle/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"time"
)

type AddressArgs struct {
	Address string `arg:"--address,env:SERVICE_ADDRESS" default:"localhost" help:"service host address"`
	Port    uint16 `arg:"--port,env:SERVICE_PORT" default:"50051" help:"service host port"`
	Socket  string `arg:"--socket,env:SERVICE_SOCKET" default:"/tmp/pinochle/pinochle.sock" help:"service socket"`
}

type KeyArgs struct {
	Id string `arg:"positional,required:id" help:"id or slug of game"`
}

func (args *AddressArgs) NewClient() (pb.PinochleServiceClient, error) {
	var target = fmt.Sprintf("%s:%d", args.Address, args.Port)
	var timeout = 1 * time.Second
	var fields = log.Fields{
		"target":  target,
		"timeout": timeout,
	}

	log.WithFields(fields).Debugf("creating client: %s", target)

	var ctx, cancel = context.WithTimeout(context.Background(), timeout)
	var connection, err = grpc.DialContext(ctx, target, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	defer cancel()

	if err != nil {
		return nil, err
	}

	var client = pb.NewPinochleServiceClient(connection)

	return client, nil
}

func (args *AddressArgs) NewSocketClient() (pb.PinochleServiceClient, error) {
	var target = args.Socket
	var timeout = 1 * time.Second
	var fields = log.Fields{
		"target":  target,
		"timeout": timeout,
	}

	log.WithFields(fields).Debugf("creating client: %s", target)

	dialer := func(addr string, t time.Duration) (net.Conn, error) {
		return net.Dial("unix", addr)
	}

	//lint:ignore SA1019 no context yet
	connection, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDialer(dialer)) //nolint
	if err != nil {
		return nil, err
	}

	var client = pb.NewPinochleServiceClient(connection)

	return client, nil
}
