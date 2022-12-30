/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

const (
	AdmissionWebhookAnnotationValidateKey = "lxcfs-admission-webhook.kubernetes.io/validate"
	AdmissionWebhookAnnotationMutateKey   = "lxcfs-admission-webhook.kubernetes.io/mutate"
	AdmissionWebhookAnnotationStatusKey   = "lxcfs-admission-webhook.kubernetes.io/status"
	AdmissionWebhookAnnotationPatchKey    = "lxcfs-admission-webhook.kubernetes.io/applied-patch"

	NameLabel      = "app.kubernetes.io/name"
	InstanceLabel  = "app.kubernetes.io/instance"
	VersionLabel   = "app.kubernetes.io/version"
	ComponentLabel = "app.kubernetes.io/component"
	PartOfLabel    = "app.kubernetes.io/part-of"
	ManagedByLabel = "app.kubernetes.io/managed-by"

	NA = "not_available"
)
