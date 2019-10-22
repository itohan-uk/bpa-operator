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
       "k8s.io/apimachinery/pkg/runtime"
       "k8s.io/apimachinery/pkg/types"
       "k8s.io/client-go/kubernetes/scheme"
       //"sigs.k8s.io/controller-runtime/pkg/client"
       "sigs.k8s.io/controller-runtime/pkg/client/fake"
       "sigs.k8s.io/controller-runtime/pkg/reconcile"
       //fakeclientset "k8s.io/client-go/kubernetes/fake"
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

     bmhHost :=  newHost()

     bpaSpec := &bpav1alpha1.ProvisioningSpec{
               Masters: masterList,

               }

    provisioning := newBPA(name, namespace, "test-cluster", bpaSpec)

    objs := []runtime.Object{provisioning,}

    // Register operatory tyes with the runtime schem
    sc := scheme.Scheme

    /*if err := metal3v1alpha1.AddToScheme(sc); err != nil {
        t.Fatalf("Unable to add baremetalhost scheme: (%v)", err)
    }*/
    sc.AddKnownTypes(bpav1alpha1.SchemeGroupVersion, provisioning, bmhHost)



    // Create a fake client
    //fakeClient := fake.NewFakeClient(objs...)
    fakeClient := fake.NewFakeClientWithScheme(sc, objs...)

    err := fakeClient.Create(context.TODO(), bmhHost)
    if err != nil {
       t.Fatalf("Error occured while create baremetal host: (%v)", err)
    }

    bmh := &metal3v1alpha1.BareMetalHost{}
    err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: bmhHost.Name, Namespace: bmhHost.Namespace}, bmh)
    if err != nil {
        t.Fatalf("Error occured while getting baremetalhost: (%v)", err)
    }

   fmt.Printf("found bmh\n")
    provisioningObj := &ReconcileProvisioning{client: fakeClient, scheme: sc}


    req := simulateRequest(provisioning)


    res, err := provisioningObj.Reconcile(req)

    if err != nil {
       t.Fatalf("reconcile: (%v)", err)
    }

    if res != (reconcile.Result{}) {
        t.Error("reconcile did not return an empty Result")
    }


   //expectedJob := createKUDinstallerJob("test-cluster", "default", map[string]string{ "cluster": "test-cluster", }, fakeClientSet)
   // Check if a job was created
   /*bmh := &metal3v1alpha1.BareMetalHost{}
   err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: bmhHost.Name, Namespace: bmhHost.Namespace}, bmh)
    if err != nil {
        t.Fatalf("Error occured while getting baremetalhost: (%v)", err)
    }*/


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
func newHost() *metal3v1alpha1.BareMetalHost {


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
                        Name:      "test-bmh",
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
