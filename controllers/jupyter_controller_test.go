package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	operatorsv2 "convect.ai/notebook-crd/api/v2"
)

var _ = Describe("Notebook controller", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		Name      = "test-notebook"
		Namespace = "default"
		timeout   = time.Second * 10
		interval  = time.Millisecond * 250
	)

	Context("When validating the notebook controller", func() {
		It("Should create replicas", func() {
			By("By creating a new notebook")
			ctx := context.Background()
			notebook := &operatorsv2.Jupyter{
				ObjectMeta: metav1.ObjectMeta{
					Name:      Name,
					Namespace: Namespace,
				},
				Spec: operatorsv2.JupyterSpec{
					Template: operatorsv2.JupyterTemplate{
						Spec: v1.PodSpec{
							Containers: []v1.Container{{
								Name:  "busybox",
								Image: "busybox",
							}},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, notebook)).Should(Succeed())

			notebookLookupKey := types.NamespacedName{Name: Name, Namespace: Namespace}
			createdNotebook := &operatorsv2.Jupyter{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, notebookLookupKey, createdNotebook)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("By checking that the notebook has statefulset")
			Eventually(func() (bool, error) {
				sts := &appsv1.StatefulSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      Name,
						Namespace: Namespace,
					},
				}

				err := k8sClient.Get(ctx, notebookLookupKey, sts)
				if err != nil {
					return false, err
				}
				return true, nil
			}, timeout, interval).Should(BeTrue())
		})
	})
})
