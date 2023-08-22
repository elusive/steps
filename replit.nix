{ pkgs }: {
    deps = [
        pkgs.git status
        pkgs.go
        pkgs.gopls
    ];
}