package native

import (
	"context"
	"fmt"
	"os"

	"github.com/peers-touch/peers-touch/station/frame/core/client"
	"github.com/peers-touch/peers-touch/station/frame/core/config"
	"github.com/peers-touch/peers-touch/station/frame/core/logger"
	"github.com/peers-touch/peers-touch/station/frame/core/option"
	"github.com/peers-touch/peers-touch/station/frame/core/plugin"
	"github.com/peers-touch/peers-touch/station/frame/core/registry"
	"github.com/peers-touch/peers-touch/station/frame/core/server"
	"github.com/peers-touch/peers-touch/station/frame/core/util/log"
)

// Init initialises options. Additionally, it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *native) Init(ctx context.Context, opts ...option.Option) error {
	// process options
	for _, o := range opts {
		s.opts.Apply(o)
	}

	if len(s.opts.BeforeInit) > 0 {
		for _, f := range s.opts.BeforeInit {
			err := f(s.opts)
			if err != nil {
				log.Fatalf("init peers err: %s", err)
			}
		}
	}

	if s.opts.Name == "" {
		hostName, err := os.Hostname()
		if err != nil {
			logger.Errorf(ctx, "get hostname failed: %v", err)
			hostName = "unknown"
		}
		s.opts.Name = "peers-touch-go:s:name:" + hostName
	}

	if s.opts.Id == "" {
		s.opts.Id = "peers-touch-go:s:id:" + s.opts.Name
	}

	// begin init
	if err := s.initComponents(ctx); err != nil {
		log.Fatalf("init peers' components err: %s", err)
	}

	return nil
}

func (s *native) initComponents(ctx context.Context) error {
	// set Logger
	// only change if we have the logger and type differs
	if s.opts.Logger == nil {
		// use slogs as default logger
		loggerName := config.Get("peers.logger.name").String("slogrus")
		l, ok := plugin.LoggerPlugins[loggerName]
		if !ok {
			return fmt.Errorf("logger [%s] not found", loggerName)
		}

		s.opts.Logger = l.New()
	}
	if err := s.opts.Logger.Init(ctx, s.opts.LoggerOptions...); err != nil {
		return err
	}

	// init store
	if s.opts.Store == nil {
		storeName := config.Get("peers.store.name").String(plugin.NativePluginName)
		if len(storeName) > 0 {
			if plugin.StorePlugins[storeName] == nil {
				logger.Errorf(ctx, "store %s not found, use native by default", storeName)
				storeName = plugin.NativePluginName
			}
		}

		logger.Infof(ctx, "initial store's name is: %s", storeName)

		s.opts.Store = plugin.StorePlugins[storeName].New()
	}
	if err := s.opts.Store.Init(ctx, s.opts.StoreOptions...); err != nil {
		return err
	}

	// init registry
	if s.opts.Registry == nil {
		registryName := config.Get("peers.registry.name").String(plugin.NativePluginName)
		if len(registryName) > 0 {
			if plugin.RegistryPlugins[registryName] == nil {
				logger.Errorf(ctx, "registry %s not found, use native by default", registryName)
				registryName = plugin.NativePluginName
			}
		}

		logger.Infof(ctx, "initial registry's name is: %s", registryName)

		s.opts.Registry = plugin.RegistryPlugins[registryName].New(registry.WithStore(s.opts.Store))
	}

	// todo, set private key in advance.
	if err := s.opts.Registry.Init(ctx, append(s.opts.RegistryOptions, registry.WithPrivateKey(s.opts.PrivateKey))...); err != nil {
		return fmt.Errorf("init registry err: %v", err)
	}

	// init server
	if s.opts.Server == nil {
		serverName := config.Get("peers.node.server.name").String(plugin.NativePluginName)
		if len(serverName) > 0 {
			if plugin.ServerPlugins[serverName] == nil {
				logger.Errorf(ctx, "server %s not found, use native by default", serverName)
				serverName = plugin.NativePluginName
			}
		}

		logger.Infof(ctx, "initial server's name is: %s", serverName)
		s.opts.Server = plugin.ServerPlugins[serverName].New()

		// Inject plugin subservers into subServerNewFunctions
		for name, p := range plugin.SubserverPlugins {
			if p.Enabled() {
				s.opts.ServerOptions = append(s.opts.ServerOptions, server.WithSubServer(name, p.New))
			} else {
				logger.Infof(ctx, "subserver [%s] is disabled", name)
			}
		}
	}
	if err := s.opts.Server.Init(s.opts.ServerOptions...); err != nil {
		return err
	}

	// init client
	if s.opts.Client == nil {
		clientName := config.Get("peers.client.name").String(plugin.NativePluginName)
		if len(clientName) > 0 {
			if plugin.ClientPlugins[clientName] == nil {
				logger.Errorf(ctx, "client %s not found, use native by default", clientName)
				clientName = plugin.NativePluginName
			}
		}

		logger.Infof(ctx, "initial client's name is: %s", clientName)
		s.opts.Client = plugin.ClientPlugins[clientName].New()
	}

	if err := s.opts.Client.Init(append([]option.Option{client.Registry(s.opts.Registry)}, s.opts.ClientOptions...)...); err != nil {
		return err
	}

	return nil
}
