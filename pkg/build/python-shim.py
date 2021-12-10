# This file includes a shim that will execute your task code.

try:
    import airplane
except ModuleNotFoundError:
    pass
import importlib.util as util
import json
import os
import sys

def run(args):
    sys.path.append("{{.TaskRoot}}")
    
    if len(args) != 2:
        raise Exception("usage: python ./shim.py <args>")

    os.chdir("{{.TaskRoot}}")
    spec = util.spec_from_file_location("mod.main", "{{ .Entrypoint }}")
    mod = util.module_from_spec(spec)
    spec.loader.exec_module(mod)

    try:
        ret = mod.main(json.loads(args[1]))
        if ret is not None:
            try:
                airplane.set_output(ret)
            except NameError:
                # airplanesdk is not installed - gracefully print to stdout instead.
                # This makes it easier to use the shim in a dev environment. We ensure airplanesdk
                # is installed in production images.
                sys.stdout.flush()
                print("The airplanesdk package must be installed to set return values as task output.", file=sys.stderr)
                print("Printing return values to stdout instead.", file=sys.stderr)
                sys.stderr.flush()
                print(json.dumps(ret, indent=2))
    except Exception as e:
        raise Exception("executing {{.Entrypoint}}") from e

if __name__ == "__main__":
    run(sys.argv)
