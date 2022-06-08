# Environment Custom Resource Definition

The **Environment** is a special place where is possible to create some resources, like schema-registry, kafka-cluster, etc.

An environment has Clusters, Network Management, Schema Registry and Environment Settings, is possible create more than one cluster by environment, but only one schema-registry;

In a environment is possible to create RBAC rules, to add more security betteween envrionments, in other words and enviroment provides an isolation, a security space with specifc rules;

## CRD Parameters

| Parameter | Value | Description |
|-------------|-------------|------------------------------------------------|
| name        | **string**  | Defines the names of Environment |
|resourceExist| **boolean** | Defines if the environment exists or not exist | 

## **CRD Examples**

Below, we have two kind of Custom Resource Definion, the capabilities of this CRD's is create and Environment at Cloud Clonfluent or just retrive an existent.

When we have this Environment syncronized with the Operator, it will create a secret inside of the namespace with the following content:

```yaml

    environmentName: kubbee-pay
    environmentId: uire78

```

The Environment name should be unique, and inside of the environment we can have more than one cluster kafka, example:

 + development
 + QualityAssurance
 + Production

> <span style="color:orange">**WARNING**</span>: &nbsp;Be careful, the name of the environmets should be unique at Cloud Confluent.

### **Resource not exist at cloud confluent**

When is necessary to create a new _**Environment**_ at Cloud Confluent, we can use the CRD if the flag _**resourceExist: false**_ the environment will be create at Cloud Confluent and a refenrece also be create inside the namespace defined during the execution of CRD.


```yaml
apiVersion: messages.kubbee.tech/v1alpha1
kind: CCloudEnvironment
metadata:
  name: ccloudenvironment-kubbee
spec:
  name: kubbee-pay
  environmentResource:
    resourceExist: false
```

### **Resource exists at cloud confluent**

If the environment already exists at Cloud Confluent, we can use the flag _**resourceExist: true**_ the CRD will create a reference of this environment inside a namespace defined during the execution of CRD.

```yaml
apiVersion: messages.kubbee.tech/v1alpha1
kind: CCloudEnvironment
metadata:
  name: ccloudenvironment-kubbee
spec:
  name: kubbee-pay
  environmentResource:
    resourceExist: true

```