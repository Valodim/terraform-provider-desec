{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "terraform-provider-desec";
  version = "0.6.1";

  src = ./.;
  vendorHash = "sha256-APeSigqPXoAoPLmC9YCRIfo0vkChzHQocqM3LqoA1/w=";

  meta = {
    description = "Terraform/OpenTofu provider for deSEC DNS";
    homepage = "https://github.com/Valodim/terraform-provider-desec";
    license = lib.licenses.mpl20;
    mainProgram = "terraform-provider-desec";
  };
}
