// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	configv1 "github.com/openshift/api/config/v1"
	versioned "github.com/openshift/client-go/config/clientset/versioned"
	internalinterfaces "github.com/openshift/client-go/config/informers/externalversions/internalinterfaces"
	v1 "github.com/openshift/client-go/config/listers/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// BuildInformer provides access to a shared informer and lister for
// Builds.
type BuildInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.BuildLister
}

type buildInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewBuildInformer constructs a new informer for Build type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBuildInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBuildInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredBuildInformer constructs a new informer for Build type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBuildInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().Builds().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().Builds().Watch(context.TODO(), options)
			},
		},
		&configv1.Build{},
		resyncPeriod,
		indexers,
	)
}

func (f *buildInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBuildInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *buildInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configv1.Build{}, f.defaultInformer)
}

func (f *buildInformer) Lister() v1.BuildLister {
	return v1.NewBuildLister(f.Informer().GetIndexer())
}
