package main

import (
	"context"
	"sync"

	goutils "go.viam.com/utils"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/datamanager"
)

var Model = resource.NewModel("ncs", "sensor", "nick-selective-sync-sensor")
var (
	Reset = "\033[0m"
	Green = "\033[32m"
	Cyan  = "\033[36m"
)

type fake struct {
	mu   sync.Mutex
	sync bool
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger
}

type Config struct {
	Sync bool `json:"sync,omitempty"`
}

func (c *Config) Validate(path string) ([]string, error) {
	return nil, nil
}

func newSensor(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger logging.Logger,
) (sensor.Sensor, error) {
	c, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return nil, err
	}

	return &fake{
		Named:  conf.ResourceName().AsNamed(),
		sync:   c.Sync,
		logger: logger,
	}, nil
}

func (f *fake) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return map[string]interface{}{
		datamanager.ShouldSyncKey: f.sync,
	}, nil
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	resource.RegisterComponent(
		sensor.API,
		Model,
		resource.Registration[sensor.Sensor, *Config]{Constructor: newSensor})

	module, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}
	if err := module.AddModelFromRegistry(ctx, sensor.API, Model); err != nil {
		return err
	}

	err = module.Start(ctx)
	defer module.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func main() {
	goutils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs(Model.String()))
}
