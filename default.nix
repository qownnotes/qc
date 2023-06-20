{ buildGoModule, installShellFiles, lib }:

buildGoModule rec {
  pname = "qc";
  version = "0.5.0";

  src = builtins.fetchGit ./.;

  vendorSha256 = "sha256-7t5rQliLm6pMUHhtev/kNrQ7AOvmA/rR93SwNQhov6o=";

  ldflags = [
    "-s" "-w" "-X=github.com/qownnotes/qc/cmd.version=${version}"
  ];

  doCheck = false;

  subPackages = [ "." ];

  nativeBuildInputs = [
    installShellFiles
  ];

  postInstall = ''
    installShellCompletion --cmd qc \
      --zsh ./misc/completions/zsh/_qc
  '';

  meta = with lib; {
    description = "QOwnNotes command-line snippet manager";
    homepage = "https://github.com/qownnotes/qc";
    license = licenses.mit;
    maintainers = with maintainers; [ pbek ];
    platforms = platforms.linux ++ platforms.darwin;
  };
}
