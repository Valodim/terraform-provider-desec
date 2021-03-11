let
  pkgs = import (builtins.fetchTarball {
    # Descriptive name to make the store path easier to identify
    name = "nixos-20.09-2020-11-05";
    # Commit hash for specific nixpkgs revision we pin to
    url = "https://github.com/nixos/nixpkgs/archive/d105075a1fd870b1d1617a6008cb38b443e65433.tar.gz";
    # Hash obtained using `nix-prefetch-url --unpack <url>`
    sha256 = "1jcs44wn0s6mlf2jps25bvcai1rij9b2dil6zcj8zqwf2i8xmqjh";
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
