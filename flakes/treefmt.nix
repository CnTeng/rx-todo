{ inputs, ... }:
{
  imports = [ inputs.treefmt.flakeModule ];

  perSystem = {
    treefmt = {
      projectRootFile = "flake.nix";

      programs = {
        gofmt.enable = true;
        nixfmt.enable = true;
        prettier.enable = true;
      };
    };
  };
}
