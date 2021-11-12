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
                raise Exception("airplanesdk package must be installed to output return values - add airplanesdk to requirements.txt or remove return value from task")
    except Exception as e:
        raise Exception("executing {{.Entrypoint}}") from e

if __name__ == "__main__":
    run(sys.argv)
