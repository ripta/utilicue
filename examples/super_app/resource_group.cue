package super_app

#ID: string

#ResourceGroup: {
    kubernetes: #KubernetesResources
}

#KubernetesResources: {
    deployments: [#ID]: #Deployment
    services:    [#ID]: #Service
}
