package backup

import (
	"context"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/exec"
	"time"
)

type Etcd struct {
	Etcd struct {
		External struct{
			CaFile string `yaml:"caFile"`
			CertFile string `yaml:"certFile"`
			KeyFile string `yaml:"keyFile"`
			Endpoints []string `yaml:"endpoints"`
		}   `yaml:"external"`
	} `yaml:"etcd"`
}

func Backup(outputDir , kubeConfig ,configMapName string)  {
		ctx := context.Background()
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}
		cm, err := clientset.CoreV1().ConfigMaps("kube-system").Get(ctx,configMapName,v1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	   etcd := Etcd{}
	   if yaml.Unmarshal([]byte(cm.Data["ClusterConfiguration"]), &etcd) != nil {
		   fmt.Println("ConfigMaps parser error,please check etcd in kubeadm-config on kube-system")
	   	   os.Exit(1)
	   }
	   if etcd.Etcd.External.CertFile == "" || etcd.Etcd.External.KeyFile == "" || etcd.Etcd.External.CaFile == ""{
		   fmt.Println("cannot found etcd cert key")
		   os.Exit(1)
	   }
	   endpoint := "--endpoints="
	   for index,item := range etcd.Etcd.External.Endpoints {
	   	   if len(etcd.Etcd.External.Endpoints) == index + 1{
			   endpoint += item
		   }else{
			   endpoint += item + ","
		   }
	   }
	   now := time.Now().Format("2006.01.02-15:04:05")
	   outputDir = outputDir + "/" + now + "-etcd-snapshot.db"

	   secret := fmt.Sprintf(" --cert=%s --key=%s --cacert=%s snapshot save %s",etcd.Etcd.External.CertFile,etcd.Etcd.External.KeyFile,etcd.Etcd.External.CaFile,outputDir)
	   cmd := exec.Command("/bin/bash","-c","etcdctl " + endpoint + secret)
	   output, err := cmd.CombinedOutput()
	   if err != nil {
		 fmt.Println(fmt.Sprint(err) + ": " + string(output))
		   os.Exit(1)
	   }
	  fmt.Println(string(output))

}
