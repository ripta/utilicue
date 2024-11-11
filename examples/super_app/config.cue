package super_app

#Config: Config={
    name: string
    team: string

    commonAnnotations: [string]: string

    commonLabels: [string]: string
    commonLabels: {
        "app":  Config.name
        "team": Config.team
    }

    commonSelectors: [string]: string
    commonSelectors: {
        "app":  Config.name
    }
}
