package kubeClient

import (
	"context"
	"log"
	"os"

	"github.com/infracloudio/krius/pkg/random"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

type KubeConfig struct {
	Namespace string
	Context   string
}

func BuildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func (k KubeConfig) InitClient() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	kubeconfig := dirname + "/.kube/config"
	config, err := BuildConfigFromFlags(k.Context, kubeconfig)
	if err != nil {
		return err
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func (k KubeConfig) CreateSecret(secretSpec map[string][]byte, secretName string) error {
	secretsClient := clientset.CoreV1().Secrets(k.Namespace)
	if k.HasSecret(secretName) {
		secretName = secretName + random.RandStringRunes(4)
	}
	// Create secret
	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: k.Namespace,
		},
		Data: secretSpec,
	}
	_, err := secretsClient.Create(context.Background(), secret, metav1.CreateOptions{
		FieldManager: "objStore",
	})
	if err != nil {
		return err
	}
	return nil
}

func (k KubeConfig) HasSecret(name string) bool {
	secretsClient := clientset.CoreV1().Secrets(k.Namespace)
	cm, err := secretsClient.Get(context.TODO(), name, metav1.GetOptions{})
	return cm != nil && err == nil
}

func (k KubeConfig) CreateNSIfNotExist() error {
	err := k.CheckNamespaceExist()
	if err != nil {
		nsName := &apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: k.Namespace,
			},
		}
		_, err := clientset.CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (k KubeConfig) CheckNamespaceExist() error {
	_, err := clientset.CoreV1().Namespaces().Get(context.Background(), k.Namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (k KubeConfig) GetServiceInfo(svcName string) []string {
	list, err := clientset.CoreV1().Services(k.Namespace).Get(context.Background(), svcName, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	targets := []string{}
	for _, v := range list.Status.LoadBalancer.Ingress {
		targets = append(targets, v.Hostname)
	}
	return targets

}
