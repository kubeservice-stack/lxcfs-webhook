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

package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/kubeservice-stack/lxcfs-webhook/pkg/common"
	"github.com/kubeservice-stack/lxcfs-webhook/pkg/lxcfs"
	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	v1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter = runtime.ObjectDefaulter(runtimeScheme)
)

var (
	ignoredNamespaces = []string{
		metav1.NamespaceSystem,
		metav1.NamespacePublic,
	}
	requiredLabels = []string{
		common.NameLabel,
		common.InstanceLabel,
		common.VersionLabel,
		common.ComponentLabel,
		common.PartOfLabel,
		common.ManagedByLabel,
	}
	addLabels = map[string]string{
		common.NameLabel:      common.NA,
		common.InstanceLabel:  common.NA,
		common.VersionLabel:   common.NA,
		common.ComponentLabel: common.NA,
		common.PartOfLabel:    common.NA,
		common.ManagedByLabel: common.NA,
	}
)

type WebhookServer struct {
	Server *http.Server
}

// Webhook Server parameters
type WhSvrParameters struct {
	Port           int    // webhook server port
	CertFile       string // path to the x509 certificate for https
	KeyFile        string // path to the x509 private key matching `CertFile`
	SidecarCfgFile string // path to sidecar injector configuration file
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1.AddToScheme(runtimeScheme)
	// defaulting with webhooks:
	// https://github.com/kubernetes/kubernetes/issues/57982
	_ = v1.AddToScheme(runtimeScheme)
}

func admissionRequired(ignoredList []string, admissionAnnotationKey string, metadata *metav1.ObjectMeta) bool {
	// skip special kubernetes system namespaces
	for _, namespace := range ignoredList {
		if metadata.Namespace == namespace {
			glog.Infof("Skip validation for %v for it's in special namespace:%v", metadata.Name, metadata.Namespace)
			return false
		}
	}

	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	var required bool
	switch strings.ToLower(annotations[admissionAnnotationKey]) {
	default:
		required = true
	case "n", "no", "false", "off":
		required = false
	}
	return required
}

func mutationRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	required := admissionRequired(ignoredList, common.AdmissionWebhookAnnotationMutateKey, metadata)
	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	status := annotations[common.AdmissionWebhookAnnotationStatusKey]

	if strings.ToLower(status) == "mutated" {
		required = false
	}

	glog.Infof("Mutation policy for %v/%v: required:%v", metadata.Namespace, metadata.Name, required)
	return required
}

func validationRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	required := admissionRequired(ignoredList, common.AdmissionWebhookAnnotationValidateKey, metadata)
	glog.Infof("Validation policy for %v/%v: required:%v", metadata.Namespace, metadata.Name, required)
	return required
}

// Serve method for webhook server
func (whsvr *WebhookServer) Serve(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		glog.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("Content-Type=%s, expect application/json", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	glog.Infof("Received AdmissionReview: %s\n", string(body))

	var admissionResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		glog.Errorf("Can't decode body: %v", err)
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		fmt.Println(r.URL.Path)
		if r.URL.Path == "/mutate" {
			admissionResponse = whsvr.mutatePod(&ar)
		} else if r.URL.Path == "/validate" {
			admissionResponse = whsvr.validatePod(&ar)
		}
	}

	admissionReview := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		glog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	glog.Infof("Ready to write response ...")
	if _, err := w.Write(resp); err != nil {
		glog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

// main mutation process
func (whsvr *WebhookServer) mutatePod(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	var (
		objectMeta                      *metav1.ObjectMeta
		resourceNamespace, resourceName string
	)

	glog.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, resourceName, req.UID, req.Operation, req.UserInfo)

	var pod corev1.Pod

	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		glog.Errorf("Could not unmarshal raw object to pod: %v", err)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}
	resourceName, resourceNamespace, objectMeta = pod.Name, pod.Namespace, &pod.ObjectMeta

	if !mutationRequired(ignoredNamespaces, objectMeta) {
		glog.Infof("Skipping validation for %s/%s due to policy check", resourceNamespace, resourceName)
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	patchBytes, err := createPodPatch(&pod)
	if err != nil {
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	patchType := v1beta1.PatchTypeJSONPatch

	glog.Infof("AdmissionResponse: patch=%v\n", string(patchBytes))
	ret := &v1beta1.AdmissionResponse{
		UID:       req.UID,
		Allowed:   true,
		Patch:     patchBytes,
		PatchType: &patchType,
	}
	glog.Infof("AdmissionResponse: body=%v\n", ret)
	return ret
}

func createPodPatch(pod *corev1.Pod) ([]byte, error) {

	var patches []patchOperation

	var op = patchOperation{
		Path: "/metadata/annotations",
		Value: map[string]string{
			common.AdmissionWebhookAnnotationStatusKey: "mutated",
		},
	}

	if pod.Annotations == nil {
		op.Op = "add"
	} else {
		op.Op = "add"
		if pod.Annotations[common.AdmissionWebhookAnnotationStatusKey] != "" {
			op.Op = "replace"
		}
		op.Path = "/metadata/annotations/" + escapeJSONPointerValue(common.AdmissionWebhookAnnotationStatusKey)
		op.Value = "mutated"
	}

	patches = append(patches, op)

	containers := pod.Spec.Containers

	// Modify the Pod spec to include the LXCFS volumes, then op the original pod.
	for i := range containers {
		if containers[i].VolumeMounts == nil {
			path := fmt.Sprintf("/spec/containers/%d/volumeMounts", i)
			op = patchOperation{
				Op:    "add",
				Path:  path,
				Value: lxcfs.VolumeMountsTemplate,
			}
			patches = append(patches, op)
		} else {
			path := fmt.Sprintf("/spec/containers/%d/volumeMounts/-", i)
			for _, volumeMount := range lxcfs.VolumeMountsTemplate {
				op = patchOperation{
					Op:    "add",
					Path:  path,
					Value: volumeMount,
				}
				patches = append(patches, op)
			}
		}
	}

	if pod.Spec.Volumes == nil {
		op = patchOperation{
			Op:    "add",
			Path:  "/spec/volumes",
			Value: lxcfs.VolumesTemplate,
		}
		patches = append(patches, op)
	} else {
		for _, volume := range lxcfs.VolumesTemplate {
			op = patchOperation{
				Op:    "add",
				Path:  "/spec/volumes/-",
				Value: volume,
			}
			patches = append(patches, op)
		}
	}

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		glog.Warningf("error in json.Marshal %s: %v", pod.Name, err)
		return nil, err
	}
	return patchBytes, nil
}

// validate deployments and services
func (whsvr *WebhookServer) validatePod(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Allowed: true,
	}
}

func escapeJSONPointerValue(in string) string {
	step := strings.Replace(in, "~", "~0", -1)
	return strings.Replace(step, "/", "~1", -1)
}
