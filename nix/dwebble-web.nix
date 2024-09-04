{ self }:
{ config, lib, pkgs, ... }:
with lib; let
  cfg = config.services.dwebble-web;
in
{
  options.services.dwebble-web = {
    enable = mkEnableOption "Enable the dwebble-web service";

    port = mkOption {
      type = types.port;
      default = 7551;
      description = "port to listen on";
    };

    host = mkOption {
      type = types.str;
      default = "";
      description = "hostname or address to listen on";
    };

    apiAddress = mkOption {
      type = types.str;
      description = "address to the api server";
    };

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.system}.dwebble-web;
      description = "package to use for this service (defaults to the one in the flake)";
    };

    user = mkOption {
      type = types.str;
      default = "dwebble-web";
      description = lib.mdDoc "user to use for this service";
    };

    group = mkOption {
      type = types.str;
      default = "dwebble-web";
      description = lib.mdDoc "group to use for this service";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.dwebble-web = {
      description = "Frontend for dwebble";
      wantedBy = [ "multi-user.target" ];

      environment = {
        PORT = "${toString cfg.port}";
        HOST = "${cfg.host}";
        API_ADDRESS = "${cfg.apiAddress}";
        HOST_HEADER = "x-forwarded-host";
      };

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;

        ExecStart = "${cfg.package}/bin/dwebble-web";

        Restart = "on-failure";
        RestartSec = "5s";

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

    users.users = mkIf (cfg.user == "dwebble-web") {
      dwebble-web = {
        group = cfg.group;
        isSystemUser = true;
      };
    };

    users.groups = mkIf (cfg.group == "dwebble-web") {
      dwebble-web = {};
    };
  };
}
