package kal

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/authorization/v1"
)

var KALCmd = &cobra.Command{
	Use:   "iam",
	Short: "List permissions",
	Run:   Run,
}

func init() {
	KALCmd.Flags().BoolP("show-all", "S", false, "show all results")
}

func Run(cmd *cobra.Command, _ []string) {
	authToken, _ := cmd.Flags().GetString("token")
	server, _ := cmd.Flags().GetString("server")
	insecureTLS, _ := cmd.Flags().GetBool("insecure")
	ns, _ := cmd.Flags().GetString("namespace")
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	noRateLimit, _ := cmd.Flags().GetBool("no-rate-limit")
	ua, _ := cmd.Flags().GetString("user-agent")

	authListOptions := &Options{
		AuthToken:   authToken,
		InsecureTLS: insecureTLS,
		Kubeconfig:  kubeconfig,
		Namespace:   ns,
		NoRateLimit: noRateLimit,
		Server:      server,
		UserAgent:   ua,
	}
	log.Debug().Msgf("auth list options: %#v\n", authListOptions)

	mgr, err := New(authListOptions)
	if err != nil {
		log.Error().Err(err).Msg("failed to create authorization list manager")
		return
	}

	processingChannels := ProcessingProperties{
		ResultChan: make(chan *v1.SelfSubjectAccessReview, 100),
		OutputChan: make(chan *v1.SelfSubjectAccessReview, 100),
	}

	var outputWg sync.WaitGroup
	outputWg.Add(1)
	go func(outChan chan *v1.SelfSubjectAccessReview) {
		defer outputWg.Done()

		for accessReview := range outChan {
			if accessReview.Status.Allowed {
				log.Warn().
					Str("request_verb", accessReview.Spec.ResourceAttributes.Verb).
					Str("group_name", accessReview.Spec.ResourceAttributes.Group).
					Str("group_version", accessReview.Spec.ResourceAttributes.Version).
					Str("resource_name", accessReview.Spec.ResourceAttributes.Name).
					Str("namespace", accessReview.Spec.ResourceAttributes.Namespace).
					Msg("can access resource")
			} else {
				showAll, _ := cmd.Flags().GetBool("show-all")
				if showAll {
					log.Info().
						Str("request_verb", accessReview.Spec.ResourceAttributes.Verb).
						Str("group_name", accessReview.Spec.ResourceAttributes.Group).
						Str("group_version", accessReview.Spec.ResourceAttributes.Version).
						Str("resource_name", accessReview.Spec.ResourceAttributes.Name).
						Str("namespace", accessReview.Spec.ResourceAttributes.Namespace).
						Msg("cannot access resource")
				}
			}
		}
	}(processingChannels.OutputChan)

	ctx := context.WithValue(context.TODO(), ProcessingPropertiesCtxType{},
		processingChannels)
	_, _ = mgr.Start(ctx)
	mgr.Finish(ctx)
	outputWg.Wait()
}
