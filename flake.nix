{
  description = "Devshell for dwebble";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";
    templ.url        = "github:a-h/templ";
  };

  outputs = { self, nixpkgs, templ, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        swag = pkgs.buildGoModule {
          name = "swag";
          vendorHash = "sha256-sxW4WgtsuVQQiacM24wU81cXKEJuCxT39uUqh2nMZ0k=";
          src = pkgs.fetchFromGitHub {
            owner = "swaggo";
            repo = "swag";
            rev = "v1.16.2";
            hash = "sha256-lLsrwWuSfyU5C8cJUNVihSqQTbr28yVcTVej8fW5oP4=";
          };

          subPackages = [ "cmd/swag" ];
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            nodejs_20
            ffmpeg
            air
            templ.packages.${system}.templ
            sqlc
            go
            swag
            gotools
          ];
        };
      }
    );
}
