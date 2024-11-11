package super_app

import (
    appsv1 "k8s.io/api/apps/v1"
)

#Deployment: appsv1.#Deployment & {
    apiVersion: "apps/v1"
    kind:       "Deployment"
}
