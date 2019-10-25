package provisioning

import (

       "fmt"
       "context"
       "testing"

       bpav1alpha1 "github.com/bpa-operator/pkg/apis/bpa/v1alpha1"
       metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
       //corev1 "k8s.io/api/core/v1"
       batchv1 "k8s.io/api/batch/v1"
       logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
       metal3v1alpha1 "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
       metal3 "github.com/metal3-io/baremetal-operator/pkg/apis"
       "k8s.io/apimachinery/pkg/runtime"
       "k8s.io/apimachinery/pkg/types"
       "k8s.io/client-go/kubernetes/scheme"
       "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
       //"sigs.k8s.io/controller-runtime/pkg/client"
       "sigs.k8s.io/controller-runtime/pkg/client/fake"
       "sigs.k8s.io/controller-runtime/pkg/reconcile"
       fakedynamic "k8s.io/client-go/dynamic/fake"
       fakeclientset "k8s.io/client-go/kubernetes/fake"
)

func TestProvisioningController(t *testing.T) {

     logf.SetLogger(logf.ZapLogger(true))
     name := "bpa-test-cr"
     namespace := "default"
     masterMap := make(map[string]bpav1alpha1.Master)

     masterMap["test-master"] = bpav1alpha1.Master{
                                    MACaddress: "08:00:27:00:ab:c0",
                                 }

     masterList := make([]map[string]bpav1alpha1.Master, 1)
     masterList[0] = masterMap

     //bmhHost :=  newHost()

     bmhList := newHostList()
     //bmhList := newBMList()
     //t.Logf("\n\n%+v\n\n UUUUUUUU", bmhListLLL)



     bpaSpec := &bpav1alpha1.ProvisioningSpec{
               Masters: masterList,

               }

    provisioning := newBPA(name, namespace, "test-cluster", bpaSpec)


    // Objects to track in the fake Client
    objs := []runtime.Object{provisioning}

    // Register operator types with the runtime scheme
    sc := scheme.Scheme

    // Add baremetalHost to Scheme
    if err := metal3.AddToScheme(sc); err != nil {
        t.Fatalf("Unable to add BareMetalHost scheme: (%v)", err)
    }
    sc.AddKnownTypes(bpav1alpha1.SchemeGroupVersion, provisioning)

    //Create Fake Client and Clientset
    fakeClient := fake.NewFakeClient(objs...)
    fakeDyn := fakedynamic.NewSimpleDynamicClient(sc, bmhList )
    fakeClientSet := fakeclientset.NewSimpleClientset()


    r := &ReconcileProvisioning{client: fakeClient, scheme: sc, clientset: fakeClientSet, bmhClient: fakeDyn}


   // Mock request to simulate Reconcile() being called on an event for a
    // watched resource .
    req := reconcile.Request{
        NamespacedName: types.NamespacedName{
            Name:      name,
            Namespace: namespace,
        },
    }

   _, err := r.Reconcile(req)
    if err != nil {
       t.Fatalf("reconcile: (%v)", err)
    }




   jb := &batchv1.Job{}
   jobName := types.NamespacedName{
		Name: "kud-test-cluster",
		Namespace: "default",
    }

   err = fakeClient.Get(context.TODO(), jobName, jb)
    if err != nil {
        t.Fatalf("Error occured while getting job: (%v)", err)
    }

   fmt.Printf("Got job...ending now!!")

}

func simulateRequest(bpaCR *bpav1alpha1.Provisioning) reconcile.Request {
	namespacedName := types.NamespacedName{
		Name:      bpaCR.ObjectMeta.Name,
		Namespace: bpaCR.ObjectMeta.Namespace,
	}
	return reconcile.Request{NamespacedName: namespacedName}
}

func newBPA(name, namespace, clusterName string, spec *bpav1alpha1.ProvisioningSpec) *bpav1alpha1.Provisioning {

     provisioningCR := &bpav1alpha1.Provisioning{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: namespace,
            Labels: map[string]string{
                "cluster": clusterName,
            },
        },
        Spec: *spec,
    }

    return provisioningCR
}

