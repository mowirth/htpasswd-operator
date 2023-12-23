package operator

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mowirth/htpasswd-operator/pkg/client/clientset/versioned"
	htpasswdv1 "github.com/mowirth/htpasswd-operator/pkg/client/clientset/versioned/typed/htpasswduser/v1"
	"github.com/mowirth/htpasswd-operator/pkg/client/informers/externalversions"
	"github.com/sirupsen/logrus"
	v13 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

type HtpasswdOperator struct {
	KubeConfig   string
	PasswordFile string
	PidFile      string

	coreClient         *kubernetes.Clientset
	htpasswdUserClient htpasswdv1.HtpasswdUsersGetter
	htpasswdInformer   cache.SharedIndexInformer
	secretInformer     cache.SharedIndexInformer
	configMapInformer  cache.SharedIndexInformer
	userRefMap         map[cache.ObjectName]*CredentialReference
	userRefMapMutex    sync.RWMutex
	passwordMap        map[cache.ObjectName]*BasicAuthCredentials
	passwordMapMutex   sync.RWMutex
	recorder           record.EventRecorder
	queue              workqueue.RateLimitingInterface
}

func (s *HtpasswdOperator) Init(ctx context.Context) error {
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

	coreClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	s.coreClient = coreClient

	htpasswdClient, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}

	coreInformer := informers.NewSharedInformerFactoryWithOptions(coreClient, 0)
	htpasswdInformer := externalversions.NewSharedInformerFactoryWithOptions(htpasswdClient, 0)

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	runtime.Must(v13.AddToScheme(scheme.Scheme))
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logrus.Tracef)
	eventBroadcaster.StartRecordingToSink(&corev1.EventSinkImpl{Interface: s.coreClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v13.EventSource{Component: "htpasswdusers"})

	informer, err := s.watchHtpasswdUsers(htpasswdInformer, queue)
	if err != nil {
		return err
	}

	secretInformer, err := s.watchSecrets(ctx, coreInformer, queue)
	if err != nil {
		return err
	}

	configMapInformer, err := s.watchConfigMaps(ctx, coreInformer, queue)
	if err != nil {
		return err
	}

	s.htpasswdUserClient = htpasswdClient.FlangaV1()
	s.htpasswdInformer = informer
	s.secretInformer = secretInformer
	s.configMapInformer = configMapInformer
	s.recorder = recorder
	s.queue = queue
	s.userRefMap = make(map[cache.ObjectName]*CredentialReference)
	s.passwordMap = make(map[cache.ObjectName]*BasicAuthCredentials)

	return nil
}

func (s *HtpasswdOperator) Generate(ctx context.Context) error {
	users, err := s.htpasswdUserClient.HtpasswdUsers("").List(ctx, metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("Failed to fetch htpasswd user: %v", err)
		return err
	}

	for _, user := range users.Items {
		if err := s.upsertUser(ctx, &user); err != nil {
			logrus.Errorf("Failed to upsert user: %v", err)
			s.recorder.Eventf(&user, v13.EventTypeWarning, "UpdateFailed", err.Error())
			return err
		}
	}

	return nil
}

func (s *HtpasswdOperator) Run(ctx context.Context) {
	defer runtime.HandleCrash()
	defer s.queue.ShutDown()

	go s.htpasswdInformer.Run(ctx.Done())
	go s.secretInformer.Run(ctx.Done())
	go s.configMapInformer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done()) {
		runtime.HandleError(fmt.Errorf("timeout on waiting for cache sync"))
		return
	}

	wait.Until(func() {
		s.processQueue(ctx, s.queue)
	}, time.Second, ctx.Done())

	logrus.Errorf("Shutting down controller")
}

func (s *HtpasswdOperator) processQueue(ctx context.Context, queue workqueue.RateLimitingInterface) bool {
	key, exit := queue.Get()
	if exit {
		return false
	}

	defer queue.Done(key)
	if err := s.processUser(ctx, key.(cache.ObjectName)); err != nil {
		if s.queue.NumRequeues(key) < 5 {
			logrus.Errorf("Failed to update key, retrying later: %v", err)
			s.queue.AddRateLimited(key)
		} else {
			logrus.Errorf("Failed to update key, not retrying: %v", err)
			s.queue.Forget(key)
		}

		return true
	}

	s.queue.Forget(key)
	return true
}

func (s *HtpasswdOperator) processUser(ctx context.Context, key cache.ObjectName) (updateErr error) {
	object, exists, err := s.htpasswdInformer.GetIndexer().GetByKey(key.String())
	if err != nil {
		return err
	}

	if !exists {
		logrus.Infof("HtpasswdUser %v is gone, removing", key)
		return s.removeUserInConfig(ctx, key)
	}

	htpasswdUser, err := convertHtpasswdUser(object)
	if err != nil {
		return err
	}

	defer func() {
		s.updateHtpasswdUserStatus(ctx, htpasswdUser, updateErr)
	}()
	if err := s.upsertUser(ctx, htpasswdUser); err != nil {
		s.recorder.Eventf(htpasswdUser, v13.EventTypeWarning, "UpdateFailed", err.Error())
		return err
	}

	s.recorder.Event(htpasswdUser, v13.EventTypeNormal, "Success", "HtpasswdUser was configured")
	return nil
}
