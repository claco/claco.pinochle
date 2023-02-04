package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/claco/claco.pinochle/build"
	"github.com/claco/claco.pinochle/pb"
	"github.com/claco/claco.pinochle/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type ServiceArgs struct {
	AddressArgs

	CheckArgs *CheckArgs `arg:"subcommand:check" help:"check service health"`
	RunArgs   *RunArgs   `arg:"subcommand:run" help:"run service interactively"`
}

type CheckArgs struct{}

type RunArgs struct {
	Migrate bool `arg:"--migrate" help:"apply database migrations"`
}

func (sargs *ServiceArgs) Execute() (exitCode int, err error) {
	if sargs.CheckArgs != nil {
		return sargs.Check(sargs.CheckArgs)
	} else if sargs.RunArgs != nil {
		return sargs.Run(sargs.RunArgs)
	}

	return 127, errors.New("command not implemented: service")
}

func (sargs *ServiceArgs) Check(args *CheckArgs) (exitCode int, err error) {
	request := &grpc_health_v1.HealthCheckRequest{}

	dialer := func(addr string, t time.Duration) (net.Conn, error) {
		return net.Dial("unix", addr)
	}

	//lint:ignore SA1019 no context yet
	connection, err := grpc.Dial(sargs.Socket, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDialer(dialer))
	if err != nil {
		return 1, err
	}

	client := grpc_health_v1.NewHealthClient(connection)
	if err != nil {
		return 1, err
	}

	response, err := client.Check(context.Background(), request)
	if err != nil {
		return 1, err
	} else {
		log.Info(response.Status)
	}

	return 0, nil
}

func (sargs *ServiceArgs) Run(args *RunArgs) (exitCode int, err error) {
	socketPath := sargs.Socket
	serviceAddress := fmt.Sprintf("%s:%d", sargs.AddressArgs.Address, sargs.Port)
	tcpAddress, err := net.ResolveTCPAddr("tcp", serviceAddress)
	fields := log.Fields{
		"address":  sargs.Address,
		"port":     sargs.Port,
		"version":  build.GetVersion(),
		"listener": tcpAddress.String(),
		"socket":   socketPath,
	}

	log.WithFields(fields).Infof("running service: tcp://%s unix:%s", tcpAddress, socketPath)

	if err != nil {
		return 1, err
	}

	service := service.NewPinochleService()
	service.ConnectDatabase()

	if args.Migrate {
		if err := service.InitializeDatabase(); err != nil {
			return 1, err
		}
	}

	log.WithField("socket", socketPath).Debugf("creating socket: %s", socketPath)

	socketFolder := filepath.Dir(socketPath)
	os.Mkdir(socketFolder, os.FileMode(0755))

	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	log.WithField("listener", tcpAddress).Debugf("creating listener: %s", tcpAddress)
	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		return 1, err
	}
	defer listener.Close()

	server := grpc.NewServer()

	grpc_health_v1.RegisterHealthServer(server, service)
	pb.RegisterPinochleServiceServer(server, service)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c

		log.Info("stopping service")

		os.RemoveAll(socketFolder)
		server.GracefulStop()
	}()

	go func() {
		err := server.Serve(socket)
		if err != nil {
			log.Error(err)
		}
	}()

	err = server.Serve(listener)
	if err != nil {
		log.Error(err)
	}

	return 0, nil
}
