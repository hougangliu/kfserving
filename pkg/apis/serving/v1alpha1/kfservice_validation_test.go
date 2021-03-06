/*
Copyright 2019 kubeflow.org.

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

package v1alpha1

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestRejectMultipleModelSpecs(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	kfsvc := TFExampleKFService.DeepCopy()
	kfsvc.Spec.Default.ScikitLearn = &ScikitLearnSpec{ModelURI: "gs://testbucket/testmodel"}
	g.Expect(kfsvc.ValidateCreate()).Should(gomega.MatchError("Exactly one of [Custom, Tensorflow, ScikitLearn, XGBoost] may be specified in ModelSpec"))
}
func TestRejectModelSpecMissing(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	kfsvc := TFExampleKFService.DeepCopy()
	kfsvc.Spec.Default.Tensorflow = nil
	g.Expect(kfsvc.ValidateCreate()).Should(gomega.MatchError("Exactly one of [Custom, Tensorflow, ScikitLearn, XGBoost] may be specified in ModelSpec"))
}
func TestRejectMultipleCanaryModelSpecs(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	kfsvc := TFExampleKFService.DeepCopy()
	kfsvc.Spec.Canary = &CanarySpec{ModelSpec: ModelSpec{
		ScikitLearn: &ScikitLearnSpec{ModelURI: "gs://testbucket/testmodel"},
		Tensorflow:  kfsvc.Spec.Default.Tensorflow,
	}}
	g.Expect(kfsvc.ValidateCreate()).Should(gomega.MatchError("Exactly one of [Custom, Tensorflow, ScikitLearn, XGBoost] may be specified in ModelSpec"))
}

func TestRejectCanaryModelSpecMissing(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	kfsvc := TFExampleKFService.DeepCopy()
	kfsvc.Spec.Canary = &CanarySpec{ModelSpec: ModelSpec{}}
	g.Expect(kfsvc.ValidateCreate()).Should(gomega.MatchError("Exactly one of [Custom, Tensorflow, ScikitLearn, XGBoost] may be specified in ModelSpec"))
}

func TestRejectBadCanaryTrafficValues(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	kfsvc := TFExampleKFService.DeepCopy()
	kfsvc.Spec.Canary = &CanarySpec{
		TrafficPercent: -1,
		ModelSpec:      kfsvc.Spec.Default,
	}
	g.Expect(kfsvc.ValidateCreate()).Should(gomega.MatchError("TrafficPercent must be between [0, 100]"))
	kfsvc.Spec.Canary.TrafficPercent = 101
	g.Expect(kfsvc.ValidateCreate()).Should(gomega.MatchError("TrafficPercent must be between [0, 100]"))
}

func TestBadReplicaValues(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	kfsvc := TFExampleKFService.DeepCopy()
	kfsvc.Spec.MinReplicas = -1
	g.Expect(kfsvc.ValidateCreate()).Should(gomega.MatchError("MinReplicas cannot be less than 0"))
	kfsvc.Spec.MinReplicas = 1
	kfsvc.Spec.MaxReplicas = -1
	g.Expect(kfsvc.ValidateCreate()).Should(gomega.MatchError("MaxReplicas cannot be less than 0"))
	kfsvc.Spec.MinReplicas = 2
	kfsvc.Spec.MaxReplicas = 1
	g.Expect(kfsvc.ValidateCreate()).Should(gomega.MatchError("MinReplicas cannot be greater than MaxReplicas"))
}
