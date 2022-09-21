package ingestion

import (
	"testing"

	"github.com/cilium/cilium/operator/pkg/model"
	slim_networkingv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/networking/v1"
	slim_metav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	"github.com/stretchr/testify/assert"
)

var exactPathType = slim_networkingv1.PathTypeExact

var prefixPathType = slim_networkingv1.PathTypePrefix

var testAnnotations = map[string]string{
	"service.beta.kubernetes.io/dummy-load-balancer-backend-protocol":    "http",
	"service.beta.kubernetes.io/dummy-load-balancer-access-log-enabled":  "true",
	"service.alpha.kubernetes.io/dummy-load-balancer-access-log-enabled": "true",
}

// Add the ingress objects in
// https://github.com/kubernetes-sigs/ingress-controller-conformance/tree/master/features
// as test fixtures

var defaultBackend = &slim_networkingv1.Ingress{
	ObjectMeta: slim_metav1.ObjectMeta{
		Name:      "load-balancing",
		Namespace: "random-namespace",
	},
	Spec: slim_networkingv1.IngressSpec{
		IngressClassName: stringp("cilium"),
		DefaultBackend: &slim_networkingv1.IngressBackend{
			Service: &slim_networkingv1.IngressServiceBackend{
				Name: "default-backend",
				Port: slim_networkingv1.ServiceBackendPort{
					Number: 8080,
				},
			},
		},
	},
}

