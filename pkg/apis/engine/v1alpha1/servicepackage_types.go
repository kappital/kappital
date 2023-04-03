/*
 * Copyright 2022 Huawei Cloud Computing Technologies Co., Ltd
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// PendingPhase pending status of the ServicePackage (SP), it means the SP is installing
	PendingPhase = "Pending"
	// RunningPhase running status of the ServicePackage (SP), it means the SP is running
	RunningPhase = "Running"
	// SucceededPhase succeeded status of the ServicePackage (SP), it means the SP has already deployed succeeded
	SucceededPhase = "Succeeded"
	// FailedPhase failed status of the ServicePackage (SP), it means the SP is deploying or upgrading failed
	FailedPhase = "Failed"
	// UnknownPhase unknown status of the ServicePackage (SP), it means the SP has unknown reasons to run with exception
	UnknownPhase = "Unknown"
	// UpgradingPhase upgrading status of the ServicePackage (SP), it means the SP is upgrading
	UpgradingPhase = "Upgrading"
	// DeletingPhase deleting status of the ServicePackage (SP), it means the SP is during deleting
	DeletingPhase = "Deleting"
	// DeletedPhase deleted status of the ServicePackage (SP), it means the SP has already deleted
	DeletedPhase = "Deleted"
)

// ServicePackageSpec defines the desired state of ServicePackage
type ServicePackageSpec struct {
	// ServiceID is the unique id for the service which from service package
	ServiceID string `json:"serviceID"`
	// Name is the service name
	Name string `json:"name"`
	// Version is the version of this service
	Version string `json:"version,omitempty"`
	// Resources is a base64 binary code, it will be analysis by engine to the service instance resources,
	// such as custom resource, cluster role, cluster role binding, and workload etc.
	Resources string `json:"resources,omitempty"`
	// RawResources is a slice of non-kappital packages resources
	RawResources []RawResource `json:"rawResources,omitempty"`
}

// RawResource is the resources for the service except Kappital’s package, such as helm, etc.
type RawResource struct {
	Type string `json:"type"`
	Raw  string `json:"raw"`
}

// ServicePackageStatus defines the observed state of ServicePackage
type ServicePackageStatus struct {
	CurrentVersion   string       `json:"currentVersion,omitempty"`
	Phase            string       `json:"phase,omitempty"`
	Reason           string       `json:"reason,omitempty"`
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

// ServicePackage is the Schema for the servicepackages API
type ServicePackage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServicePackageSpec   `json:"spec,omitempty"`
	Status ServicePackageStatus `json:"status,omitempty"`
}

// ServicePackageList contains a list of ServicePackage
type ServicePackageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServicePackage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServicePackage{}, &ServicePackageList{})
}

// VerifyStatus will verify the service package status
func (in *ServicePackage) VerifyStatus() {
	in.Status.Reason = ""
	if len(in.Status.CurrentVersion) == 0 && !in.IsDeleting() {
		in.SetToPending()
		return
	}
	if len(in.Spec.Version) == 0 {
		in.SetToDeleting()
		return
	}
	if in.Spec.Version != in.Status.CurrentVersion && !in.IsDeleting() {
		in.SetToUpgrading()
	}
}

// SetToPending set the service package to pending, which means the service engine is installing this service package.
func (in *ServicePackage) SetToPending() {
	in.Status.Phase = PendingPhase
	in.Status.CurrentVersion = in.Spec.Version
	in.Status.Reason = ""
	in.Status.LastScheduleTime = &metav1.Time{Time: time.Now().UTC()}
}

// SetToRunning set the service package to running, which means the service package is running.
func (in *ServicePackage) SetToRunning() {
	in.Status.Phase = RunningPhase
	in.Status.Reason = ""
	in.Status.CurrentVersion = in.Spec.Version
	in.Status.LastScheduleTime = &metav1.Time{Time: time.Now().UTC()}
}

// SetToSucceeded set the service package to succeeded, which means the service package is deployed success or
// upgraded success.
func (in *ServicePackage) SetToSucceeded() {
	in.Status.Phase = SucceededPhase
	in.Status.Reason = ""
	in.Status.CurrentVersion = in.Spec.Version
	in.Status.LastScheduleTime = &metav1.Time{Time: time.Now().UTC()}
}

// SetToFailed set the service package to failed, which means the service package is deployed failed or upgraded failed.
func (in *ServicePackage) SetToFailed(reason string) {
	if in.isException() {
		in.Status.Reason = fmt.Sprintf("%s; %s", in.Status.Reason, reason)
	} else {
		in.Status.Reason = reason
	}
	in.Status.Phase = FailedPhase
	in.Status.LastScheduleTime = &metav1.Time{Time: time.Now().UTC()}
}

// SetToUnknown set the service package to unknown, which means the service package has unknown reasons make it may not
// provide the all functions.
func (in *ServicePackage) SetToUnknown(reason string) {
	if in.isException() {
		in.Status.Reason = fmt.Sprintf("%s; %s", in.Status.Reason, reason)
	} else {
		in.Status.Reason = reason
	}
	if in.Status.Phase != FailedPhase {
		in.Status.Phase = UnknownPhase
	}
	in.Status.LastScheduleTime = &metav1.Time{Time: time.Now().UTC()}
}

// SetToUpgrading set to the service package to upgrading, which means the service package is during the upgrading.
func (in *ServicePackage) SetToUpgrading() {
	in.Status.Phase = UpgradingPhase
	in.Status.Reason = "Upgrade the operator"
	in.Status.LastScheduleTime = &metav1.Time{Time: time.Now().UTC()}
}

// SetToDeleting set the service package to deleting, which means the service package is during the deleting.
func (in *ServicePackage) SetToDeleting() {
	in.Status.Phase = DeletingPhase
	now := time.Now().UTC()
	in.Status.Reason = fmt.Sprintf("begin [%s] to delete the service instance [%s]", now, in.Spec.Name)
	in.Status.LastScheduleTime = &metav1.Time{Time: now}
}

// SetToDeleted set the service package to deleted, which means the service package is already deleted. The service
// package will wait for manage do the garbage collection.
func (in *ServicePackage) SetToDeleted() {
	in.Status.Phase = DeletedPhase
	in.Status.CurrentVersion = ""
	now := time.Now().UTC()
	in.Status.Reason = fmt.Sprintf("at [%s] the service instance [%s] has been deleted", now, in.Spec.Name)
	in.Status.LastScheduleTime = &metav1.Time{Time: now}
}

// UpdateStatus update the service package with errors. This is the final step for reconciler. It will update the
// status to the correct status for pre-processes.
func (in *ServicePackage) UpdateStatus(err error) {
	if err == nil {
		switch in.Status.Phase {
		case PendingPhase, UpgradingPhase:
			in.SetToSucceeded()
		case DeletingPhase:
			in.SetToDeleted()
		}
		return
	}

	in.SetToFailed(err.Error())
}

// NeedCheckRuntime does the service package is need to check the application objects runtime status.
func (in ServicePackage) NeedCheckRuntime() bool {
	return !in.IsDeleting() && !in.IsUpgrading() && in.Status.Phase != PendingPhase
}

// IsDeleting does the service package during the deleting relation status.
func (in ServicePackage) IsDeleting() bool {
	return in.Status.Phase == DeletingPhase || in.Status.Phase == DeletedPhase
}

// IsDeleted does the service package is deleted.
func (in ServicePackage) IsDeleted() bool {
	return in.Status.Phase == DeletedPhase
}

// IsUpgrading does the service package during the upgrading status.
func (in ServicePackage) IsUpgrading() bool {
	return in.Status.Phase == UpgradingPhase || (in.Status.CurrentVersion != in.Spec.Version && in.Spec.Version != "")
}

func (in ServicePackage) isException() bool {
	return in.Status.Phase == FailedPhase || in.Status.Phase == UnknownPhase
}
