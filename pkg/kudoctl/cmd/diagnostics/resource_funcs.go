package diagnostics

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/kudobuilder/kudo/pkg/apis/kudo/v1beta1"
	"github.com/kudobuilder/kudo/pkg/kudoctl/env"
	"github.com/kudobuilder/kudo/pkg/kudoctl/kudoinit"
	"github.com/kudobuilder/kudo/pkg/kudoctl/util/kudo"
	kudoutil "github.com/kudobuilder/kudo/pkg/util/kudo"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// resourceFuncsConfig - a wrapper for Kube and Kudo clients and common invocation parameters
// for loading Kube and Kudo resources
type resourceFuncsConfig struct {
	c           *kudo.Client
	ns          string
	instanceObj *v1beta1.Instance
	opts        metav1.ListOptions
	logOpts     corev1.PodLogOptions
}

// newInstanceResources is a configuration for instance-related resources
func newInstanceResources(instanceName string, opts *Options, c *kudo.Client, s *env.Settings) (*resourceFuncsConfig, error) {
	instance, err := c.GetInstance(instanceName, s.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance %s/%s: %v", s.Namespace, instanceName, err)
	}
	if instance == nil {
		return nil, fmt.Errorf("instance %s/%s not found", s.Namespace, instanceName)
	}
	logOpts := corev1.PodLogOptions{}
	if opts.LogSince > 0 {
		logOpts.SinceSeconds = &opts.LogSince
	}
	return &resourceFuncsConfig{
		c:           c,
		ns:          s.Namespace,
		instanceObj: instance,
		opts:        metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", kudoutil.OperatorLabel, instance.Labels[kudoutil.OperatorLabel])},
		logOpts:     logOpts,
	}, nil
}

// newKudoResources is a configuration for Kudo controller related resources
// panics if used to load Kudo CRDs (e.g. instance etc.)
func newKudoResources(c *kudo.Client, logSince int64) (*resourceFuncsConfig, error) {
	opts := metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", kudoinit.DefaultKudoLabel)}
	ns, err := c.CoreV1().Namespaces().List(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get kudo system namespace: %v", err)
	}
	if ns == nil || len(ns.Items) == 0 {
		return nil, fmt.Errorf("kudo system namespace not found")
	}
	logOpts := corev1.PodLogOptions{}
	if logSince > 0 {
		logOpts.SinceSeconds = &logSince
	}
	return &resourceFuncsConfig{
		c:       c,
		ns:      ns.Items[0].Name,
		opts:    opts,
		logOpts: corev1.PodLogOptions{},
	}, nil
}

// object implements runtime.object and
// metav1.object interfaces.
// copied from K8 internal type metaRuntimeInterface
type object interface {
	runtime.Object
	metav1.Object
}

type stringGetter func() string

func (r *resourceFuncsConfig) instance() (runtime.Object, error) {
	return r.instanceObj, nil
}

func (r *resourceFuncsConfig) operatorVersion(name stringGetter) func() (runtime.Object, error) {
	return func() (runtime.Object, error) {
		return r.c.GetOperatorVersion(name(), r.ns)
	}
}

func (r *resourceFuncsConfig) operator(name stringGetter) func() (runtime.Object, error) {
	return func() (runtime.Object, error) {
		return r.c.GetOperator(name(), r.ns)
	}
}

func (r *resourceFuncsConfig) deployments() (runtime.Object, error) {
	obj, err := r.c.AppsV1().Deployments(r.ns).List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) pods() (runtime.Object, error) {
	obj, err := r.c.CoreV1().Pods(r.ns).List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) services() (runtime.Object, error) {
	obj, err := r.c.CoreV1().Services(r.ns).List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) replicaSets() (runtime.Object, error) {
	obj, err := r.c.AppsV1().ReplicaSets(r.ns).List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) statefulSets() (runtime.Object, error) {
	obj, err := r.c.AppsV1().StatefulSets(r.ns).List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) serviceAccounts() (runtime.Object, error) {
	obj, err := r.c.CoreV1().ServiceAccounts(r.ns).List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) clusterRoleBindings() (runtime.Object, error) {
	obj, err := r.c.RbacV1().ClusterRoleBindings().List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) roleBindings() (runtime.Object, error) {
	obj, err := r.c.RbacV1().RoleBindings(r.ns).List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) clusterRoles() (runtime.Object, error) {
	obj, err := r.c.RbacV1().ClusterRoles().List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) roles() (runtime.Object, error) {
	obj, err := r.c.RbacV1().Roles(r.ns).List(r.opts)
	return obj, err
}

func (r *resourceFuncsConfig) log(podName string) (io.ReadCloser, error) {
	req := r.c.CoreV1().Pods(r.ns).GetLogs(podName, &r.logOpts)
	// a hack for tests: fake client returns rest.Request{} for GetLogs and Stream panics with null-pointer
	if reflect.DeepEqual(*req, rest.Request{}) {
		return ioutil.NopCloser(&bytes.Buffer{}), nil
	}
	return req.Stream()
}
