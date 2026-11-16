"""Traditional Paramiko usage through the local Go SSH proxy."""
import os
import paramiko

PROXY_HOST = "127.0.0.1"
PROXY_PORT = 2222

client = paramiko.SSHClient()
client.load_system_host_keys()
client.set_missing_host_key_policy(paramiko.RejectPolicy())

try:
    client.connect(
        hostname=PROXY_HOST,
        port=PROXY_PORT,
        username=os.environ["SSH_USERNAME"],
        password=os.environ["SSH_PASSWORD"],
        look_for_keys=False,
        allow_agent=False,
    )
    stdin, stdout, stderr = client.exec_command("uname -a")
    print(stdout.read().decode(), end="")
    error = stderr.read().decode()
    if error:
        print(error, end="")
finally:
    client.close()
