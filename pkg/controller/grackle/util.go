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
	"crypto/sha256"
	"encoding/base64"
	"os"
	"time"

	"github.com/jmckind/grackle-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// AnnotationIngestHash is the annotation key for the ingest hash value.
	AnnotationIngestHash = "k8s.mkz.io/grackle-ingest-hash"

	// DefaultGrackleImageName is the defualt image to use for Grackle.
	DefaultGrackleImageName = "quay.io/jmckind/grackle"

	// DefaultGrackleImageTag is the defualt tag to use for the Grackle image.
	DefaultGrackleImageTag = "latest"

	// DefaultWebReplicas is the number of Web UI pods to create by default.
	DefaultWebReplicas int32 = 2

	// EnvOperatorPodName is the environment variable containig the operator pod name.
	EnvOperatorPodName = "POD_NAME"

	// LabelComponentKey is the key for the component label.
	LabelComponentKey = "component"

	// LabelComponentIngest is the ingest value for the component label.
	LabelComponentIngest = "ingest"

	// LabelComponentWeb is the web value for the component label.
	LabelComponentWeb = "web"
)

// defaultLabels returns the default set of labels for the cluster.
func defaultLabels(cr *v1alpha1.Grackle) map[string]string {
	return map[string]string{
		"app":     "grackle",
		"cluster": cr.Name,
	}
}

// labelsForCluster returns the labels for all cluster resources.
func labelsForCluster(cr *v1alpha1.Grackle) map[string]string {
	labels := defaultLabels(cr)
	for key, val := range cr.ObjectMeta.Labels {
		labels[key] = val
	}
	return labels
}

// hashValue with return a URL encoded SHA256 hash of the given value.
func hashValue(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// newEvent returns a new event for the given Grackle resource.
func newEvent(cr *v1alpha1.Grackle) *corev1.Event {
	t := time.Now()
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cr.Name + "-",
			Namespace:    cr.Namespace,
		},
		InvolvedObject: corev1.ObjectReference{
			APIVersion:      cr.APIVersion,
			Kind:            cr.Kind,
			Name:            cr.Name,
			Namespace:       cr.Namespace,
			UID:             cr.UID,
			ResourceVersion: cr.ResourceVersion,
		},
		Source: corev1.EventSource{
			Component: os.Getenv(EnvOperatorPodName),
		},
		// Each cluster event is unique so it should not be collapsed with other events
		FirstTimestamp: metav1.Time{Time: t},
		LastTimestamp:  metav1.Time{Time: t},
		Count:          int32(1),
	}
}
