package operator

import (
	"context"

	"github.com/sirupsen/logrus"
	v1 "htpasswd-operator/pkg/apis/htpasswduser/v1"
	"htpasswd-operator/pkg/client/informers/externalversions"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Listens on the htpasswd users and handles changes to them
func (s *HtpasswdOperator) watchHtpasswdUsers(factory externalversions.SharedInformerFactory, queue workqueue.RateLimitingInterface) (cache.SharedIndexInformer, error) {
	htpasswdInformer := factory.Flanga().V1().HtpasswdUsers().Informer()
	_, err := htpasswdInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.ObjectToName(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.ObjectToName(newObj)
			if err == nil {
				if htpasswdUserChanged(oldObj, newObj) {
					queue.Add(key)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.ObjectToName(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	return htpasswdInformer, err
}

// Listens on secrets and handles changes to them. Only updates the users if the secret is referenced by one of them.
func (s *HtpasswdOperator) watchSecrets(ctx context.Context, coreInformer informers.SharedInformerFactory, _ workqueue.RateLimitingInterface) (cache.SharedIndexInformer, error) {
	informer := coreInformer.Core().V1().Secrets().Informer()
	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.ObjectToName(obj)
			if err == nil {
				s.checkSecretChange(ctx, key, s.fetchAndUpsertHtpasswdUser)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.ObjectToName(newObj)
			if err == nil {
				s.checkSecretChange(ctx, key, s.fetchAndUpsertHtpasswdUser)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.ObjectToName(obj)
			if err == nil {
				s.checkSecretChange(ctx, key, s.removeUserFromPasswordConfig)
			}
		},
	})

	return informer, err
}

// Listens on configMaps and handles changes to them. Only updates the users if the configMap is referenced by one of them.
func (s *HtpasswdOperator) watchConfigMaps(ctx context.Context, coreInformer informers.SharedInformerFactory, _ workqueue.RateLimitingInterface) (cache.SharedIndexInformer, error) {
	informer := coreInformer.Core().V1().ConfigMaps().Informer()
	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.ObjectToName(obj)
			if err == nil {
				s.checkConfigmapChange(ctx, key, s.fetchAndUpsertHtpasswdUser)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.ObjectToName(newObj)
			if err == nil {
				s.checkConfigmapChange(ctx, key, s.fetchAndUpsertHtpasswdUser)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.ObjectToName(obj)
			if err == nil {
				s.checkConfigmapChange(ctx, key, s.removeUserFromPasswordConfig)
			}
		},
	})

	return informer, err
}

func checkRefChange(ctx context.Context, userKey cache.ObjectName, key cache.ObjectName, ref *v1.Reference, onChange func(ctx context.Context, userKey cache.ObjectName) error) {
	if ref.Name == key.Name {
		if err := onChange(ctx, userKey); err != nil {
			logrus.Errorf("Failed at onChange for configmap change: %v", err)
		}
	}
}
