package configfx

import (
	"fmt"

	"github.com/raf555/kbbi-api/internal/config"
	salomeconfig "github.com/raf555/salome/config/v1"
	"github.com/raf555/salome/config/v1/providers/infisical"
	"github.com/raf555/salome/config/v1/providers/os"
	"github.com/raf555/salome/config/v1/providers/osdotenv"
	"go.uber.org/fx"
)

type infisicalConfig struct {
	SiteUrl     string `env:"INFISICAL_SITE_URL"`
	IdentityID  string `env:"INFISICAL_KUBERNETES_IDENTITY_ID" validate:"required_with=SiteUrl"`
	ProjectSlug string `env:"INFISICAL_PROJECT_SLUG" validate:"required_with=SiteUrl"`
	Environment string `env:"INFISICAL_ENVIRONMENT" validate:"required_with=SiteUrl"`
}

var Module = fx.Module("config",
	fx.Provide(
		fx.Annotate(
			os.New,
			fx.As(fx.Self()),
			fx.As(new(salomeconfig.Provider)),
			fx.ResultTags(`name:"config_provider.os"`),
		),
	),

	fx.Provide(
		fx.Annotate(
			salomeconfig.LoadConfigTo[infisicalConfig],
			fx.ParamTags(`name:"config_provider.os"`),
		),
	),

	fx.Provide(
		func(infisicalCfg infisicalConfig, lc fx.Lifecycle) (salomeconfig.Provider, error) {
			if infisicalCfg.SiteUrl != "" { // load from cloud if provided
				infCfg, err := infisical.NewWithOptions(infisicalCfg.SiteUrl, infisical.SecretConfig{
					ProjectSlug: infisicalCfg.ProjectSlug,
					Environment: infisicalCfg.Environment,
					ConfigPath:  "",
				},
					infisical.WithKubernetesAuth(infisicalCfg.IdentityID, ""),
				)
				if err != nil {
					return nil, fmt.Errorf("infisical new: %w", err)
				}

				lc.Append(fx.StopHook(func() {
					infCfg.Close()
				}))

				return infCfg, nil
			}

			// local env
			defaultCfg, err := osdotenv.New(".env")
			if err != nil {
				return nil, fmt.Errorf("osdotenv.New: %w", err)
			}

			return defaultCfg, nil
		},
	),

	fx.Provide(
		fx.Annotate(
			func(provider salomeconfig.Provider) (*salomeconfig.Dynamic, error) {
				return salomeconfig.NewDynamic(provider)
			},
			fx.As(fx.Self()),
			fx.As(new(salomeconfig.DynamicConfigManager)),
			fx.OnStart(func(mgr *salomeconfig.Dynamic) {
				mgr.Start()
			}),
			fx.OnStop(func(mgr *salomeconfig.Dynamic) {
				mgr.Close()
			}),
		),
	),

	fx.Provide(salomeconfig.LoadConfigTo[config.ServerConfig]),
)
