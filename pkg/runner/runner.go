package runner

import (
	"context"
	"strings"
	"sync"

	myK8s "github.com/ing-bank/kal/pkg/kubernetes"
	"github.com/ing-bank/kal/pkg/types"
	"github.com/projectdiscovery/gologger"
	"golang.org/x/sync/semaphore"
	v1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Exec will execute the procedure to list all permissions from a given configuration
//
// return : resourcePermissions map[resourceName+version+subresource][]ApiVerb
func (r *Runner) Exec() (resourcePermissions map[string][]string) {
	resourcePermissions = make(map[string][]string, 0)

	gologger.Info().Msgf("running from namespace = %s\n", r.Namespace)

	apiGroupsList := make([]string, 0)

	groupList, _, _, err := r.KubernetesClient.GroupsAndMaybeResources()
	if err != nil {
		gologger.Error().Msgf("could not list api resources")
		return
	}

	for _, group := range groupList.Groups {
		for _, version := range group.Versions {
			apiGroupsList = append(apiGroupsList, version.GroupVersion)
		}
	}

	resources := make([]*Resource, 0)
	for _, group := range apiGroupsList {
		groupResources, err := r.KubernetesClient.DiscoveryClient.ServerResourcesForGroupVersion(group)

		if err != nil {
			gologger.Error().Msgf("could not get resources from group. error: %s\n", err)
			continue
		}

		for _, resource := range groupResources.APIResources {
			y := strings.Split(groupResources.GroupVersion, "/")

			groupName := y[0]
			var groupVersion string
			if len(y) > 1 {
				groupVersion = y[1]
			}

			resourceName := resource.Name
			subResource := ""
			if strings.Contains(resourceName, "/") {
				// example: helmchartrepositories/status -> status is the sub-resource of helmchartrepositories

				x := strings.Split(resourceName, "/")
				resourceName = x[0]
				subResource = x[1]
			}

			resourceItem := &Resource{
				GroupName:    groupName,
				GroupVersion: groupVersion,
				Name:         resourceName,
				Namespaced:   resource.Namespaced,
				SubResource:  subResource,
			}
			if groupName == "v1" {
				resourceItem.GroupVersion = groupName
				resourceItem.GroupName = ""
			}

			resources = append(
				resources,
				resourceItem,
			)

		}
	}

	// output processor start

	r.outputWg.Add(1)
	go func(output chan *Result, rp *map[string][]string) {
		defer r.outputWg.Done()

		for outputResult := range output {
			if len(outputResult.AllowedVerbs) > 0 {
				if _, ok := (*rp)[outputResult.Resource.Name]; !ok {
					(*rp)[outputResult.Resource.Name] = make([]string, 0)
				}

				resultKey := outputResult.Resource.Name

				if outputResult.Resource.GroupVersion != "" {
					resultKey += "/" + outputResult.Resource.GroupVersion
				}

				if outputResult.Resource.SubResource != "" {
					resultKey += "/" + outputResult.Resource.SubResource
				}

				(*rp)[resultKey] = append((*rp)[outputResult.Resource.Name], outputResult.AllowedVerbs...)
			}

			if !r.ShowAll && len(outputResult.AllowedVerbs) == 0 {
				continue
			}

			if len(outputResult.str) > 0 {
				gologger.Silent().Msgf("%s\n", outputResult.str)
			}
		}
	}(r.outputChan, &resourcePermissions)

	var analysisWg sync.WaitGroup
	sem := semaphore.NewWeighted(1)
	semCtx := context.TODO()

	gologger.Info().Msgf("found %d resources and sub-resources\n", len(resources))
	for _, resource := range resources {
		_ = sem.Acquire(semCtx, 1)
		analysisWg.Add(1)
		go func() {
			defer sem.Release(1)
			defer analysisWg.Done()
			r.outputChan <- r.analysis(resource)
		}()
	}

	analysisWg.Wait()

	close(r.outputChan)
	r.outputWg.Wait()
	return
}

func (r *Runner) analysis(resource *Resource) (result *Result) {
	result = &Result{
		Resource:                       resource,
		Namespace:                      r.Namespace,
		SelfSubjectAccessReviewResults: make([]*v1.SelfSubjectAccessReview, 0),
		AllowedVerbs:                   make([]string, 0),
	}

	var verbWg sync.WaitGroup
	var verbChanWg sync.WaitGroup
	verbChan := make(chan string)
	accessReviewChan := make(chan *v1.SelfSubjectAccessReview)

	verbChanWg.Add(1)
	go func() {
		for allowedVerb := range verbChan {
			result.AllowedVerbs = append(result.AllowedVerbs, allowedVerb)
			result.SelfSubjectAccessReviewResults = append(result.SelfSubjectAccessReviewResults, <-accessReviewChan)
		}
		verbChanWg.Done()
	}()

	for _, verb := range myK8s.ApiVerbs {
		gologger.Debug().Msgf("testing resource [%s] -> VERB[%s] NS[%s]\n", resource.String(), verb, r.Namespace)

		ns := r.Namespace
		if !resource.Namespaced {
			ns = ""
		}

		verbWg.Add(1)
		go func(vb, nspace string, resource *Resource) {
			defer verbWg.Done()

			verbAccessReview := r.requestAccessReview(vb, nspace, resource)
			if verbAccessReview.Status.Allowed {
				verbChan <- verb
				accessReviewChan <- verbAccessReview
			}
		}(verb, ns, resource)

	}
	verbWg.Wait()
	close(verbChan)
	close(accessReviewChan)
	verbChanWg.Wait()

	builder := &strings.Builder{}

	builder.WriteString(resource.String())

	if len(result.AllowedVerbs) > 0 {
		builder.WriteString(" [")
		builder.WriteString(types.AU.Green(strings.Join(result.AllowedVerbs, ",")).String())
		builder.WriteRune(']')
	} else {
		builder.WriteString(" [")
		builder.WriteString(types.AU.Red("NO_ALLOWED_VERBS").String())
		builder.WriteRune(']')
	}

	ns := ""
	if result.Resource.Namespaced {
		ns = result.Namespace
	} else {
		ns = "CLUSTER_WIDE"
	}
	builder.WriteString(" [")
	builder.WriteString(types.AU.Blue(ns).String())
	builder.WriteRune(']')

	if r.ShowReason {
		builder.WriteString(" [")
		reasons := make([]string, 0)
		for _, review := range result.SelfSubjectAccessReviewResults {
			reasons = append(reasons, review.Status.Reason)
		}
		builder.WriteString(types.AU.Magenta(strings.Join(reasons, ";")).String())
		builder.WriteRune(']')
	}

	result.str = builder.String()

	return result
}

func (r *Runner) requestAccessReview(verb, ns string, resource *Resource) *v1.SelfSubjectAccessReview {
	sar := &v1.SelfSubjectAccessReview{
		Spec: v1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &v1.ResourceAttributes{
				Verb:        verb,
				Resource:    resource.Name,
				Group:       resource.GroupName,
				Subresource: resource.SubResource,
				Namespace:   ns,
				Name:        resource.Name,
			},
		},
	}

	if resource.Namespaced {
		sar.Spec.ResourceAttributes.Namespace = ns
	}

	accessReviewResponse, err := r.KubernetesClient.
		AuthorizationV1().
		SelfSubjectAccessReviews().
		Create(
			r.Context, // All the requests use the same context for rate limit control
			sar,
			metav1.CreateOptions{},
		)

	if err != nil {
		gologger.Error().Msgf("could not analyze resource [%s] -> [%s] (%s)\n", verb, resource.Name, ns)
		return accessReviewResponse
	}

	return accessReviewResponse
}
