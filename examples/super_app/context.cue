package super_app

#Context: Context={
    git_url: string
    subpath: string | *"/"

    namespace: string & =~"^[a-z]([a-z0-9-]+)$"

    region: "us-east-1" | "us-east-2" | "us-west-2"
}
