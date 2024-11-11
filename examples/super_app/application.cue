package super_app

#ApplicationSpec: {
    image: {
        repository: string
        tag:        string
    }
    replicas: int & >=1
    resources: #ResourceRequirements
}

#Application: Application={
    #Basis

    spec: #ApplicationSpec

    _Deployment: #Deployment & {
        _InDeploy: _

        metadata: _InDeploy.metadata

        spec: replicas: _InDeploy.spec.replicas
        spec: strategy: type: "RollingUpdate"

        spec: selector: matchLabels: _InDeploy.config.commonSelectors

        spec: template: {
            metadata: _InDeploy.metadata

            spec: containers: [{
                name:  _InDeploy.config.name
                image: _InDeploy.spec.image.repository + ":" + _InDeploy.spec.image.tag
                resources: _InDeploy.spec.resources
                ports: [{
                    containerPort: 8080
                }]
            }]

            spec: securityContext: {
                runAsNonRoot: true
                runAsUser:    1001
            }
        }
    }

    resources: kubernetes: deployments: (Application.config.name): _Deployment & {
        _InDeploy: {
            config:   Application.config
            metadata: Application.metadata
            spec:     Application.spec
        }
    }
}
