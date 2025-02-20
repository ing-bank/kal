package runner

import (
	"context"
	"strings"
	"sync"

	v1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/kubernetes"
)

const kalUserAgent string = "KAL"

// Runner is the structure holding information about KAL's Runner
type Runner struct {
	KubernetesClient *kubernetes.Clientset
	Namespace        string
	Context          context.Context

	WideOutput bool
	JSONOutput bool
	ShowReason bool
	ShowAll    bool

	outputWg   sync.WaitGroup
	outputChan chan *Result
}

// Resource is the abstraction of a Kubernetes resource
type Resource struct {
	// AllowedVerbs []string
	GroupName    string
	GroupVersion string
	Name         string
	Namespaced   bool
	SubResource  string
}

// String return the string representation of a Resource
func (r *Resource) String() string {
	sb := &strings.Builder{}

	sb.WriteString(r.Name)

	if r.GroupName != "" {
		sb.WriteString("." + r.GroupName)
	}

	sb.WriteString("/" + r.GroupVersion)

	if r.SubResource != "" {
		sb.WriteString("/" + r.SubResource)
	}

	return sb.String()
}

// Result is the structure used to hold information used to present
// the analysis result for a user
type Result struct {
	Namespace                      string
	Resource                       *Resource
	SelfSubjectAccessReviewResults []*v1.SelfSubjectAccessReview
	str                            string
	AllowedVerbs                   []string
}

// AnalysisResult is the structure that contains the analysis information of a Resource
type AnalysisResult struct {
	Allowed bool
	Verb    string
}
