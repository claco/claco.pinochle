package main

import (
	"bytes"
	"os"
	"strings"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"

	"github.com/alexflint/go-arg"
	"github.com/claco/claco.pinochle/build"
	"github.com/claco/claco.pinochle/commands"
	"golang.org/x/exp/slices"
)

type mainArgs struct {
	DebugArg bool `arg:"--debug" help:"show debug information"`
	JsonArg  bool `arg:"--json" help:"output json instead of plain text"`
	LogsArg  bool `arg:"--logs" help:"output logs instead of plain text"`

	GamesArgs   *commands.GamesArgs   `arg:"subcommand:games" help:"game related commands"`
	ServiceArgs *commands.ServiceArgs `arg:"subcommand:service" help:"service related commands"`
}

func main() {
	var args mainArgs = mainArgs{}
	var code int = 0
	var err error = nil

	configureLogging()

	var parser = parseArguments(&args)

	if args.ServiceArgs != nil {
		code, err = args.ServiceArgs.Execute(parser)
	} else if args.GamesArgs != nil {
		code, err = args.GamesArgs.Execute(parser)
	}

	if err != nil {
		log.Error(err)
	}

	log.Exit(code)
}

func configureLogging() {
	log.SetFormatter(&easy.Formatter{LogFormat: "%msg%\n"})

	if slices.Contains(os.Args, "--debug") || strings.ToUpper(os.Getenv("LOG_LEVEL")) == "DEBUG" {
		log.SetLevel(log.DebugLevel)
		// log.SetReportCaller(true)
	}

	if slices.Contains(os.Args, "--json") || strings.ToLower(os.Getenv("LOG_FORMAT")) == "json" {
		log.SetFormatter(&log.JSONFormatter{DataKey: "detail", DisableHTMLEscape: true})
	} else if slices.Contains(os.Args, "--logs") || strings.ToLower(os.Getenv("LOG_FORMAT")) == "logs" {
		log.SetFormatter(&nested.Formatter{NoColors: true, ShowFullLevel: true, TimestampFormat: time.StampMilli})
	}
}

func parseArguments(args *mainArgs) *arg.Parser {
	parser, _ := arg.NewParser(arg.Config{}, args)
	err := parser.Parse(os.Args[1:])
	help := bytes.NewBufferString("")
	usage := bytes.NewBufferString("")

	parser.WriteHelp(help)
	parser.WriteUsage(usage)

	if err == arg.ErrHelp {
		log.Info(help)
		log.Exit(0)
	} else if err == arg.ErrVersion {
		log.Info(build.GetVersion())
	} else if err != nil {
		log.Error(err)
		log.Exit(1)
	} else if parser.Subcommand() == nil {
		log.Warn(strings.Replace(usage.String(), "\n", "", 1))
		log.Exit(1)
	}

	return parser
}
