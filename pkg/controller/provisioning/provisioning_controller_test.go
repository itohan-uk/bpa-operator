package provisioning

import (

       "testing"

       bpav1alpha1 "github.com/bpa-operator/pkg/apis/bpa/v1alpha1"
       metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
       logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
       "k8s.io/apimachinery/pkg/runtime"
       "k8s.io/apimachinery/pkg/types"
       "k8s.io/client-go/kubernetes/scheme"
       "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
       "sigs.k8s.io/controller-runtime/pkg/client/fake"
       "sigs.k8s.io/controller-runtime/pkg/reconcile"
       fakedynamic "k8s.io/client-go/dynamic/fake"
       fakeclientset "k8s.io/client-go/kubernetes/fake"
)

func TestProvisioningController(t *testing.T) {

     logf.SetLogger(logf.ZapLogger(true))
     name := "bpa-test-cr"
     namespace := "default"
     clusterName := "test-cluster"

     // Create Fake baremetalhost
     bmhList := newBMList()

    // Create Fake Provisioning CR
    provisioning := newBPA(name, namespace, clusterName)

    // Objects to track in the fake Client
    objs := []runtime.Object{provisioning}

    // Register operator types with the runtime scheme
    sc := scheme.Scheme

    sc.AddKnownTypes(bpav1alpha1.SchemeGroupVersion, provisioning)

    // Create Fake Clients and Clientset
    fakeClient := fake.NewFakeClient(objs...)
    fakeDyn := fakedynamic.NewSimpleDynamicClient(sc, bmhList,)
    fakeClientSet := fakeclientset.NewSimpleClientset()

    r := &ReconcileProvisioning{client: fakeClient, scheme: sc, clientset: fakeClientSet, bmhClient: fakeDyn}

    // Mock request to simulate Reconcile() being called on an event for a watched resource 

    req := simulateRequest(provisioning)
    _, err := r.Reconcile(req)
    if err != nil {
       t.Fatalf("reconcile: (%v)", err)
    }

   jobClient := r.clientset.BatchV1().Jobs(namespace)
   job, err := jobClient.Get("kud-test-cluster", metav1.GetOptions{})

    if err != nil {
        t.Fatalf("Error occured while getting job: (%v)", err)
    }

   jobClusterName := job.Labels["cluster"]
   if jobClusterName != clusterName {
      t.Fatalf("Job cluster Name is wrong")
   }

}

func simulateRequest(bpaCR *bpav1alpha1.Provisioning) reconcile.Request {
	namespacedName := types.NamespacedName{
		Name:      bpaCR.ObjectMeta.Name,
		Namespace: bpaCR.ObjectMeta.Namespace,
	}
	return reconcile.Request{NamespacedName: namespacedName}
}



func newBPA(name, namespace, clusterName string) *bpav1alpha1.Provisioning {

     provisioningCR := &bpav1alpha1.Provisioning{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: namespace,
            Labels: map[string]string{
                "cluster": clusterName,
            },
        },
        Spec: bpav1alpha1.ProvisioningSpec{
               Masters: []map[string]bpav1alpha1.Master{
                         map[string]bpav1alpha1.Master{
                           "test-master" : bpav1alpha1.Master{
                                 MACaddress: "08:00:27:00:ab:2c",
                            },

               },
              },
       },

    }
    return provisioningCR
}


func newBMList() *unstructured.UnstructuredList{

	bmMap := map[string]interface{}{
			   "apiVersion": "metal3.io/v1alpha1",
			   "kind": "BareMetalHostList",
			   "metaDatai": map[string]interface{}{
			       "continue": "",
				   "resourceVersion": "11830058",
				   "selfLink": "/apis/metal3.io/v1alpha1/baremetalhosts",

		 },
		 }




	metaData := map[string]interface{}{
			 "creationTimestamp": "2019-10-24T04:51:15Z",
			 "generation":"1",
			 "name": "fake-test-bmh",
			 "namespace": "default",
			 "resourceVersion": "11829263",
			 "selfLink": "/apis/metal3.io/v1alpha1/namespaces/default/baremetalhosts/bpa-test-bmh",
			 "uid": "e92cb312-f619-11e9-90bc-00219ba0c77a",
	}



	nicMap1 := map[string]interface{}{
			"ip": "",
			 "mac": "08:00:27:00:ab:2c",
			 "model": "0x8086 0x1572",
			 "name": "eth3",
			 "pxe": "false",
			 "speedGbps": "0",
			 "vlanId": "0",
	}

	specMap  := map[string]interface{}{
			  "status" : map[string]interface{}{
				   "errorMessage": "",
					"hardware": map[string]interface{}{
					   "nics": nicMap1,
			  },
			  },


	}

	itemMap := map[string]interface{}{
			   "apiVersion": "metal3.io/v1alpha1",
			   "kind": "BareMetalHost",
			   "metadata": metaData,
			   "spec": specMap,
		 }
	itemU := unstructured.Unstructured{
			 Object: itemMap,
		   }

	itemsList := []unstructured.Unstructured{itemU,}

	bmhList := &unstructured.UnstructuredList{
					Object: bmMap,
					Items: itemsList,
	 }


      return bmhList
}


// Create DHCP file for testing

