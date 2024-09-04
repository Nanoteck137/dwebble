{
  description = "Devshell for dwebble";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    gitignore.url = "github:hercules-ci/gitignore.nix";
    gitignore.inputs.nixpkgs.follows = "nixpkgs";

    pyrin.url        = "github:nanoteck137/pyrin/v0.7.0";
    pyrin.inputs.nixpkgs.follows = "nixpkgs";

    devtools.url     = "github:nanoteck137/devtools";
    devtools.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, gitignore, devtools, pyrin, ... }:
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
          subPackages = ["cmd/dwebble"];

          ldflags = [
            "-X github.com/nanoteck137/dwebble.Version=${version}"
            "-X github.com/nanoteck137/dwebble.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          ];

          vendorHash = "sha256-ITNTEG7Esspp0eRdMC9G37GIinI8UAqQH3ukXAL940g=";
        };

        dwebbleWeb = pkgs.buildNpmPackage {
          name = "dwebble-web";
          version = fullVersion;

          src = gitignore.lib.gitignoreSource ./web;
          npmDepsHash = "sha256-0R0zcjevu3yrKC/+7nJsOq4eXEpE0/Y/4IJ/lgtU9oY=";

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

            pyrin.packages.${system}.default
            tools.publishVersion
          ];
        };
      }
    ) // {
      nixosModules.default = import ./nix/dwebble.nix { inherit self; };
      nixosModules.dwebble-web = import ./nix/dwebble-web.nix { inherit self; };
    };
}
