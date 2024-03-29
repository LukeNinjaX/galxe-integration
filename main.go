package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/artela-network/galxe-integration/api"
	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/db"
	"github.com/artela-network/galxe-integration/logging"
	_ "github.com/artela-network/galxe-integration/logging"
	"github.com/artela-network/galxe-integration/onchain/faucet"
	"github.com/artela-network/galxe-integration/onchain/rug"
	"github.com/artela-network/galxe-integration/onchain/updater"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	serviceConf := flag.String("config", "./config.json", "monitor config json file path")
	flag.Parse()

	if !*debug {
		logging.SetupLogRotate()
	} else {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug mode enabled")
	}

	log.Info("starting the service...")

	ctx, cancel := context.WithCancel(context.Background())

	conf := loadConfig(*serviceConf)
	biz.GoPlus_Config = conf.GoPlus
	biz.Recaptcha_Config = conf.Recaptcha

	conn, driver, err := db.GetDB(ctx, conf.DB)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// indexers := make([]common.Indexer, len(conf.Indexers))
	// chainFetcher, err := fetcher.NewFetcher(ctx, conf.Fetcher, driver, conn)
	// if err != nil {
	// 	log.Fatalf("failed to create fetcher: %v", err)
	// }
	// for i, indexerConf := range conf.Indexers {
	// 	indexerInstance, err := indexer.GetRegistry().GetIndexer(ctx, indexerConf, driver, conn)
	// 	if err != nil {
	// 		log.Fatalf("failed to create indexer: %v", err)
	// 	}
	// 	chainFetcher.RegisterIndexer(indexerInstance)
	// 	indexers[i] = indexerInstance
	// }
	// chainFetcher.Start()

	apiServer := api.NewServer(ctx, conf, driver, conn, nil, nil)
	apiServer.Start()

	rugServ, err := rug.NewRug(conn, conf.Rug)
	if err != nil {
		log.Error("failed to start rug service", err)
		os.Exit(-1)
	}
	rugServ.Start()

	faucetServ, err := faucet.NewFaucet(conn, conf.Faucet)
	if err != nil {
		log.Error("failed to start faucet service", err)
		os.Exit(-1)
	}
	faucetServ.Start()

	updaterServ, err := updater.NewUpdater(conn, conf.Updater)
	if err != nil {
		log.Error("failed to start updater service", err)
		os.Exit(-1)
	}
	updaterServ.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGINT)

	<-c

	apiServer.Stop()

	cancel()

	log.Info("service exited")
}

func loadConfig(configFile string) *config.Config {
	// load config file
	if strings.HasPrefix(configFile, ".") {
		cwd, _ := os.Getwd()
		configFile = path.Join(cwd, configFile)
	}

	log.Infof("loading config file: %s ...", configFile)
	f, err := os.Open(configFile)
	if err != nil {
		log.Panicf("failed to load config: %v", err)
	}
	defer f.Close()

	rawConf, err := io.ReadAll(f)
	if err != nil {
		log.Panicf("failed to read config: %v", err)
	}

	conf := &config.Config{}
	if err := json.Unmarshal(rawConf, conf); err != nil {
		log.Panicf("failed to parse config: %v", err)
	}

	return conf
}
