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

package engine

import (
	"context"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
)

// reconcilePermission will process the permission-related resources, such as service account, cluster role, and
// cluster role binding.
func (r *ServicePackageReconciler) reconcilePermission(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	permissions []enginev1alpha1.Permission) error {
	if err := r.reconcileServiceAccount(ctx, pack, permissions); err != nil {
		return err
	}
	return r.reconcileClusterRoleAndBinding(ctx, pack, permissions)
}
