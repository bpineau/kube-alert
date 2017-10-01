package cmd

import (
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/log"
	"github.com/bpineau/kube-alert/pkg/run"
)

const AppName = "kube-alert"

var (
	cfgFile   string
	apiServer string
	kubeConf  string
	dryRun    bool
	logLevel  string
	logOutput string
	logServer string
	ddApiKey  string
	ddAppKey  string
	healthP   int

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   AppName,
		Short: "Monitor pods",
		Long:  "Monitor pods and alert on failure",

		Run: func(cmd *cobra.Command, args []string) {
			config := &config.AlertConfig{
				DryRun:     viper.GetBool("dry-run"),
				Logger:     log.New(viper.GetString("log.level"), viper.GetString("log.server"), viper.GetString("log.output")),
				DdAppKey:   viper.GetString("datadog.app-key"),
				DdApiKey:   viper.GetString("datadog.api-key"),
				HealthPort: viper.GetInt("healthcheck-port"),
			}
			config.Init(viper.GetString("api-server"), viper.GetString("kube-config"))
			run.Run(config)
		},
	}
)

// Execute adds all child commands to the root command and sets their flags.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	defaultCfg := "/etc/" + AppName + "/" + AppName + ".yaml"
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", defaultCfg, "configuration file")

	rootCmd.PersistentFlags().StringVarP(&apiServer, "api-server", "s", "", "kube api server url")
	viper.BindPFlag("api-server", rootCmd.PersistentFlags().Lookup("api-server"))

	rootCmd.PersistentFlags().StringVarP(&kubeConf, "kube-config", "k", "", "kube config path")
	viper.BindPFlag("kube-config", rootCmd.PersistentFlags().Lookup("kube-config"))
	viper.BindEnv("kube-config", "KUBECONFIG")

	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry-run mode")
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "v", "debug", "log level")
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().StringVarP(&logOutput, "log-output", "l", "stderr", "log output")
	viper.BindPFlag("log.output", rootCmd.PersistentFlags().Lookup("log-output"))

	rootCmd.PersistentFlags().StringVarP(&logServer, "log-server", "r", "", "log server (if using syslog)")
	viper.BindPFlag("log.server", rootCmd.PersistentFlags().Lookup("log-server"))

	rootCmd.PersistentFlags().StringVarP(&ddAppKey, "datadog-app-key", "a", "", "datadog app key")
	viper.BindPFlag("datadog.app-key", rootCmd.PersistentFlags().Lookup("datadog-app-key"))

	rootCmd.PersistentFlags().StringVarP(&ddApiKey, "datadog-api-key", "i", "", "datadog api key")
	viper.BindPFlag("datadog.api-key", rootCmd.PersistentFlags().Lookup("datadog-api-key"))

	rootCmd.PersistentFlags().IntVarP(&healthP, "healthcheck-port", "p", 0, "port for answering healthchecks")
	viper.BindPFlag("healthcheck-port", rootCmd.PersistentFlags().Lookup("healthcheck-port"))
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName(AppName)

	// all possible config file paths, by priority
	viper.AddConfigPath("/etc/" + AppName + "/")
	if home, err := homedir.Dir(); err == nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(".")

	// prefer the config file path provided by cli flag, if any
	if _, err := os.Stat(cfgFile); !os.IsNotExist(err) {
		viper.SetConfigFile(cfgFile)
	}

	// allow config params through prefixed env variables
	viper.SetEnvPrefix("KUBE_ALERT")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file: ", viper.ConfigFileUsed())
	}
}
