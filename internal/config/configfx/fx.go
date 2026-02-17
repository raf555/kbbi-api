package configfx

import (
	"fmt"

	config "github.com/raf555/salome/config/v1"
	"github.com/raf555/salome/config/v1/providers/infisical"
	"github.com/raf555/salome/config/v1/providers/os"
	"github.com/raf555/salome/config/v1/providers/osdotenv"
	"go.uber.org/fx"
)

type infisicalConfig struct {
	SiteUrl      string `env:"INFISICAL_SITE_URL"`
	ClientID     string `env:"INFISICAL_CLIENT_ID" validate:"required_with=SiteUrl"`
	ClientSecret string `env:"INFISICAL_CLIENT_SECRET" validate:"required_with=SiteUrl"`
	ProjectSlug  string `env:"INFISICAL_PROJECT_SLUG" validate:"required_with=SiteUrl"`
	Environment  string `env:"INFISICAL_ENVIRONMENT" validate:"required_with=SiteUrl"`
}

var Module = fx.Module("config",
	fx.Provide(
		fx.Annotate(
			os.New,
			fx.As(fx.Self()),
			fx.As(new(config.Provider)),
			fx.ResultTags(`name:"config_provider.os"`),
		),
	),

	fx.Provide(
		fx.Annotate(
			config.LoadConfigTo[infisicalConfig],
			fx.ParamTags(`name:"config_provider.os"`),
		),
	),

	fx.Provide(
		func(infisicalCfg infisicalConfig, lc fx.Lifecycle) (config.Provider, error) {
			if infisicalCfg.SiteUrl != "" { // load from cloud if provided
				infCfg, err := infisical.New(infisical.Config{
					SiteUrl:      infisicalCfg.SiteUrl,
					ClientID:     infisicalCfg.ClientID,
					ClientSecret: infisicalCfg.ClientSecret,
					ProjectSlug:  infisicalCfg.ProjectSlug,
					Environment:  infisicalCfg.Environment,
				})
				if err != nil {
					return nil, fmt.Errorf("infisical.New: %w", err)
				}

				lc.Append(fx.StopHook(func() {
					infCfg.Close()
				}))

				return infCfg, nil
			}

			defaultCfg, err := osdotenv.New(".env")
			if err != nil {
				return nil, fmt.Errorf("osdotenv.New: %w", err)
			}

			return defaultCfg, nil
		},
	),

	fx.Provide(
		fx.Annotate(
			func(provider config.Provider) (*config.Dynamic, error) {
				return config.NewDynamic(provider)
			},
			fx.As(fx.Self()),
			fx.As(new(config.DynamicConfigManager)),
			fx.OnStart(func(mgr *config.Dynamic) {
				mgr.Start()
			}),
			fx.OnStop(func(mgr *config.Dynamic) {
				mgr.Close()
			}),
		),
	),
)
