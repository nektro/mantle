with import <nixpkgs> {};

pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    gcc
    nodejs
  ];

  hardeningDisable = [ "all" ];
}
