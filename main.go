package main

import (
	"fmt"

	"github.com/ing-bank/kal/internal/logger"
	"github.com/ing-bank/kal/internal/utilz"
	kalpkg "github.com/ing-bank/kal/pkg/kal"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Configure(rootCmd)
		log.Error().Err(err).Msg("failed to execute kal")
	}
}

var rootCmd = &cobra.Command{
	Use:   "kal",
	Short: "Kubernetes Authorization Listing - KAL",
	Long:  fmt.Sprintf("%s\n\n%s", utilz.BANNER, utilz.DISCLAIMER),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Configure logging approach
		logger.Configure(cmd)
	},
}

func init() {
	kalpkg.KALCmd.Use = "run"
	rootCmd.AddCommand(kalpkg.KALCmd)

	kubeConfigDir := fmt.Sprintf("%s/.kube/config", homedir.HomeDir())
	rootCmd.PersistentFlags().String("kubeconfig", kubeConfigDir, "Path to the kubeconfig file")

	rootCmd.PersistentFlags().BoolP("color", "c", true, "colored output")
	rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Skip TLS verification")
	rootCmd.PersistentFlags().BoolP("no-rate-limit", "L", false, "disable rate limiting")
	rootCmd.PersistentFlags().String("as", "", "Username to impersonate")
	rootCmd.PersistentFlags().String("token", "", "Bearer token for API authentication")
	rootCmd.PersistentFlags().StringP("namespace", "n", "default", "Namespace scope")
	rootCmd.PersistentFlags().StringP("server", "s", "", "Server URL")
	rootCmd.PersistentFlags().StringP("user-agent", "A", "KAET", "HTTP User-Agent header")
}
