package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bpineau/kube-alert/config"
	klog "github.com/bpineau/kube-alert/pkg/log"
	"github.com/bpineau/kube-alert/pkg/run"
)

const appName = "kube-alert"

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
	msgPrefix string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   appName,
		Short: "Monitor pods",
		Long:  "Monitor pods and alert on failure",

		Run: func(cmd *cobra.Command, args []string) {
			config := &config.AlertConfig{
				DryRun:     viper.GetBool("dry-run"),
				Logger:     klog.New(viper.GetString("log.level"), viper.GetString("log.server"), viper.GetString("log.output")),
				DdAppKey:   viper.GetString("datadog.app-key"),
				DdApiKey:   viper.GetString("datadog.api-key"),
				HealthPort: viper.GetInt("healthcheck-port"),
				MsgPrefix:  viper.GetString("messages-prefix"),
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

	defaultCfg := "/etc/" + appName + "/" + appName + ".yaml"
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", defaultCfg, "configuration file")

	rootCmd.PersistentFlags().StringVarP(&apiServer, "api-server", "s", "", "kube api server url")
	if err := viper.BindPFlag("api-server", rootCmd.PersistentFlags().Lookup("api-server")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&kubeConf, "kube-config", "k", "", "kube config path")
	if err := viper.BindPFlag("kube-config", rootCmd.PersistentFlags().Lookup("kube-config")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}
	if err := viper.BindEnv("kube-config", "KUBECONFIG"); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry-run mode")
	if err := viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "v", "debug", "log level")
	if err := viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&logOutput, "log-output", "l", "stderr", "log output")
	if err := viper.BindPFlag("log.output", rootCmd.PersistentFlags().Lookup("log-output")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&logServer, "log-server", "r", "", "log server (if using syslog)")
	if err := viper.BindPFlag("log.server", rootCmd.PersistentFlags().Lookup("log-server")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&ddAppKey, "datadog-app-key", "a", "", "datadog app key")
	if err := viper.BindPFlag("datadog.app-key", rootCmd.PersistentFlags().Lookup("datadog-app-key")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&ddApiKey, "datadog-api-key", "i", "", "datadog api key")
	if err := viper.BindPFlag("datadog.api-key", rootCmd.PersistentFlags().Lookup("datadog-api-key")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().IntVarP(&healthP, "healthcheck-port", "p", 0, "port for answering healthchecks")
	if err := viper.BindPFlag("healthcheck-port", rootCmd.PersistentFlags().Lookup("healthcheck-port")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}

	rootCmd.PersistentFlags().StringVarP(&msgPrefix, "messages-prefix", "m", "", "prefix appended to notifications")
	if err := viper.BindPFlag("messages-prefix", rootCmd.PersistentFlags().Lookup("messages-prefix")); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName(appName)

	// all possible config file paths, by priority
	viper.AddConfigPath("/etc/" + appName + "/")
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
