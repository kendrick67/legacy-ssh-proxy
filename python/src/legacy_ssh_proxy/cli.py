import subprocess
import sys
from . import get_binary_path

def main():
    binary = get_binary_path()
    result = subprocess.run([str(binary), *sys.argv[1:]])
    sys.exit(result.returncode)
