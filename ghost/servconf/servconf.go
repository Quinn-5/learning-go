package servconf

import (
	"bytes"
	"errors"
	"flag"
	"path/filepath"
	"regexp"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeconfig *string

// Server config with private values
// for use in backend functions.
// Initialize with servconf.New()
type ServerConfig struct {

	// Name of user requesting server
	username string

	// Name of new server
	serverName string

	// Type of server requested
	serverType string

	// Number of CPU cores to assign
	cpu resource.Quantity

	// Number of GiB RAM to reserve
	ram resource.Quantity

	// Number of MiB disk space to reserve
	disk resource.Quantity

	// IP address to connect to
	ip string

	// Internal port to be exposed
	internalPort int32

	// External port to connect
	externalPort int32

	// Protocol used for communication
	protocol apiv1.Protocol

	// kubeconfig
	clientset *kubernetes.Clientset
}

// Server config with public values
// for use in web handlers.
// Generate with servconf.WebConfig()
type WebConfig struct {
	Username     string
	ServerName   string
	ServerType   string
	CPU          string
	RAM          string
	Disk         string
	IP           string
	InternalPort int32
	ExternalPort int32
}

func New(username string, serverName string) *ServerConfig {
	cfg := &ServerConfig{}

	if kubeconfig == nil {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	cfg.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	cfg.setUsername(username)
	cfg.setServerName(serverName)

	return cfg
}

func (cfg *ServerConfig) GetUsername() string {
	return cfg.username
}

func (cfg *ServerConfig) setUsername(username string) error {
	exp := regexp.MustCompile(`[a-z]([-a-z0-9]*[a-z0-9])?`)
	if !bytes.Equal(exp.Find([]byte(username)), []byte(username)) {
		return errors.New("username must contain only alphanumeric, lowercase characters")
	} else {
		cfg.username = strings.ToLower(username)
		return nil
	}
}

func (cfg *ServerConfig) GetServerName() string {
	return cfg.serverName
}

func (cfg *ServerConfig) setServerName(serverName string) error {
	exp := regexp.MustCompile(`[a-z]([-a-z0-9]*[a-z0-9])?`)
	if !bytes.Equal(exp.Find([]byte(serverName)), []byte(serverName)) {
		return errors.New("servername must contain only alphanumeric, lowercase characters")
	} else {
		cfg.serverName = strings.ToLower(serverName)
		return nil
	}
}

func (cfg *ServerConfig) GetServerType() string {
	return cfg.serverType
}

func (cfg *ServerConfig) SetType(serverType string) error {
	exp := regexp.MustCompile(`[a-z]([-a-z0-9]*[a-z0-9])?`)
	if !bytes.Equal(exp.Find([]byte(serverType)), []byte(serverType)) {
		return errors.New("servertype must contain only alphanumeric, lowercase characters")
	} else {
		cfg.serverType = strings.ToLower(serverType)
		return nil
	}
}

func (cfg *ServerConfig) GetInternalPort() int32 {
	return cfg.internalPort
}

func (cfg *ServerConfig) SetInternalPort(port int32) {
	cfg.internalPort = port
}

func (cfg *ServerConfig) GetExternalPort() int32 {
	return cfg.externalPort
}

func (cfg *ServerConfig) SetExternalPort(port int32) {
	cfg.externalPort = port
}

func (cfg *ServerConfig) GetIP() string {
	return cfg.ip
}

func (cfg *ServerConfig) SetIP(ip string) {
	cfg.ip = ip
}

func (cfg *ServerConfig) GetProtocol() apiv1.Protocol {
	return cfg.protocol
}

func (cfg *ServerConfig) SetProtocol(protocol apiv1.Protocol) {
	cfg.protocol = protocol
}

func (cfg *ServerConfig) GetCPU() resource.Quantity {
	return cfg.cpu
}

func (cfg *ServerConfig) SetCPU(cpu string) {
	if n, err := resource.ParseQuantity(cpu); err == nil {
		cfg.cpu = n
	}
}

func (cfg *ServerConfig) GetRAM() resource.Quantity {
	return cfg.ram
}

func (cfg *ServerConfig) SetRAM(ram string) {
	if n, err := resource.ParseQuantity(ram + "Gi"); err == nil {
		cfg.ram = n
	}
}

func (cfg *ServerConfig) GetDisk() resource.Quantity {
	return cfg.disk
}

func (cfg *ServerConfig) SetDisk(disk string) {
	if n, err := resource.ParseQuantity(disk + "Gi"); err == nil {
		cfg.disk = n
	}
}

func (cfg *ServerConfig) GetKubeConfig() *kubernetes.Clientset {
	return cfg.clientset
}

func (cfg *ServerConfig) WebConfig() *WebConfig {
	webconf := &WebConfig{
		Username:     cfg.GetUsername(),
		ServerName:   cfg.GetServerName(),
		ServerType:   cfg.GetServerType(),
		CPU:          cfg.GetCPU().OpenAPISchemaFormat(),
		RAM:          cfg.GetRAM().OpenAPISchemaFormat(),
		Disk:         cfg.GetRAM().OpenAPISchemaFormat(),
		IP:           cfg.GetIP(),
		InternalPort: cfg.GetInternalPort(),
		ExternalPort: cfg.GetExternalPort(),
	}
	return webconf
}
