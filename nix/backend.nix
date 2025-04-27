{ self }: 
{ config, lib, pkgs, ... }:
with lib; let
  cfg = config.services.dwebble;

  dwebbleConfig = pkgs.writeText "config.toml" ''
    listen_addr = "${cfg.host}:${toString cfg.port}"
    data_dir = "/var/lib/dwebble"
    library_dir = "${cfg.libraryDir}"
    username = "${cfg.username}"
    initial_password = "${cfg.initialPassword}"
    jwt_secret = "${cfg.jwtSecret}"
  '';
in
{
  options.services.dwebble = {
    enable = mkEnableOption "Enable the dwebble service";

    port = mkOption {
      type = types.port;
      default = 7550;
      description = "port to listen on";
    };

    host = mkOption {
      type = types.str;
      default = "";
      description = "hostname or address to listen on";
    };

    username = mkOption {
      type = types.str;
      description = "username of the first user";
    };

    initialPassword = mkOption {
      type = types.str;
      description = "initial password of the first user (should change after the first login)";
    };

    jwtSecret = mkOption {
      type = types.str;
      description = "jwt secret";
    };

    libraryDir = mkOption {
      type = types.path;
      description = "library directory";
    };

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.system}.backend;
      description = "package to use for this service (defaults to the one in the flake)";
    };

    user = mkOption {
      type = types.str;
      default = "dwebble";
      description = "user to use for this service";
    };

    group = mkOption {
      type = types.str;
      default = "dwebble";
      description = "group to use for this service";
    };

    openFirewall = mkOption {
      type = types.bool;
      default = false;
      description = "open the ports in the firewall";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.dwebble = {
      description = "dwebble";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;

        StateDirectory = "dwebble";

        ExecStart = "${cfg.package}/bin/dwebble serve -c '${dwebbleConfig}'";

        Restart = "on-failure";
        RestartSec = "5s";

        PrivateTmp = true;
        ProtectHome = true;
        ProtectHostname = true;
        ProtectKernelLogs = true;
        ProtectKernelModules = true;
        ProtectKernelTunables = true;
        ProtectProc = "invisible";
        ProtectSystem = "strict";
        RestrictAddressFamilies = [ "AF_INET" "AF_INET6" "AF_UNIX" ];
        RestrictNamespaces = true;
        RestrictRealtime = true;
        RestrictSUIDSGID = true;
      };
    };

    networking.firewall = lib.mkIf cfg.openFirewall {
      allowedTCPPorts = [ cfg.port ];
    };

    users.users = mkIf (cfg.user == "dwebble") {
      dwebble = {
        group = cfg.group;
        isSystemUser = true;
      };
    };

    users.groups = mkIf (cfg.group == "dwebble") {
      dwebble = {};
    };
  };
}
