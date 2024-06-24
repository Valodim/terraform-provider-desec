let
  pkgs = import <nixpkgs> {
    overlays = [
      # for pinning specific versions (currently not needed)
      # (self: super: { nodejs = self.nodejs-12_x; })
    ];
};
in
  pkgs.mkShell {
    buildInputs = with pkgs; [
      go
      goreleaser
      terraform
    ];
  }
