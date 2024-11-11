package super_app

import (
    corev1 "k8s.io/api/core/v1"
)

#Service: corev1.#Service & {
    apiVersion: "v1"
    kind:       "Service"
}
