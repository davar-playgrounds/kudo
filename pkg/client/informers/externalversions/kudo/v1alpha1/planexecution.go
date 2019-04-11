/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	kudov1alpha1 "github.com/kudobuilder/kudo/pkg/apis/kudo/v1alpha1"
	versioned "github.com/kudobuilder/kudo/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kudobuilder/kudo/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kudobuilder/kudo/pkg/client/listers/kudo/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// PlanExecutionInformer provides access to a shared informer and lister for
// PlanExecutions.
type PlanExecutionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.PlanExecutionLister
}

type planExecutionInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewPlanExecutionInformer constructs a new informer for PlanExecution type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPlanExecutionInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPlanExecutionInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredPlanExecutionInformer constructs a new informer for PlanExecution type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPlanExecutionInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KudoV1alpha1().PlanExecutions(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KudoV1alpha1().PlanExecutions(namespace).Watch(options)
			},
		},
		&kudov1alpha1.PlanExecution{},
		resyncPeriod,
		indexers,
	)
}

func (f *planExecutionInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPlanExecutionInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *planExecutionInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&kudov1alpha1.PlanExecution{}, f.defaultInformer)
}

func (f *planExecutionInformer) Lister() v1alpha1.PlanExecutionLister {
	return v1alpha1.NewPlanExecutionLister(f.Informer().GetIndexer())
}