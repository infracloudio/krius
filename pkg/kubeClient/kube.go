package kubeClient

import (
	"context"
	"log"
	"os"

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

	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: k.Namespace,
		},
		Data: secretSpec,
	}
	// Create secret
	_, err := secretsClient.Create(context.TODO(), secret, metav1.CreateOptions{
		FieldManager: "objStore",
	})
	if err != nil {
		return err
	}
	return nil
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
