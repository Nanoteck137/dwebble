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

        swag = pkgs.buildGoModule rec {
          name = "swag";
          vendorHash = "sha256-BxWmEcx5IIT/yI46CJGE0vE1BRm5zwngc0x1dVy/04s=";
          src = pkgs.fetchFromGitHub {
            owner = "swaggo";
            repo = "swag";
            rev = "76695ca";
            hash = "sha256-+YjmYf0BYsjXL8rztPTdQtQhfiTwznGXL/73EzuUu6g=";
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
          ];
        };
      }
    );
}