//func newHost(name string, spec *metal3v1alpha1.BareMetalHostSpec, status *metal3v1alpha1.BareMetalHostStatus) *metal3v1alpha1.BareMetalHost {
func newHostList() *metal3v1alpha1.BareMetalHostList {


    nicList := make([]metal3v1alpha1.NIC, 2)

                                nic1 := metal3v1alpha1.NIC{
                                          Name: "eth1",
                                          Model: "0x80860x1572",
                                          MAC: "08:00:27:00:ab:c0",
                                          IP: "",
                                }

                                nic2 := metal3v1alpha1.NIC{
                                      Name: "eth2",
                                          Model: "0x80860x37d2",
                                          MAC: "a4:bf:01:64:86:6e",
                                          IP: "",
                                }

                nicList[0] = nic1
                nicList[1] = nic2


    bmh := &metal3v1alpha1.BareMetalHostList{Items: []metal3v1alpha1.BareMetalHost{{


                ObjectMeta: metav1.ObjectMeta{
                        Name:      "test-bmh-1",
                        Namespace: "default",
                },
                Spec: metal3v1alpha1.BareMetalHostSpec{
                        BMC: metal3v1alpha1.BMCDetails{
                                Address:         "",
                                CredentialsName: "bmc-creds-valid",
                        },
                },
                Status: metal3v1alpha1.BareMetalHostStatus{
                          HardwareDetails: &metal3v1alpha1.HardwareDetails {
                                         NIC: nicList,
                                          Hostname: "fake-host",
                          },

                            },

                                },}}

    return bmh
}



func newHost()  *metal3v1alpha1.BareMetalHost {

      nicList := make([]metal3v1alpha1.NIC, 2)

                                nic1 := metal3v1alpha1.NIC{
                                          Name: "eth1",
                                          Model: "0x80860x1572",
                                          MAC: "08:00:27:00:ab:c0",
                                          IP: "",
                                }

                                nic2 := metal3v1alpha1.NIC{
                                      Name: "eth2",
                                          Model: "0x80860x37d2",
                                          MAC: "a4:bf:01:64:86:6e",
                                          IP: "",
                                }

                nicList[0] = nic1
                nicList[1] = nic2

     bmh := &metal3v1alpha1.BareMetalHost{



                ObjectMeta: metav1.ObjectMeta{
                        Name:      "test-bmh-2",
                        Namespace: "default",
                },
                Spec: metal3v1alpha1.BareMetalHostSpec{
                        BMC: metal3v1alpha1.BMCDetails{
                                Address:         "",
                                CredentialsName: "bmc-creds-valid",
                        },
                },
                Status: metal3v1alpha1.BareMetalHostStatus{
                          HardwareDetails: &metal3v1alpha1.HardwareDetails {
                                         NIC: nicList,
                                          Hostname: "fake-host",
                          },

                            },

                                }
         return bmh
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
			 "name": "bpa-test-bmh",
			 "namespace": "default",
			 "resourceVersion": "11829263",
			 "selfLink": "/apis/metal3.io/v1alpha1/namespaces/default/baremetalhosts/bpa-test-bmh",
			 "uid": "e92cb312-f619-11e9-90bc-00219ba0c77a",
	}



	nicMap1 := map[string]interface{}{
			"ip": "",
			 "mac": "08:00:27:00:ab:c0",
			 "model": "0x8086 0x1572",
			 "name": "eth3",
			 "pxe": "false",
			 "speedGbps": "0",
			 "vlanId": "0",
	}

	nicMap2 := map[string]interface{}{
			"ip": "",
			 "mac": "a4:bf:01:64:86:6e",
			 "model": "0x8086 0x37d2",
			 "name": "eth4",
			 "pxe": "false",
			 "speedGbps": "0",
			 "vlanId": "0",
	}

	nicList := []map[string]interface{}{
			   nicMap1,
			   nicMap2,
			 }

	specMap  := map[string]interface{}{
			  "status" : map[string]interface{}{
				   "errorMessage": "",
					"hardware": map[string]interface{}{
					   "nics": nicList,
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
