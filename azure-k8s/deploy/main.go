package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/joho/godotenv"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// 定义全局变量
var (
	subscriptionID string // Azure订阅ID
	resourceGroup  string // Azure资源组名称
	clusterName    string // AKS集群名称
	namespace      string // Kubernetes命名空间
)

// init函数初始化配置
// 从.env文件加载Azure认证信息，并设置默认的集群名称和命名空间
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroup = os.Getenv("AZURE_RESOURCE_GROUP")
	clusterName = "learn-ask-tmp1" // AKS集群名称
	namespace = "default"
	// 这里的明明空间可根据需要进行修改
}

// getK8sClient 获取Kubernetes客户端
// 通过Azure认证信息创建Kubernetes客户端，用于与AKS集群交互
func getK8sClient() (*kubernetes.Clientset, error) {
	// 创建Azure认证凭据
	// 使用服务主体（Service Principal）的客户端ID和密钥进行认证
	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("AZURE_TENANT_ID"),
		os.Getenv("AZURE_CLIENT_ID"),
		os.Getenv("AZURE_CLIENT_SECRET"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// 创建AKS客户端
	// 用于管理AKS集群
	aksClient, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	// 获取集群凭据
	// 获取用于访问Kubernetes集群的kubeconfig信息
	clusterCredential, err := aksClient.ListClusterUserCredentials(context.Background(), resourceGroup, clusterName, nil)
	if err != nil {
		return nil, err
	}

	// 创建临时kubeconfig文件
	// 用于存储集群访问凭据
	kubeconfig, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		return nil, err
	}
	defer os.Remove(kubeconfig.Name())

	// 写入kubeconfig内容
	// 将集群凭据写入临时文件
	if _, err := kubeconfig.Write(clusterCredential.Kubeconfigs[0].Value); err != nil {
		return nil, err
	}

	// 从kubeconfig创建REST配置
	// 用于创建Kubernetes客户端
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig.Name())
	if err != nil {
		return nil, err
	}

	// 创建Kubernetes客户端
	// 返回可用于与集群交互的客户端实例
	return kubernetes.NewForConfig(config)
}

// createDeployment 创建Deployment
// 参数：
// - clientset: Kubernetes客户端
// - name: Deployment名称
// - replicas: 副本数量
// 返回：错误信息
func createDeployment(clientset *kubernetes.Clientset, name string, replicas int32) error {
	// 创建Deployment对象
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			// 定义Pod选择器
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			// 定义Pod模板
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: "nginx:latest", // 这里使用的是nginx，创建后可以直接访问一个基础的web页面
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// 在指定命名空间中创建Deployment
	_, err := clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	return err
}

// createService 创建Service
// 参数：
// - clientset: Kubernetes客户端
// - name: Service名称
// 返回：错误信息
func createService(clientset *kubernetes.Clientset, name string) error {
	// 创建Service对象
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ServiceSpec{
			// 选择器，用于匹配Pod
			Selector: map[string]string{
				"app": name,
			},
			// 定义服务端口
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
			// 使用LoadBalancer类型，创建外部负载均衡器
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}

	// 在指定命名空间中创建Service
	_, err := clientset.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
	return err
}

// deleteDeployment 删除Deployment
// 参数：
// - clientset: Kubernetes客户端
// - name: Deployment名称
// 返回：错误信息
func deleteDeployment(clientset *kubernetes.Clientset, name string) error {
	return clientset.AppsV1().Deployments(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

// deleteService 删除Service
// 参数：
// - clientset: Kubernetes客户端
// - name: Service名称
// 返回：错误信息
func deleteService(clientset *kubernetes.Clientset, name string) error {
	return clientset.CoreV1().Services(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

// listDeployments 列出所有Deployment
// 参数：
// - clientset: Kubernetes客户端
// 返回：错误信息
func listDeployments(clientset *kubernetes.Clientset) error {
	// 获取指定命名空间中的所有Deployment
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// 打印Deployment信息
	fmt.Println("\n列出所有Deployment：")
	for _, deployment := range deployments.Items {
		fmt.Printf("名称: %s\n", deployment.Name)
		fmt.Printf("副本数: %d\n", *deployment.Spec.Replicas)
		fmt.Printf("镜像: %s\n", deployment.Spec.Template.Spec.Containers[0].Image)
		fmt.Println("---")
	}
	return nil
}

// main函数：程序入口
// 支持的命令行参数：
// - action: 操作类型（create/delete/list）
// - name: 应用名称
// - replicas: 副本数量
func main() {
	// 定义命令行参数
	action := flag.String("action", "list", "操作类型：create/delete/list")
	name := flag.String("name", "demo-app", "应用名称")
	replicas := flag.Int("replicas", 2, "副本数量")
	flag.Parse()

	// 获取Kubernetes客户端
	clientset, err := getK8sClient()
	if err != nil {
		log.Fatal(err)
	}

	// 根据action参数执行相应操作
	switch *action {
	case "create":
		fmt.Printf("开始创建应用 %s...\n", *name)
		if err := createDeployment(clientset, *name, int32(*replicas)); err != nil {
			log.Fatal(err)
		}
		if err := createService(clientset, *name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("应用创建成功！\n")
	case "delete":
		fmt.Printf("开始删除应用 %s...\n", *name)
		if err := deleteDeployment(clientset, *name); err != nil {
			log.Fatal(err)
		}
		if err := deleteService(clientset, *name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("应用删除成功！\n")
	case "list":
		if err := listDeployments(clientset); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("无效的操作类型。请使用 -action 参数指定操作类型：create/delete/list")
		flag.PrintDefaults()
	}
}
