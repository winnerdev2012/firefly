// Copyright © 2021 Kaleido, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"os"

	"github.com/spf13/viper"
)

// Key are the known configuration keys
type Key string

const (
	Lang              Key = "lang"
	LogLevel          Key = "log.level"
	LogColor          Key = "log.color"
	DebugPort         Key = "debug.port"
	HttpAddress       Key = "http.address"
	HttpPort          Key = "http.port"
	HttpReadTimeout   Key = "http.readTimeout"
	HttpWriteTimeout  Key = "http.writeTimeout"
	HttpTLSEnabled    Key = "http.tls.enabled"
	HttpTLSClientAuth Key = "http.tls.clientAuth"
	HttpTLSCertsFile  Key = "http.tls.certsFile"
	HttpTLSKeyFile    Key = "http.tls.keyFile"
)

func Reset() {
	viper.Reset()

	// Set defaults
	viper.SetDefault(string(Lang), "en")
	viper.SetDefault(string(LogLevel), "info")
	viper.SetDefault(string(LogColor), true)
	viper.SetDefault(string(DebugPort), -1)
	viper.SetDefault(string(HttpAddress), "127.0.0.1")
	viper.SetDefault(string(HttpPort), 5000)
	viper.SetDefault(string(HttpReadTimeout), 15)
	viper.SetDefault(string(HttpWriteTimeout), 15)
}

// ReadConfig initializes the config
func ReadConfig(cfgFile string) error {
	Reset()

	// Set precedence order for reading config location
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		f, err := os.Open(cfgFile)
		if err == nil {
			defer f.Close()
			err = viper.ReadConfig(f)
		}
		return err
	} else {
		viper.SetConfigName("firefly.core")
		viper.AddConfigPath("/etc/firefly/")
		viper.AddConfigPath("$HOME/.firefly")
		viper.AddConfigPath(".")
		return viper.ReadInConfig()
	}
}

// GetString gets a configuration string
func GetString(key Key) string {
	return viper.GetString(string(key))
}

// GetBool gets a configuration bool
func GetBool(key Key) bool {
	return viper.GetBool(string(key))
}

// GetUInt gets a configuration uint
func GetUint(key Key) uint {
	return viper.GetUint(string(key))
}

// GetInt gets a configuration uint
func GetInt(key Key) int {
	return viper.GetInt(string(key))
}

// Set allows runtime setting of config (used in unit tests)
func Set(key Key, value interface{}) {
	viper.Set(string(key), value)
}