var defaultBackendListeners = []model.HTTPListener{
	{
		Sources: []model.FullyQualifiedResource{
			{
				Name:      "load-balancing",
				Namespace: "random-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
		},
		Port:     80,
		Hostname: "*",
		Routes: []model.HTTPRoute{
			{
				Backends: []model.Backend{
					{
						Name:      "default-backend",
						Namespace: "random-namespace",
						Port: &model.BackendPort{
							Port: 8080,
						},
					},
				},
			},
		},
	},
}

var hostRules = &slim_networkingv1.Ingress{
	ObjectMeta: slim_metav1.ObjectMeta{
		Name:      "host-rules",
		Namespace: "random-namespace",
	},
	Spec: slim_networkingv1.IngressSpec{
		TLS: []slim_networkingv1.IngressTLS{
			{
				Hosts:      []string{"foo.bar.com"},
				SecretName: "conformance-tls",
			},
		},
		Rules: []slim_networkingv1.IngressRule{
			{
				Host: "*.foo.com",
				IngressRuleValue: slim_networkingv1.IngressRuleValue{
					HTTP: &slim_networkingv1.HTTPIngressRuleValue{
						Paths: []slim_networkingv1.HTTPIngressPath{
							{
								Path: "/",
								Backend: slim_networkingv1.IngressBackend{
									Service: &slim_networkingv1.IngressServiceBackend{
										Name: "wildcard-foo-com",
										Port: slim_networkingv1.ServiceBackendPort{
											Number: 8080,
										},
									},
								},
								PathType: &prefixPathType,
							},
						},
					},
				},
			},
			{
				Host: "foo.bar.com",
				IngressRuleValue: slim_networkingv1.IngressRuleValue{
					HTTP: &slim_networkingv1.HTTPIngressRuleValue{
						Paths: []slim_networkingv1.HTTPIngressPath{
							{
								Path: "/",
								Backend: slim_networkingv1.IngressBackend{
									Service: &slim_networkingv1.IngressServiceBackend{
										Name: "foo-bar-com",
										Port: slim_networkingv1.ServiceBackendPort{
											Name: "http",
										},
									},
								},
								PathType: &prefixPathType,
							},
						},
					},
				},
			},
		},
	},
}

var hostRulesListeners = []model.HTTPListener{
	{
		Name: "ing-host-rules-random-namespace-*.foo.com",
		Sources: []model.FullyQualifiedResource{
			{
				Name:      "host-rules",
				Namespace: "random-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
		},
		Port:     80,
		Hostname: "*.foo.com",
		Routes: []model.HTTPRoute{
			{
				PathMatch: model.StringMatch{
					Prefix: "/",
				},
				Backends: []model.Backend{
					{
						Name:      "wildcard-foo-com",
						Namespace: "random-namespace",
						Port: &model.BackendPort{
							Port: 8080,
						},
					},
				},
			},
		},
	},
	{
		Name: "ing-host-rules-random-namespace-foo.bar.com",
		Sources: []model.FullyQualifiedResource{
			{
				Name:      "host-rules",
				Namespace: "random-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
		},
		Port:     80,
		Hostname: "foo.bar.com",
		Routes: []model.HTTPRoute{
			{
				PathMatch: model.StringMatch{
					Prefix: "/",
				},
				Backends: []model.Backend{
					{
						Name:      "foo-bar-com",
						Namespace: "random-namespace",
						Port: &model.BackendPort{
							Name: "http",
						},
					},
				},
			},
		},
	},
	{
		Name: "ing-host-rules-random-namespace-foo.bar.com",
		Sources: []model.FullyQualifiedResource{
			{
				Name:      "host-rules",
				Namespace: "random-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
			// TODO(youngnick): Add some deduplication logic into the Sources
			// field. (Maybe an add method or something?)
			{
				Name:      "host-rules",
				Namespace: "random-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
		},
		Port:     443,
		Hostname: "foo.bar.com",
		TLS: &model.TLSSecret{
			Name:      "conformance-tls",
			Namespace: "random-namespace",
		},
		Routes: []model.HTTPRoute{
			{
				PathMatch: model.StringMatch{
					Prefix: "/",
				},
				Backends: []model.Backend{
					{
						Name:      "foo-bar-com",
						Namespace: "random-namespace",
						Port: &model.BackendPort{
							Name: "http",
						},
					},
				},
			},
		},
	},
}

var complexIngress = &slim_networkingv1.Ingress{
	ObjectMeta: slim_metav1.ObjectMeta{
		Name:        "dummy-ingress",
		Namespace:   "dummy-namespace",
		Annotations: testAnnotations,
		UID:         "d4bd3dc3-2ac5-4ab4-9dca-89c62c60177e",
	},
	Spec: slim_networkingv1.IngressSpec{
		IngressClassName: stringp("cilium"),
		DefaultBackend: &slim_networkingv1.IngressBackend{
			Service: &slim_networkingv1.IngressServiceBackend{
				Name: "default-backend",
				Port: slim_networkingv1.ServiceBackendPort{
					Number: 8080,
				},
			},
		},
		TLS: []slim_networkingv1.IngressTLS{
			{
				Hosts:      []string{"very-secure.server.com"},
				SecretName: "tls-very-secure-server-com",
			},
			{
				Hosts:      []string{"another-very-secure.server.com"},
				SecretName: "tls-another-very-secure-server-com",
			},
		},
		Rules: []slim_networkingv1.IngressRule{
			{
				IngressRuleValue: slim_networkingv1.IngressRuleValue{
					HTTP: &slim_networkingv1.HTTPIngressRuleValue{
						Paths: []slim_networkingv1.HTTPIngressPath{
							{
								Path: "/dummy-path",
								Backend: slim_networkingv1.IngressBackend{
									Service: &slim_networkingv1.IngressServiceBackend{
										Name: "dummy-backend",
										Port: slim_networkingv1.ServiceBackendPort{
											Number: 8080,
										},
									},
								},
								PathType: &exactPathType,
							},
							{
								Path: "/another-dummy-path",
								Backend: slim_networkingv1.IngressBackend{
									Service: &slim_networkingv1.IngressServiceBackend{
										Name: "another-dummy-backend",
										Port: slim_networkingv1.ServiceBackendPort{
											Number: 8081,
										},
									},
								},
								PathType: &prefixPathType,
							},
						},
					},
				},
			},
		},
	},
}

var complexIngressListeners = []model.HTTPListener{
	{
		Sources: []model.FullyQualifiedResource{
			{
				Name:      "dummy-ingress",
				Namespace: "dummy-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
			{
				Name:      "dummy-ingress",
				Namespace: "dummy-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
		},
		Port:     80,
		Hostname: "*",
		Routes: []model.HTTPRoute{
			{
				Backends: []model.Backend{
					{
						Name:      "default-backend",
						Namespace: "dummy-namespace",
						Port: &model.BackendPort{
							Port: 8080,
						},
					},
				},
			},
			{
				PathMatch: model.StringMatch{
					Exact: "/dummy-path",
				},
				Backends: []model.Backend{
					{
						Name:      "dummy-backend",
						Namespace: "dummy-namespace",
						Port: &model.BackendPort{
							Port: 8080,
						},
					},
				},
			},
			{
				PathMatch: model.StringMatch{
					Prefix: "/another-dummy-path",
				},
				Backends: []model.Backend{
					{
						Name:      "another-dummy-backend",
						Namespace: "dummy-namespace",
						Port: &model.BackendPort{
							Port: 8081,
						},
					},
				},
			},
		},
	},
	{
		Sources: []model.FullyQualifiedResource{
			{
				Name:      "dummy-ingress",
				Namespace: "dummy-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
			{
				Name:      "dummy-ingress",
				Namespace: "dummy-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
		},
		Port:     443,
		Hostname: "another-very-secure.server.com",
		TLS: &model.TLSSecret{
			Name:      "tls-another-very-secure-server-com",
			Namespace: "dummy-namespace",
		},
		Routes: []model.HTTPRoute{
			{
				Backends: []model.Backend{
					{
						Name:      "default-backend",
						Namespace: "dummy-namespace",
						Port: &model.BackendPort{
							Port: 8080,
						},
					},
				},
			},
			{
				PathMatch: model.StringMatch{
					Exact: "/dummy-path",
				},
				Backends: []model.Backend{
					{
						Name:      "dummy-backend",
						Namespace: "dummy-namespace",
						Port: &model.BackendPort{
							Port: 8080,
						},
					},
				},
			},
			{
				PathMatch: model.StringMatch{
					Prefix: "/another-dummy-path",
				},
				Backends: []model.Backend{
					{
						Name:      "another-dummy-backend",
						Namespace: "dummy-namespace",
						Port: &model.BackendPort{
							Port: 8081,
						},
					},
				},
			},
		},
	},
	{
		Sources: []model.FullyQualifiedResource{
			{
				Name:      "dummy-ingress",
				Namespace: "dummy-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
			{
				Name:      "dummy-ingress",
				Namespace: "dummy-namespace",
				Version:   "v1",
				Kind:      "Ingress",
			},
		},
		Port:     443,
		Hostname: "very-secure.server.com",
		TLS: &model.TLSSecret{
			Name:      "tls-very-secure-server-com",
			Namespace: "dummy-namespace",
		},
		Routes: []model.HTTPRoute{
			{
				Backends: []model.Backend{
					{
						Name:      "default-backend",
						Namespace: "dummy-namespace",
						Port: &model.BackendPort{
							Port: 8080,
						},
					},
				},
			},
			{
				PathMatch: model.StringMatch{
					Exact: "/dummy-path",
				},
				Backends: []model.Backend{
					{
						Name:      "dummy-backend",
						Namespace: "dummy-namespace",
						Port: &model.BackendPort{
							Port: 8080,
						},
					},
				},
			},
			{
				PathMatch: model.StringMatch{
					Prefix: "/another-dummy-path",
				},
				Backends: []model.Backend{
					{
						Name:      "another-dummy-backend",
						Namespace: "dummy-namespace",
						Port: &model.BackendPort{
							Port: 8081,
						},
					},
				},
			},
		},
	},
}

func stringp(in string) *string {
	return &in
}

type testcase struct {
	ingress slim_networkingv1.Ingress
	want    []model.HTTPListener
}

func TestIngress(t *testing.T) {

	tests := map[string]testcase{
		"defaultBackend": {
			ingress: *defaultBackend,
			want:    defaultBackendListeners,
		},
		"conformance host rules test": {
			ingress: *hostRules,
			want:    hostRulesListeners,
		},
		"cilium test ingress": {
			ingress: *complexIngress,
			want:    complexIngressListeners,
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {
			listeners := Ingress(tc.ingress)
			assert.Equal(t, listeners, tc.want, "Listeners did not match")
		})
	}
}
