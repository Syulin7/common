// Copyright 2023 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	commonutil "github.com/kubeflow/common/pkg/util"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	schedulerpluginsv1alpha1 "sigs.k8s.io/scheduler-plugins/apis/scheduling/v1alpha1"
)

// KoordinatorReconciler defines a gang-scheduling reconciler for koordinator.sh/koordinator
type KoordinatorReconciler struct {
	*SchedulerFrameworkReconciler
}

func BareKoordinatorReconciler(client client.Client, bgReconciler *BaseGangReconciler, enabled bool) *KoordinatorReconciler {
	return &KoordinatorReconciler{
		SchedulerFrameworkReconciler: BareSchedulerFrameworkReconciler(client, bgReconciler, enabled),
	}
}

// GetGangSchedulerName returns the name of Gang Scheduler will be used, which is "koord-scheduler" for SchedulerFrameworkReconciler
func (r *KoordinatorReconciler) GetGangSchedulerName() string {
	return "koord-scheduler"
}

// DecoratePodForGangScheduling decorates the podTemplate before it's used to generate a pod with information for gang-scheduling
func (r *KoordinatorReconciler) DecoratePodForGangScheduling(
	rtype string,
	podTemplate *corev1.PodTemplateSpec,
	job client.Object,
) {
	if podTemplate.Spec.SchedulerName == "" || podTemplate.Spec.SchedulerName == r.GetGangSchedulerName() {
		podTemplate.Spec.SchedulerName = r.GetGangSchedulerName()
	} else {
		warnMsg := "Another scheduler is specified when gang-scheduling is enabled and it will not be overwritten"
		commonutil.LoggerForReplica(job, rtype).Warn(warnMsg)
		r.GetRecorder().Event(job, corev1.EventTypeWarning, "PodTemplateSchedulerNameAlreadySet", warnMsg)
	}

	if podTemplate.Labels == nil {
		podTemplate.Labels = make(map[string]string)
	}
	podTemplate.Labels[schedulerpluginsv1alpha1.PodGroupLabel] = job.GetName()
}
