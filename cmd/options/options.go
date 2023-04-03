/*
 * Copyright 2022 Huawei Cloud Computing Technologies Co., Ltd
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package options

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/models"
	"github.com/kappital/kappital/pkg/routers/flowcontroller"
	"github.com/kappital/kappital/pkg/utils/file"
	"github.com/kappital/kappital/pkg/utils/gateway"
	"github.com/kappital/kappital/pkg/utils/version"
)

const (
	managerEnvPrefix = "MANAGER"

	enableHTTPSEnvKey       = "ENABLE_HTTPS"
	httpsCertFilePathEnv    = "HTTPS_CERT_FILE"
	httpsKeyFilePathEnv     = "HTTPS_KEY_FILE"
	httpsTrustCaFilePathEnv = "HTTPS_TRUST_CA_FILE"
	enableMutualHTTPSEnv    = "ENABLE_MUTUAL_HTTPS"
	tlsConfigEnv            = "TLS_CONFIG"

	noClientCert               = "NO_CLIENT_CERT"
	requestClientCert          = "REQUEST_CLIENT_CERT"
	requireAnyClientCert       = "REQUIRE_ANY_CLIENT_CERT"
	verifyClientCertIfGiven    = "VERIFY_CLIENT_CERT_IF_GIVEN"
	requireAndVerifyClientCert = "REQUIRE_AND_VERIFY_CLIENT_CERT"

	httpsPort = 30330

	minPort = 1000
	maxPort = 65535
)

// ServerRunOptions of the kappital service
type ServerRunOptions struct {
	fs *flag.FlagSet
	web.Config

	FlowControllerConfig *flowcontroller.Config
	DBConfig             *models.DatabaseConfig
	DBWatcherConfig      *models.DatabaseWatcherConfig
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters
func NewServerRunOptions(component string) (*ServerRunOptions, error) {
	var prefix string
	switch component {
	case version.ServiceNameManager:
		prefix = managerEnvPrefix
		if err := version.SetupClusterVersion(); err != nil {
			return nil, fmt.Errorf("cannot get the cluster version, error: %v", err)
		}
	default:
		return nil, fmt.Errorf("component is invalid")
	}

	ip, err := gateway.GetLocalIP()
	if err != nil {
		return nil, fmt.Errorf("get local ip failed, error: %v", err)
	}

	s := &ServerRunOptions{
		fs: flag.NewFlagSet(component, flag.ContinueOnError),
		Config: web.Config{
			Listen: web.Listen{HTTPSAddr: ip, HTTPSPort: httpsPort},
		},
		FlowControllerConfig: flowcontroller.DefaultFlowControllerConfig(),
		DBConfig:             models.DefaultDatabaseConfiguration(),
		DBWatcherConfig:      models.DefaultDatabaseWatcherConfig(),
	}
	s.initFlagSet()
	klog.InitFlags(s.fs)

	if err = s.getFlagSetValue(prefix); err != nil {
		return nil, fmt.Errorf("resolve config error: %w ", err)
	}
	if err = s.generateConfig(); err != nil {
		return nil, fmt.Errorf("cannot get the correct config, err: %v", err)
	}
	if err = s.secureServer(); err != nil {
		return nil, fmt.Errorf("cannot start the secure server, err: %v", err)
	}
	return s, nil
}

func (s *ServerRunOptions) initFlagSet() {
	// Server flags of http and https
	s.fs.IntVar(&s.Listen.HTTPSPort, "https-port", s.Listen.HTTPSPort,
		"The port on which to serve HTTPS with authentication and authorization")

	// Database flags
	s.fs.StringVar(&s.DBConfig.SQLDriver, "sql-driver", s.DBConfig.SQLDriver,
		"Which Database driver to use.")
	s.fs.IntVar(&s.DBConfig.MaxIdle, "sql-max-idle", s.DBConfig.MaxIdle,
		"DataBase max idle connection. High number means high memory usage and handles.")
	s.fs.IntVar(&s.DBConfig.MaxConn, "sql-max-conn", s.DBConfig.MaxConn,
		"DataBase max active connection. High number means high memory usage and handles.")
	s.fs.IntVar(&s.DBConfig.MaxLifetime, "db-max-lifetime", s.DBConfig.MaxLifetime,
		"DataBase connection max lifetime (in seconds). default is 1800")
	s.fs.StringVar(&s.DBConfig.SslEnable, "sql-tls-enable", s.DBConfig.SslEnable,
		"enable tls of database connection.")
	s.fs.DurationVar(&s.DBWatcherConfig.ListenerMaxReconnectInterval, "max-database-reconnect-interval",
		s.DBWatcherConfig.ListenerMaxReconnectInterval,
		"max database reconnect interval in seconds for watching table.")
	s.fs.DurationVar(&s.DBWatcherConfig.ListenerMinReconnectInterval, "min-database-reconnect-interval",
		s.DBWatcherConfig.ListenerMinReconnectInterval,
		"min database reconnect interval in seconds for watching table.")
}

func (s *ServerRunOptions) getFlagSetValue(prefix string) error {
	if err := s.fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	s.getFlagsValueFromEnv(prefix)
	return nil
}

func (s *ServerRunOptions) getFlagsValueFromEnv(prefix string) {
	flagsAlreadySet := make(map[string]bool)
	s.fs.Visit(func(f *flag.Flag) {
		flagsAlreadySet[f.Name] = true
	})

	s.fs.VisitAll(func(f *flag.Flag) {
		k := toEnvKey(prefix, f.Name)
		if v := os.Getenv(k); v != "" {
			if flagsAlreadySet[f.Name] {
				fmt.Printf("flag %s has been set explicitly, ignore environment variable %s\n", f.Name, k)
			} else {
				if err := s.fs.Set(f.Name, v); err != nil {
					fmt.Printf("invalid value %v for %s\n", v, k)
				}
				fmt.Printf("recongnized environment variable %s=%s\n", k, v)
			}
		}
	})
}

func (s *ServerRunOptions) generateConfig() error {
	s.Listen.EnableHTTP = false
	s.Listen.EnableHTTPS = parseBool(os.Getenv(enableHTTPSEnvKey))
	s.Listen.HTTPSCertFile = os.Getenv(httpsCertFilePathEnv)
	s.Listen.HTTPSKeyFile = os.Getenv(httpsKeyFilePathEnv)
	s.Listen.TrustCaFile = os.Getenv(httpsTrustCaFilePathEnv)
	if len(s.Listen.TrustCaFile) == 0 {
		s.Listen.EnableMutualHTTPS = false
	} else {
		enable, err := strconv.ParseBool(os.Getenv(enableMutualHTTPSEnv))
		if err != nil {
			klog.Warningf("cannot trans the env [%s] to the bool type, will set the value with true", enableMutualHTTPSEnv)
			enable = true
		}
		s.Listen.EnableMutualHTTPS = enable
	}
	if s.Listen.EnableMutualHTTPS {
		_, s.Listen.ClientAuth = getClientAuth()
	}
	if !s.Listen.EnableHTTPS {
		return fmt.Errorf("https is disable, invalid start method")
	}
	if !(minPort <= s.Listen.HTTPSPort && s.Listen.HTTPSPort <= maxPort) {
		s.Listen.HTTPSPort = httpsPort
	}

	if !file.IsFileExist(s.Listen.HTTPSCertFile) || !file.IsFileExist(s.Listen.HTTPSKeyFile) {
		return fmt.Errorf("https is enable, but cannot find the cert and/or key file")
	}
	if len(s.Listen.TrustCaFile) > 0 && !file.IsFileExist(s.Listen.TrustCaFile) {
		return fmt.Errorf("enable the mutual https, but cannot find the trust ca file")
	}

	return nil
}

func (s *ServerRunOptions) secureServer() error {
	web.BConfig.Listen = s.Listen
	web.BConfig.CopyRequestBody = true
	web.BeeApp.Server.TLSConfig = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		NextProtos:   []string{"h2", "http/1.1"},
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
	}
	if !s.Listen.EnableHTTPS {
		return nil
	}
	certPEMBlock, err := ioutil.ReadFile(s.Listen.HTTPSCertFile)
	if err != nil {
		return err
	}
	keyPEMBlock, err := ioutil.ReadFile(s.Listen.HTTPSKeyFile)
	if err != nil {
		return err
	}
	certificate, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return err
	}
	web.BeeApp.Server.TLSConfig.Certificates = append(web.BeeApp.Server.TLSConfig.Certificates, certificate)
	web.BeeApp.Server.TLSConfig.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return &web.BeeApp.Server.TLSConfig.Certificates[0], nil
	}
	web.BeeApp.Server.TLSConfig.ClientAuth, web.BeeApp.Cfg.Listen.ClientAuth = getClientAuth()
	web.BeeApp.Server.TLSConfig.InsecureSkipVerify = false
	return nil
}

func getClientAuth() (tls.ClientAuthType, int) {
	switch os.Getenv(tlsConfigEnv) {
	case noClientCert:
		return tls.NoClientCert, 0
	case requestClientCert:
		return tls.RequestClientCert, 1
	case requireAnyClientCert:
		return tls.RequireAnyClientCert, 2
	case verifyClientCertIfGiven:
		return tls.VerifyClientCertIfGiven, 3
	case requireAndVerifyClientCert:
		return tls.RequireAndVerifyClientCert, 4
	default:
		return tls.RequireAndVerifyClientCert, 4
	}
}

func toEnvKey(prefix, name string) string {
	return prefix + "_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
}

func parseBool(str string) bool {
	res, err := strconv.ParseBool(str)
	if err != nil {
		klog.Warningf("cannot parse the bool from the string, will use the default result, err: %v", err)
	}
	return res
}
