package runner

import (
	"context"
	"errors"
	"sync"

	myK8s "github.com/ing-bank/kal/pkg/kubernetes"
	"github.com/ing-bank/kal/pkg/types"
	"github.com/projectdiscovery/gologger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
)

// FromOptions creates a KAL runner based on provided options
func FromOptions(o *types.Options) *Runner {
	r := &Runner{
		Context:    context.Background(),
		JSONOutput: o.Output.JSON,
		Namespace:  o.Kubernetes.Namespace,
		ShowAll:    o.Output.ShowAll,
		ShowReason: o.Output.ShowReason,
		outputChan: make(chan *Result),
		outputWg:   sync.WaitGroup{},
	}
	types.InitAurora(o)

	var client *kubernetes.Clientset
	var err error

	client, err = getCustomClient(o)
	if err != nil {
		gologger.Error().
			Str("error", err.Error()).
			Msg("could not create a kubernetes custom client\n")
	}

	if client == nil {
		client, err = getKubeConfigClient(o)
		if err != nil {
			gologger.Error().
				Str("error", err.Error()).
				Msg("could not create a kubernetes kubeconfig client\n")
		}
	}

	if client == nil {
		client, err = getInPodClient(o)
		if err != nil {
			gologger.Error().
				Str("error", err.Error()).
				Msg("could not create an in pod kubernetes client\n")
		}
	}

	if client == nil {
		gologger.Fatal().
			Str("error", "invalid kubernetes options. could not create a client")
	}

	r.KubernetesClient = client

	if r.Namespace == "" && o.Kubernetes.ApiToken != "" {
		r.Namespace = myK8s.GrabNamespaceFromToken(o.Kubernetes.ApiToken)
	}

	return r
}

func getCustomClient(o *types.Options) (*kubernetes.Clientset, error) {
	if o.Kubernetes == nil {
		return nil, errors.New("invalid kubernetes options")
	}

	if o.Kubernetes.ServerURL == "" || o.Kubernetes.ApiToken == "" {
		return nil, errors.New("invalid configuration for kubernetes custom client")
	}

	config := &rest.Config{
		Host:        o.Kubernetes.ServerURL,
		BearerToken: o.Kubernetes.ApiToken,
		UserAgent:   kalUserAgent,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: o.Kubernetes.InsecureTLS,
		},
	}
	setDefaultConfigOptions(config, o)

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getKubeConfigClient(o *types.Options) (*kubernetes.Clientset, error) {
	if o.Kubernetes == nil {
		return nil, errors.New("invalid kubernetes options")
	}

	if o.Kubernetes.KubeConfigPath == "" {
		return nil, errors.New("invalid configuration file for kubernetes")
	}

	config, err := clientcmd.BuildConfigFromFlags(o.Kubernetes.ServerURL, o.Kubernetes.KubeConfigPath)
	if err != nil {
		return nil, err
	}

	if o.Kubernetes.ApiToken != "" {
		gologger.Warning().
			Msg("using provided service account token instead kubeconfig configuration")
		config.BearerToken = o.Kubernetes.ApiToken
	} else {
		o.Kubernetes.ApiToken = config.BearerToken
	}
	setDefaultConfigOptions(config, o)

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getInPodClient(o *types.Options) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	setDefaultConfigOptions(config, o)

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func setDefaultConfigOptions(config *rest.Config, o *types.Options) {
	// removing rate limiting...
	if o.Kubernetes.QPS > 0 {
		config.QPS = o.Kubernetes.QPS
	} else {
		config.QPS = 400
	}

	if o.Kubernetes.Burst > 0 {
		config.Burst = o.Kubernetes.Burst
	} else {
		config.Burst = 400
	}

	if o.Kubernetes.NoRateLimit {
		config.RateLimiter = flowcontrol.NewFakeAlwaysRateLimiter()
	}

	config.UserAgent = kalUserAgent

	if o.Kubernetes.InsecureTLS {
		config.TLSClientConfig = rest.TLSClientConfig{
			Insecure: o.Kubernetes.InsecureTLS,
		}
	}
}

// Close stops the execution of the runner
//
// It closes the output chanel and clear any waiting group
func (r *Runner) Close() {
	close(r.outputChan)
	r.outputWg.Done()
}
