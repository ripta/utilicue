package kube_deployment

import (
    appsv1 "k8s.io/api/apps/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

#Deployment: {
    apiVersion: "apps/v1"
    kind:       "Deployment"
    metadata: metav1.#ObjectMeta & {
        name: string
        labels: {
            app: string
        }
    }
    spec: appsv1.#DeploymentSpec
}
