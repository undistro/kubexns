// Copyright 2023 Undistro Authors
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

package options

import (
	"log"
	"os"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	*genericclioptions.ConfigFlags
	genericclioptions.IOStreams

	DefaultMode        os.FileMode
	Dir                string
	IgnoreNotFound     bool
	ConfigMaps         []types.NamespacedName
	Secrets            []types.NamespacedName
	ConfigMapsSelector string
	SecretsSelector    string
}

func NewFromEnv() *Options {
	opts := &Options{
		ConfigFlags: genericclioptions.NewConfigFlags(false),
		IOStreams:   genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	}
	if v, ok := os.LookupEnv("DIR"); ok {
		opts.Dir = v
	} else {
		opts.Dir = "/tmp"
	}
	if v, ok := os.LookupEnv("IGNORE_NOT_FOUND"); ok {
		opts.IgnoreNotFound, _ = strconv.ParseBool(v)
	} else {
		opts.IgnoreNotFound = false
	}
	if v, ok := os.LookupEnv("CONFIGMAPS"); ok {
		opts.ConfigMaps = parseNamespacedNames(v)
	}
	if v, ok := os.LookupEnv("SECRETS"); ok {
		opts.Secrets = parseNamespacedNames(v)
	}
	if v, ok := os.LookupEnv("CONFIGMAPS_SELECTOR"); ok {
		opts.ConfigMapsSelector = v
	}
	if v, ok := os.LookupEnv("SECRETS_SELECTOR"); ok {
		opts.SecretsSelector = v
	}
	if v, ok := os.LookupEnv("DEFAULT_MODE"); ok {
		p, err := strconv.ParseInt(v, 8, 32)
		if err != nil {
			log.Fatalln("failed to parse 'DEFAULT_MODE':", err)
		}
		opts.DefaultMode = os.FileMode(p)
	} else {
		opts.DefaultMode = 0644
	}
	return opts
}

func parseNamespacedNames(s string) []types.NamespacedName {
	var result []types.NamespacedName
	for _, nn := range strings.Split(s, ",") {
		nns := strings.Split(nn, "/")
		if len(nns) > 1 {
			result = append(result, types.NamespacedName{Namespace: strings.TrimSpace(nns[0]), Name: strings.TrimSpace(nns[1])})
		} else if nns[0] != "" {
			result = append(result, types.NamespacedName{Name: strings.TrimSpace(nns[0])})
		}
	}
	return result
}
