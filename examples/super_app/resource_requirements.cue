package super_app

import (
    corev1 "k8s.io/api/core/v1"
)

#ResourceList: corev1.#ResourceList & {
    cpu:    number
    memory: int
}

#ResourceRequirements: corev1.#ResourceRequirements & {
    limits: #ResourceList & {
        cpu:    <=94.0
        memory: <376Gi
    }

    requests: #ResourceList & {
        cpu:    >=0.1 & <=limits.cpu
        memory: >=128Mi & <=limits.memory
    }
}
