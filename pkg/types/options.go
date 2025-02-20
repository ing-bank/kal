package types

import (
	"github.com/logrusorgru/aurora/v4"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
	"k8s.io/client-go/rest"
)

// AU is the package to control color in the output
var AU *aurora.Aurora

// Options is the structure used to hold the execution options
type Options struct {
	Kubernetes *KubernetesOptions

	Output *OutputOptions

	Verbose bool
	Silent  bool
	NoLogs  bool
}

// Validate checks the provided options for constraints
func (o *Options) Validate() {
	if o.Verbose && o.Silent {
		gologger.Fatal().Msg("verbose and silent output selected")
	}

	o.Kubernetes.Validate()
	o.Output.Validate()
}

// Configure configures the support packages for KAL, based on options
func (o *Options) Configure() {
	InitAurora(o)

	if o.Verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	}

	if o.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	}

	if o.NoLogs {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelFatal)
	}
}

// Kubernetes Options is the structure for Kubernetes Options
type KubernetesOptions struct {
	ApiToken          string
	Burst             int
	InsecureTLS       bool
	KubeConfigPath    string
	Namespace         string
	NoRateLimit       bool
	QPS               float32
	ServerURL         string
	UserToImpersonate string
}

// Validate validates the Kubernetes Options for constraints
func (ko *KubernetesOptions) Validate() {
	if ko.ServerURL == "" {
		gologger.Fatal().Msg("invalid kubernetes server url")
	}

	if ko.NoRateLimit {
		ko.QPS = 400
		ko.Burst = 400
	} else {
		ko.QPS = rest.DefaultQPS
		ko.Burst = rest.DefaultBurst
	}
}

// OutputOptions is the structure for options related to output information
type OutputOptions struct {
	JSON       bool
	NoColor    bool
	ShowAll    bool
	ShowReason bool
}

// Validate validates the provided Output options
func (oo *OutputOptions) Validate() {
	gologger.DefaultLogger.SetFormatter(formatter.NewCLI(oo.NoColor))

	if oo.JSON {
		gologger.DefaultLogger.SetFormatter(&formatter.JSON{})
	}
}

// InitAurora initialize Aurora for colored logging
func InitAurora(o *Options) {
	if AU != nil {
		return
	}

	if o.Output.JSON {
		AU = aurora.New(aurora.WithColors(false))
	} else {
		AU = aurora.New(aurora.WithColors(!o.Output.NoColor))
	}
}
