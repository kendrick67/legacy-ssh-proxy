import paramiko
from legacy_ssh_proxy import legacy_ssh_proxy


real_host = "aggrouter.cissta.local.net"
port = "22"
username = "user"
password = "secret"

def paramiko_router_test():
    client = paramiko.SSHClient()
    client.load_system_host_keys()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    try:
        client.connect(
            hostname=real_host, port=port, username=username, password=password,
            look_for_keys=False, allow_agent=False,
        )
        _, stdout, _ = client.exec_command("show platform")
        print(stdout.read().decode())
    finally:
        client.close()


def proxy_test():
    with legacy_ssh_proxy(f"{real_host}:22", insecure_ignore_host_key=True) as (host, port):
        client = paramiko.SSHClient()
        client.load_system_host_keys()
        client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        print(f"we will connect paramiko to {real_host} via {host}:{port}")
        try:
            client.connect(
                hostname=host, port=port, username=username, password=password,
                look_for_keys=False, allow_agent=False,
            )
            _, stdout, _ = client.exec_command("show platform")
            print(stdout.read().decode())
        finally:
            client.close()


def main():
    # paramiko_router_test()
    # proxy_test()
    pass


main()
