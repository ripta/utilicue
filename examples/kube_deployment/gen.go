//go:generate go run ../../cmd/cue2go .
package kube_deployment

import (
	_ "k8s.io/api/apps/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
)
