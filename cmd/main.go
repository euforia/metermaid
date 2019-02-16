package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hexablock/iputil"

	"go.uber.org/zap"

	"github.com/euforia/base58"
	"github.com/euforia/gossip"
	"github.com/euforia/metermaid"
	"github.com/euforia/metermaid/storage"
)

var (
	bindAddr = flag.String("bind-addr", "127.0.0.1:8080", "")
	advAddr  = flag.String("adv-addr", "", "")
	joinPeer = flag.String("join", "", "")
)

func init() {
	flag.Parse()
}

func initGossip(logger *zap.Logger, node *metermaid.Node) (*gossip.Gossip, *gossip.Pool) {
	gconf := gossip.DefaultConfig()

	sh := sha256.Sum256([]byte(*advAddr))
	gconf.Name = string(base58.Encode(sh[:]))

	gconf.BindAddr, gconf.BindPort, _ = iputil.SplitHostPort(*bindAddr)

	advHost, advPort, err := iputil.BuildAdvertiseAddr(*advAddr, *bindAddr)
	if err != nil {
		logger.Fatal("failed to get advertise address", zap.Error(err))
	}
	gconf.AdvertiseAddr = advHost
	gconf.AdvertisePort, _ = iputil.PortFromString(advPort)

	node.Name = gconf.Name
	node.Address = advHost + ":" + advPort

	gsp, err := gossip.New(gconf)

	pconf := gossip.DefaultLANPoolConfig(222)
	gspDel := &GossipDelegate{log: logger, node: *node}
	pconf.Delegate = gspDel
	pconf.Memberlist.Events = gspDel
	gpool := gsp.RegisterPool(pconf)
	if err = gsp.Start(); err != nil {
		logger.Fatal("failed to start gossip", zap.Error(err))
	}

	if *joinPeer != "" {
		var peers []string
		// Check if we should use service discovery to find the peer
		if _, _, err = iputil.SplitHostPort(*joinPeer); err != nil {
			peers, err = getAddrViaSD(*joinPeer)
			if err != nil {
				logger.Fatal("failed to get addresses", zap.Error(err))
			}
		} else {
			peers = []string{*joinPeer}
		}

		if _, err = gpool.Join(peers); err != nil {
			logger.Info("failed to join peer", zap.Error(err))
		}
	}
	return gsp, gpool
}

func getAddrViaSD(name string) ([]string, error) {
	_, addrs, err := net.DefaultResolver.LookupSRV(context.Background(), "", "", name)
	out := make([]string, len(addrs))
	if err == nil {
		for i, addr := range addrs {
			out[i] = fmt.Sprintf("%s:%d", strings.TrimSuffix(addr.Target, "."), addr.Port)
		}
	}
	return out, err
}

func main() {

	logger, _ := zap.NewDevelopment()
	node := metermaid.NewNode()
	logger.Info("node stats",
		zap.Uint64("cpu", node.CPUShares),
		zap.Uint64("memory", node.Memory),
	)

	mm, err := metermaid.New(logger)
	if err != nil {
		logger.Fatal("failed to initialize metermaid", zap.Error(err))
	}

	gsp, gpool := initGossip(logger, node)

	contStore := storage.NewInmemContainers()

	// Update state
	go func() {
		updates := mm.Updates()
		for update := range updates {
			contStore.Set(update)
			logger.Info("update",
				zap.String("id", update.ID),
				zap.Duration("runtime", update.RunTime()),
				zap.Duration("alloctime", update.AllocatedTime()),
			)
		}
		logger.Info("update loop exiting")
	}()

	capi := &containerAPI{"/container", node, contStore}
	http.Handle("/container/", capi)
	napi := &nodeAPI{"/node", gpool}
	http.Handle("/node/", napi)
	http.HandleFunc("/", handleUI)

	go func(ln net.Listener) {
		logger.Info("http server", zap.String("address", *bindAddr))
		err := http.Serve(ln, nil)
		if err != nil {
			logger.Info("http shutdown unclean", zap.Error(err))
		}
	}(gsp.ListenTCP())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	mm.Stop()
}
