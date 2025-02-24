/*
Copyright 2023 Flant JSC

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

package apis

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	"github.com/slok/kubewebhook/v2/pkg/model"
	kwhvalidating "github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/deckhouse/deckhouse/go_lib/module"
	"github.com/deckhouse/deckhouse/modules/002-deckhouse/hooks/pkg/apis/v1alpha1"
)

func init() {
	module.RegisterValidationHandler("/validate/v1alpha1/modules", moduleValidationHandler())
}

func moduleValidationHandler() http.Handler {
	vf := kwhvalidating.ValidatorFunc(func(ctx context.Context, review *model.AdmissionReview, obj metav1.Object) (result *kwhvalidating.ValidatorResult, err error) {
		// UserInfo groups: [system:serviceaccounts system:serviceaccounts:d8-system system:authenticated]
		if review.UserInfo.Username != "system:serviceaccount:d8-system:deckhouse" {
			return &kwhvalidating.ValidatorResult{
				Valid:   false,
				Message: "manual Module change is forbidden",
			}, nil
		}

		return &kwhvalidating.ValidatorResult{
			Valid:   true,
			Message: "",
		}, nil
	})

	kl := kwhlogrus.NewLogrus(log.NewEntry(log.StandardLogger()))

	// Create webhook.
	wh, _ := kwhvalidating.NewWebhook(kwhvalidating.WebhookConfig{
		ID:        "module-operations",
		Validator: vf,
		Logger:    kl,
		Obj:       &v1alpha1.Module{},
	})

	return kwhhttp.MustHandlerFor(kwhhttp.HandlerConfig{Webhook: wh, Logger: kl})
}
