{
  description = "Devshell for dwebble";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    gitignore.url = "github:hercules-ci/gitignore.nix";
    gitignore.inputs.nixpkgs.follows = "nixpkgs";

    devtools.url     = "github:nanoteck137/devtools";
    devtools.inputs.nixpkgs.follows = "nixpkgs";

    tagopus.url      = "github:nanoteck137/tagopus/v0.1.1";
    tagopus.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, gitignore, devtools, tagopus, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        fullVersion = ''${version}-${self.dirtyShortRev or self.shortRev or "dirty"}'';

        dwebble = pkgs.buildGoModule {
          pname = "dwebble";
          version = fullVersion;
          src = ./.;
          subPackages = ["cmd/dwebble" "cmd/dwebble-import" "cmd/dwebble-dl"];

          ldflags = [
            "-X github.com/nanoteck137/dwebble.Version=${version}"
            "-X github.com/nanoteck137/dwebble.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          ];

          vendorHash = "sha256-jqGXKDv03X48T/J7PX7jU09TJVbov2DlZzHomVIrX9o=";

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postFixup = ''
            wrapProgram $out/bin/dwebble --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick ]}
            wrapProgram $out/bin/dwebble-import --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick ]}
            wrapProgram $out/bin/dwebble-dl --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ffmpeg pkgs.imagemagick tagopus.packages.${system}.default ]}
          '';
        };

        dwebbleWeb = pkgs.buildNpmPackage {
          name = "dwebble-web";
          version = fullVersion;

          src = gitignore.lib.gitignoreSource ./web;
          npmDepsHash = "sha256-fTxzl5ZIp75q6MXm8GMuyR0zuMN2T+STERQWZXvP65o=";

          PUBLIC_VERSION=version;
          PUBLIC_COMMIT=self.dirtyRev or self.rev or "no-commit";

          installPhase = ''
            runHook preInstall
            cp -r build $out/
            echo '{ "type": "module" }' > $out/package.json

            mkdir $out/bin
            echo -e "#!${pkgs.runtimeShell}\n${pkgs.nodejs}/bin/node $out\n" > $out/bin/dwebble-web
            chmod +x $out/bin/dwebble-web

            runHook postInstall
          '';
        };

        tools = devtools.packages.${system};
      in
      {
        packages = {
          default = dwebble;
          dwebble = dwebble;
          dwebble-web = dwebbleWeb;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            air
            go
            gopls
            nodejs
            imagemagick
            ffmpeg

            tagopus.packages.${system}.default
            tools.publishVersion
          ];
        };
      }
    ) // {
      nixosModules.default = import ./nix/dwebble.nix { inherit self; };
      nixosModules.dwebble-web = import ./nix/dwebble-web.nix { inherit self; };
    };
}
