package kubernetes

import (
	"context"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strconv"
	//corev1 "k8s.io/api/core/v1"
)

func clientSet() (*kubernetes.Clientset, error) {
	// Get cluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil
}

func GetIngresses() ([]string, error) {
	// Create mapping
	ingressMap := make(map[string]int)
	// Create Ingress list
	var ingressList []string
	// Get clientset
	log.Debug("getting clientset")
	clientset, err := clientSet()
	if err != nil {
		return nil, err
	}

	// Get ingresses
	log.Info("checking ingresses")
	ingresses, err := clientset.ExtensionsV1beta1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Loop through ingresses and check if shield.controller is enabled
	// Return the ingress list after done looping through
	for _, ingress := range ingresses.Items {
		for key, annotation := range ingress.Annotations {
			if key == "aws.shield.controller" && annotation == "enable" {
				for _, v := range ingress.Status.LoadBalancer.Ingress {
					if v.Hostname == "" {
						log.Infof("Ingress: %s in namespace: %s, not ready yet, skipping..", ingress.Name, ingress.Namespace)
						continue
					}
					ingressMap[v.Hostname] = 0
				}
			}
		}
	}

	if len(ingressMap) > 0 {
		for key, _ := range ingressMap {
			ingressList = append(ingressList, key)
		}
	}
	log.Infof("found %s ingresses", strconv.Itoa(len(ingressList)))
	return ingressList, nil
}
