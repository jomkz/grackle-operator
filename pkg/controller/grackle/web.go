// Copyright 2019 Grackle Operator authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grackle

import (
	"fmt"

	k8sv1alpha1 "github.com/jmckind/grackle-operator/pkg/apis/k8s/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// newGrackleWebService returns a Deployment resource for the Web UI.
func newGrackleWebDeployment(cr *k8sv1alpha1.Grackle) *appsv1.Deployment {
	labels := labelsForCluster(cr)
	labels[LabelComponentKey] = LabelComponentWeb

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", cr.Name, LabelComponentWeb),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: cr.Spec.Web.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   fmt.Sprintf("%s:%s", DefaultGrackleImageName, cr.Spec.Web.Version),
						Name:    LabelComponentWeb,
						Command: []string{"/grackle-web"},
						Env: []corev1.EnvVar{{
							Name:  "GRK_RETHINKDB_HOST",
							Value: cr.Spec.Datastore.Host,
						}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8000,
							Name:          "http",
						}},
					}},
				},
			},
		},
	}
}

// newGrackleWebService returns a Service resource for the Web UI.
func newGrackleWebService(cr *k8sv1alpha1.Grackle) *corev1.Service {
	labels := labelsForCluster(cr)
	labels[LabelComponentKey] = LabelComponentWeb

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Name:       LabelComponentWeb,
				Port:       80,
				TargetPort: intstr.FromInt(8000),
			}},
		},
	}
}
