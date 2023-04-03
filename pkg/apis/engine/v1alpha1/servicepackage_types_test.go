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
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServicePackage_IsDeleted(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       ServicePackageSpec
		Status     ServicePackageStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "ServicePackage IsDeleted (false)",
			fields: fields{
				Status: ServicePackageStatus{Phase: PendingPhase},
			},
			want: false,
		},
		{
			name: "ServicePackage IsDeleted (true)",
			fields: fields{
				Status: ServicePackageStatus{Phase: DeletedPhase},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := ServicePackage{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := in.IsDeleted(); got != tt.want {
				t.Errorf("IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServicePackage_NeedCheckRuntime(t *testing.T) {
	type fields struct {
		Status ServicePackageStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "ServicePackage NeedCheckRuntime (false for deleting status)",
			fields: fields{Status: ServicePackageStatus{Phase: DeletingPhase}},
		},
		{
			name:   "ServicePackage NeedCheckRuntime (false for deleted status)",
			fields: fields{Status: ServicePackageStatus{Phase: DeletedPhase}},
		},
		{
			name:   "ServicePackage NeedCheckRuntime (false for upgrading status)",
			fields: fields{Status: ServicePackageStatus{Phase: UpgradingPhase}},
		},
		{
			name:   "ServicePackage NeedCheckRuntime (false for pending status)",
			fields: fields{Status: ServicePackageStatus{Phase: PendingPhase}},
		},
		{
			name:   "ServicePackage NeedCheckRuntime (true for running status)",
			fields: fields{Status: ServicePackageStatus{Phase: RunningPhase}},
			want:   true,
		},
		{
			name:   "ServicePackage NeedCheckRuntime (true for succeeded status)",
			fields: fields{Status: ServicePackageStatus{Phase: SucceededPhase}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := ServicePackage{
				Status: tt.fields.Status,
			}
			if got := in.NeedCheckRuntime(); got != tt.want {
				t.Errorf("NeedCheckRuntime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServicePackage_SetToDeleted(t *testing.T) {
	tests := []struct {
		name      string
		wantPhase string
	}{
		{
			name:      "ServicePackage SetToDeleted",
			wantPhase: DeletedPhase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &ServicePackage{}
			in.SetToDeleted()
			if in.Status.Phase != DeletedPhase {
				t.Errorf("SetToDeleted() = %v, want %v", in.Status.Phase, DeletedPhase)
			}
		})
	}
}

func TestServicePackage_SetToFailed(t *testing.T) {
	type fields struct {
		Status ServicePackageStatus
	}
	tests := []struct {
		name       string
		fields     fields
		reason     string
		wantReason string
	}{
		{
			name:       "ServicePackage SetToFailed (exception status of failed)",
			fields:     fields{Status: ServicePackageStatus{Phase: FailedPhase, Reason: "x0"}},
			reason:     "x1",
			wantReason: "x0; x1",
		},
		{
			name:       "ServicePackage SetToFailed (exception status of unknown)",
			fields:     fields{Status: ServicePackageStatus{Phase: UnknownPhase, Reason: "x0"}},
			reason:     "x1",
			wantReason: "x0; x1",
		},
		{
			name:       "ServicePackage SetToFailed (non-exception status)",
			fields:     fields{Status: ServicePackageStatus{Phase: PendingPhase}},
			reason:     "x1",
			wantReason: "x1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &ServicePackage{
				Status: tt.fields.Status,
			}
			in.SetToFailed(tt.reason)
			if in.Status.Phase != FailedPhase || in.Status.Reason != tt.wantReason {
				t.Errorf("SetToFailed() = %v, wantPhase %v wantReason %v", in, FailedPhase, tt.wantReason)
			}
		})
	}
}

func TestServicePackage_SetToRunning(t *testing.T) {
	tests := []struct {
		name      string
		wantPhase string
	}{
		{
			name:      "ServicePackage SetToRunning",
			wantPhase: RunningPhase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &ServicePackage{}
			in.SetToRunning()
			if in.Status.Phase != RunningPhase {
				t.Errorf("SetToPending() = %v, want %v", in.Status.Phase, RunningPhase)
			}
		})
	}
}

func TestServicePackage_SetToSucceeded(t *testing.T) {
	type fields struct {
		Spec ServicePackageSpec
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "ServicePackage SetToSucceeded",
			fields: fields{Spec: ServicePackageSpec{Version: "x1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &ServicePackage{
				Spec: tt.fields.Spec,
			}
			in.SetToSucceeded()
			if in.Status.Phase != SucceededPhase ||
				in.Status.CurrentVersion != in.Spec.Version {
				t.Errorf("SetToSucceeded() = %v, wantPhase %v and wantVersion %v",
					in, SucceededPhase, in.Spec.Version)
			}
		})
	}
}

func TestServicePackage_SetToUnknown(t *testing.T) {
	type fields struct {
		Status ServicePackageStatus
	}
	tests := []struct {
		name   string
		fields fields
		reason string
		want   ServicePackageStatus
	}{
		{
			name:   "ServicePackage SetToFailed (failed status)",
			fields: fields{Status: ServicePackageStatus{Phase: FailedPhase, Reason: "x0"}},
			reason: "x1",
			want:   ServicePackageStatus{Phase: FailedPhase, Reason: "x0; x1"},
		},
		{
			name:   "ServicePackage SetToFailed (unknown status)",
			fields: fields{Status: ServicePackageStatus{Phase: UnknownPhase, Reason: "x0"}},
			reason: "x1",
			want:   ServicePackageStatus{Phase: UnknownPhase, Reason: "x0; x1"},
		},
		{
			name:   "ServicePackage SetToFailed (non-exception status)",
			fields: fields{Status: ServicePackageStatus{Phase: SucceededPhase}},
			reason: "x1",
			want:   ServicePackageStatus{Phase: UnknownPhase, Reason: "x1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &ServicePackage{
				Status: tt.fields.Status,
			}
			in.SetToUnknown(tt.reason)
			if in.Status.Phase != tt.want.Phase || in.Status.Reason != tt.want.Reason {
				t.Errorf("SetToUnknown() = %v, wantStatus %v", in, tt.want)
			}
		})
	}
}

func TestServicePackage_UpdateStatus(t *testing.T) {
	type fields struct {
		Status ServicePackageStatus
	}
	tests := []struct {
		name   string
		fields fields
		err    error
		want   ServicePackageStatus
	}{
		{
			name: "ServicePackage UpdateStatus (err != nil)",
			err:  fmt.Errorf("error"),
			want: ServicePackageStatus{Phase: FailedPhase},
		},
		{
			name:   "ServicePackage UpdateStatus (pending phase)",
			fields: fields{Status: ServicePackageStatus{Phase: PendingPhase}},
			want:   ServicePackageStatus{Phase: SucceededPhase},
		},
		{
			name:   "ServicePackage UpdateStatus (upgrading phase)",
			fields: fields{Status: ServicePackageStatus{Phase: UpgradingPhase}},
			want:   ServicePackageStatus{Phase: SucceededPhase},
		},
		{
			name:   "ServicePackage UpdateStatus (deleting phase)",
			fields: fields{Status: ServicePackageStatus{Phase: DeletingPhase}},
			want:   ServicePackageStatus{Phase: DeletedPhase},
		},
		{
			name:   "ServicePackage UpdateStatus (running phase)",
			fields: fields{Status: ServicePackageStatus{Phase: RunningPhase}},
			want:   ServicePackageStatus{Phase: RunningPhase},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &ServicePackage{
				Status: tt.fields.Status,
			}
			in.UpdateStatus(tt.err)
			if in.Status.Phase != tt.want.Phase {
				t.Errorf("UpdateStatus() = %v, wantStatus %v", in, tt.want)
			}
		})
	}
}

func TestServicePackage_VerifyStatus(t *testing.T) {
	type fields struct {
		Spec   ServicePackageSpec
		Status ServicePackageStatus
	}
	tests := []struct {
		name      string
		fields    fields
		wantPhase string
	}{
		{
			name:      "ServicePackage VerifyStatus (to pending)",
			fields:    fields{Spec: ServicePackageSpec{Version: "x1"}},
			wantPhase: PendingPhase,
		},
		{
			name:      "ServicePackage VerifyStatus (to deleting)",
			fields:    fields{Status: ServicePackageStatus{CurrentVersion: "x2"}},
			wantPhase: DeletingPhase,
		},
		{
			name: "ServicePackage VerifyStatus (to upgrading)",
			fields: fields{
				Spec:   ServicePackageSpec{Version: "x1"},
				Status: ServicePackageStatus{CurrentVersion: "x2"},
			},
			wantPhase: UpgradingPhase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &ServicePackage{
				Spec:   tt.fields.Spec,
				Status: tt.fields.Status,
			}
			in.VerifyStatus()
		})
	}
}
