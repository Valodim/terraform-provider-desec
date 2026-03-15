{
  description = "Terraform/OpenTofu provider for deSEC DNS";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      perSystem = {pkgs, ...}: let
        pkg = pkgs.callPackage ./default.nix {};
      in {
        packages.default = pkg;
        packages.terraform-provider-desec = pkg;

        devShells.default = pkgs.callPackage ./shell.nix {};
      };
    };
}
