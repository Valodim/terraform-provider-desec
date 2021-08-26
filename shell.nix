let
  pkgs = import (builtins.fetchTarball {
    # Descriptive name to make the store path easier to identify
    name = "nixos-21.05-2021-06-15";
    # Commit hash for specific nixpkgs revision we pin to
    url = "https://github.com/nixos/nixpkgs/archive/1f91fd1040667e9265a760b0347f8bc416249da7.tar.gz";
    # Hash obtained using `nix-prefetch-url --unpack <url>`
    sha256 = "sha256-krXZQ0lObduC95f40K3JwIT//VIBpXBwVNclqh5njtE=";
  }) {
    overlays = [
      # for pinning specific versions (currently not needed)
      # (self: super: { nodejs = self.nodejs-12_x; })
    ];
};
in
  pkgs.mkShell {
    buildInputs = with pkgs; [
      go
      terraform_0_13
    ];
  }
