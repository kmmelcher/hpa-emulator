package main

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/meta/testrestmapper"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"
	metricsfake "k8s.io/metrics/pkg/client/clientset/versioned/fake"
	cmfake "k8s.io/metrics/pkg/client/custom_metrics/fake"
	emfake "k8s.io/metrics/pkg/client/external_metrics/fake"

	scalefake "k8s.io/client-go/scale/fake"

	hpa "k8s.io/kubernetes/pkg/controller/podautoscaler"
)


func main(){

	ctx := context.TODO()

	fakeEventClient := &fake.Clientset{}
	fakeScaleClient := &scalefake.FakeScaleClient{}
	fakeClient := &fake.Clientset{}
	restMapper := testrestmapper.TestOnlyStaticRESTMapper(scheme.Scheme) // legacyscheme?

	fakeMetricsClient := &metricsfake.Clientset{}
	fakeCMClient := &cmfake.FakeCustomMetricsClient{}
	fakeEMClient := &emfake.FakeExternalMetricsClient{}
	metricsClient := metrics.NewRESTMetricsClient(
		fakeMetricsClient.MetricsV1beta1(),
		fakeCMClient,
		fakeEMClient,
	)
	resyncPeriod := 1 * time.Second
	informerFactory := informers.NewSharedInformerFactory(fakeClient, resyncPeriod)
	fakeHpaInformer := informerFactory.Autoscaling().V2().HorizontalPodAutoscalers()
	fakePodInformer := informerFactory.Core().V1().Pods()
	downScaleStabilizationWindow := 5 * time.Minute
	tolerance := 0.1
	cpuInitializationPeriod := 2 * time.Minute
	delayOfInitialReadinessStatus := 10 * time.Second

	hpaController := hpa.NewHorizontalController(
		ctx,
		fakeEventClient.CoreV1(),
		fakeScaleClient, 
		fakeClient.AutoscalingV2(), 
		restMapper, 
		metricsClient, 
		fakeHpaInformer, 
		fakePodInformer, 
		resyncPeriod, 
		downScaleStabilizationWindow, 
		tolerance, 
		cpuInitializationPeriod,
		delayOfInitialReadinessStatus,
	)

	hpaController.Run(ctx, 1)
}
