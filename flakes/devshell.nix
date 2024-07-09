{ ... }:
{
  perSystem =
    {
      config,
      pkgs,
      lib,
      ...
    }:
    {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [ go ];
        shellHook = config.pre-commit.installationScript;
      };
    };
}
