// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package ingestion

import (
	"github.com/cilium/cilium/operator/pkg/model"
	slim_networkingv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/networking/v1"
)

// Ingress translates an Ingress resource to a HTTPListener
func Ingress(ing slim_networkingv1.Ingress) model.HTTPListener {

	l := model.HTTPListener{}

	return l

}
