package listener

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.flangaapis.com/htpasswd-kubernetes-operator/apis/htpasswduser/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type BasicAuthCredentials struct {
	Username string
	Password string
}

type HtpasswdListener struct {
	KubeConfig   string
	PasswordFile string
	PidFile      string
	Namespace    string
	client       *v1.BasicAuthClient
	coreClient   *kubernetes.Clientset
	authMap      map[string]*BasicAuthCredentials
	authMapMutex sync.Mutex
}

func (s *HtpasswdListener) Init() error {
	var config *rest.Config
	var err error
	if s.KubeConfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("loading configuration from %v", s.KubeConfig)
		config, err = clientcmd.BuildConfigFromFlags("", s.KubeConfig)
	}
	if err != nil {
		return err
	}

	client, err := v1.NewForConfig(config, s.Namespace)
	if err != nil {
		return err
	}

	coreClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	s.client = client
	s.coreClient = coreClient
	s.authMap = make(map[string]*BasicAuthCredentials)
	return nil
}

func (s *HtpasswdListener) Generate(ctx context.Context) error {
	users, err := s.client.Users(ctx).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	return s.addUsersToConfig(ctx, users)
}

func (s *HtpasswdListener) Run(ctx context.Context) {
	_, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (runtime.Object, error) {
				return s.client.Users(ctx).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return s.client.Users(ctx).Watch(lo)
			},
		},
		&v1.HtpasswdUser{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				casted, ok := obj.(*v1.HtpasswdUser)
				if !ok {
					logrus.Errorf("Failed to cast create obj to htpasswd-user: %v", obj)
					return
				}

				if err := s.onCreate(ctx, casted); err != nil {
					logrus.Errorf("Failed to create htpasswd-user: %v", err)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldCasted, ok := oldObj.(*v1.HtpasswdUser)
				if !ok {
					logrus.Errorf("Failed to cast old Obj to htpasswd-user: %v", oldObj)
					return
				}

				newCasted, ok := newObj.(*v1.HtpasswdUser)
				if !ok {
					logrus.Errorf("Failed to cast new obj to htpasswd-user: %v", newObj)
					return
				}

				if err := s.onUpdate(ctx, oldCasted, newCasted); err != nil {
					log.Printf("Failed to update htpasswd-user: %v\n", err)
				}
			},
			DeleteFunc: func(obj interface{}) {
				casted, ok := obj.(*v1.HtpasswdUser)
				if !ok {
					logrus.Errorf("Failed to cast remove obj to htpasswd-user: %v", obj)
					return
				}

				if err := s.onRemove(ctx, casted); err != nil {
					logrus.Errorf("Failed to remove htpasswd-user: %v", err)
				}
			},
		},
	)

	go controller.Run(ctx.Done())
	<-ctx.Done()
}

func (s *HtpasswdListener) onCreate(ctx context.Context, in *v1.HtpasswdUser) error {
	// check if user does already exist, if it does, check if it should be updated.
	if _, ok := s.authMap[string(in.UID)]; ok {
		return s.updateUserInConfig(ctx, in.UID, in)
	}

	s.authMapMutex.Lock()
	defer s.authMapMutex.Unlock()
	return s.addAndSyncUserToConfig(ctx, in)
}

func (s *HtpasswdListener) onUpdate(ctx context.Context, _, new *v1.HtpasswdUser) error {
	s.authMapMutex.Lock()
	defer s.authMapMutex.Unlock()

	return s.onCreate(ctx, new)
}

func (s *HtpasswdListener) onRemove(ctx context.Context, x *v1.HtpasswdUser) error {
	s.authMapMutex.Lock()
	defer s.authMapMutex.Unlock()

	return s.removeUserInConfig(ctx, x.UID)
}
