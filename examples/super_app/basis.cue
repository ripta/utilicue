package super_app

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

#Basis: Basis={
    config:  #Config
    context: #Context

    metadata: metav1.#ObjectMeta & {
        name:      Basis.config.name
        namespace: Basis.context.namespace

        annotations: Basis.config.commonAnnotations
        labels:      Basis.config.commonLabels
    }

    spec: _

    resources: #ResourceGroup
}
