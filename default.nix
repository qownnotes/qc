{ buildGoModule, installShellFiles, lib }:

buildGoModule rec {
  pname = "qc";
  version = "0.5.0";

  src = builtins.path { path = ./.; name = "qc"; };

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
   installShellCompletion --cmd foobar \ 
     --bash <($out/bin/foobar --bash-completion) \ 
     --fish <($out/bin/foobar --fish-completion) \ 
     --zsh <($out/bin/foobar --zsh-completion) 
  '';

  meta = with lib; {
    description = "QOwnNotes command-line snippet manager";
    homepage = "https://github.com/qownnotes/qc";
    license = licenses.mit;
    maintainers = with maintainers; [ pbek ];
    platforms = platforms.linux ++ platforms.darwin;
  };
}
