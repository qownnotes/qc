{
  pkgs ? import <nixpkgs> { },
}:
pkgs.mkShell {
  # nativeBuildInputs is usually what you want -- tools you need to run
  nativeBuildInputs = with pkgs; [
    just
    go
  ];

  shellHook = ''
    # Symlink the pre-commit hook into the .git/hooks directory
    ln -sf ../../scripts/pre-commit.sh .git/hooks/pre-commit
  '';
}
