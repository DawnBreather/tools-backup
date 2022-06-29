package ik8s

import (
	config2 "common/config"
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

var _clientSet *kubernetes.Clientset

func Initialize(){
	var kubeconfig *string

	kubeconfig = flag.String("kubeconfig", config2.K8S.KubeConfigPath, "")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	_clientSet = clientset
}

func GetUnreadyNodes() []v1.Node{
	res, _ := _clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	var unreadyNodes = funk.Filter(res.Items, func(node v1.Node) bool{
		for _, condition := range node.Status.Conditions{
			if fmt.Sprintf(`%s`, string(condition.Type)) == "Ready" {
				if condition.Status == "Unknown"{
					return true
				}
			}
		}
		return false
	})

	return unreadyNodes.([]v1.Node)

}

func GetAllNodes() []v1.Node{
	res, _ := _clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	return res.Items

}


func DeleteNodeByName(nodeName string) (ok bool){

	if nodeName != "" {

		err := _clientSet.CoreV1().Nodes().Delete(context.TODO(), nodeName, metav1.DeleteOptions{})

		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf(`nodes "%s" not found`, nodeName)) {
				ok = true
			} else {
				logrus.Errorf("ERROR removing node {%s} from the cluster: %v", nodeName, err)
				ok = false
			}
		} else {
			ok = true
		}

		if ok {
			logrus.Infof("INFO successfully successfully posted request to the cluster for removal of the node: {%s}", nodeName)
		}
	} else {
		logrus.Warnf("WARNING empty k8s node name: deletion attempt won't be performed")

		ok = true
	}

	return ok
}

func DeleteNodesByNames(nodeNames []string){
	for _, nodeName := range nodeNames{
		go DeleteNodeByName(nodeName)
	}
}