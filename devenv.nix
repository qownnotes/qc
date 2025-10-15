{
  # https://devenv.sh/languages/
  languages.go.enable = true;

  enterShell = ''
    echo "🛠️ qc dev shell"
    echo "🐹 Go version: $(go version)"
  '';

  # https://devenv.sh/git-hooks/
  git-hooks = {
    hooks = {
      gofmt.enable = true;
      golangci-lint.enable = true;
      golines.enable = true;
      gotest.enable = true;
      govet.enable = true;
      gitlint.enable = true;
    };
  };

  # See full reference at https://devenv.sh/reference/options/
}
