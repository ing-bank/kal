/*
KAL enumerates the authorization of an user or service account JWT token.
It uses the Kubernetes API to list all available resources and verify, for each resource,
if the token has permission to interact with the resource.

Usage:

	kal [flags]

Flags:
KUBERNETES:

	-token string          kubernetes api token
	-url string            kubernetes api base url
	-k, -insecure-tls      disable TLS verification
	-n, -namespace string  namespace name
	-nrl, -no-rate-limit   remove rate limit
	-as string             user/service account to impersonate
	-c, -config string     absolute path to kubeconfig file (default "$HOME/.kube/config")

OUTPUT:

	-v, -verbose       verbose output
	-s, -silent        silent output
	-sr, -show-reason  show reasons from kubernetes API response
	-all               show all results, including ones without verbs allowed
	-j, -json          output as json
	-nc, -no-color     no color output

When KAL is not provided an authentication configuration it searches for the `$HOME/.kube/config`
file. Otherwise, it uses the provided information via CLI arguments. If KAL is executed inside a
Kubernetes POD, it will use the data saved in the folder `/var/run/secrets/kubernetes.io/serviceaccount`.
*/
package main

import (
	"os"
	"os/signal"
	"path/filepath"

	"github.com/ing-bank/kal/pkg/runner"
	"github.com/ing-bank/kal/pkg/types"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var options *types.Options

func init() {
	options = &types.Options{
		Kubernetes: &types.KubernetesOptions{},
		Output:     &types.OutputOptions{},
	}
}

func main() {
	configureFlags()

	if options.Kubernetes.KubeConfigPath != "" {
		// read file and set options
		setServerUrlFromKubeConfig(options.Kubernetes.KubeConfigPath)
	}
	options.Validate()
	options.Configure()

	printBannerAndDisclaimer()

	run := runner.FromOptions(options)

	// Setup graceful exits
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			gologger.Info().Msgf("CTRL+C pressed: Exiting\n")
			run.Close()
			os.Exit(1)
		}
	}()

	run.Exec()
}

func configureFlags() {
	set := goflags.NewFlagSet()
	set.SetDescription(types.Banner)

	setGroup(set, "kubernetes", "kubernetes",
		set.StringVar(&options.Kubernetes.ApiToken, "token", "", "kubernetes api token"),
		set.StringVar(&options.Kubernetes.ServerURL, "url", "", "kubernetes api base url"),
		set.BoolVarP(&options.Kubernetes.InsecureTLS, "insecure-tls", "k", false, "disable TLS verification"),
		set.StringVarP(&options.Kubernetes.Namespace, "namespace", "n", "", "namespace name"), // support multiple namespaces
		set.BoolVarP(&options.Kubernetes.NoRateLimit, "no-rate-limit", "nrl", false, "remove rate limit"),
		set.StringVar(&options.Kubernetes.UserToImpersonate, "as", "", "user/service account to impersonate"),
	)

	if home := homedir.HomeDir(); home != "" {
		set.
			StringVarP(
				&options.Kubernetes.KubeConfigPath,
				"config",
				"c",
				filepath.Join(home, ".kube", "config"),
				"absolute path to kubeconfig file",
			).
			Group("kubernetes")
	} else {
		set.
			StringVarP(
				&options.Kubernetes.KubeConfigPath,
				"config",
				"c",
				"",
				"absolute path to kubeconfig file",
			).
			Group("kubernetes")
	}

	setGroup(set, "output", "output",
		set.BoolVarP(&options.Verbose, "verbose", "v", false, "verbose output"),
		set.BoolVarP(&options.Silent, "silent", "s", false, "silent output"),
		set.BoolVarP(&options.Output.ShowReason, "show-reason", "sr", false, "show reasons from kubernetes API response"),
		set.BoolVar(&options.Output.ShowAll, "all", false, "show all results, including ones without verbs allowed"),
		set.BoolVarP(&options.Output.JSON, "json", "j", false, "output as json"),
		set.BoolVarP(&options.Output.NoColor, "no-color", "nc", false, "no color output"),
	)

	_ = set.Parse()
}

func setGroup(set *goflags.FlagSet, groupName, description string, flags ...*goflags.FlagData) {
	set.SetGroup(groupName, description)
	for _, currentFlag := range flags {
		currentFlag.Group(groupName)
	}
}

func setServerUrlFromKubeConfig(configPath string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		gologger.Fatal().Msgf("could not parse kube config file. error: %s\n", err)
	}
	options.Kubernetes.ServerURL = config.Host

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		gologger.Fatal().Msgf("could not create kubernetes client. error: %s\n", err)
	}

	return client
}

func printBannerAndDisclaimer() {
	gologger.Silent().Msg(types.Banner + "\n")
	gologger.Silent().Msgf("[!] legal disclaimer: %s\n\n\n", types.Disclaimer)
}
