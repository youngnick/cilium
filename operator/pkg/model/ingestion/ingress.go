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

	l.Port = 80

	for _, tls := range ing.Spec.TLS {

		l.Port = 443

		l.Hostnames = append(l.Hostnames, tls.Hosts...)

	}

	sourceResource := model.FullyQualifiedResource{
		Name:      ing.Name,
		Namespace: ing.Namespace,
		Group:     "",
		Version:   "v1",
		Kind:      "Ingress",
	}

	l.Sources = append(l.Sources, sourceResource)

	return l

}
