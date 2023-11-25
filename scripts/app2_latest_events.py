import requests

HOST = "http://localhost:8080/v1"
TOPIC = "newtopic"
USERNAME = "admin"
PASSWORD = "pass"

def main():
    # login
    login_uri = f"{HOST}/login"
    credentials = {
        "username": USERNAME,
        "password": PASSWORD
    }
    resp = requests.post(login_uri, json=credentials)
    token = resp.json()["token"]
    print("Token", token)

    # get latest event
    headers = {"Authorization": f"Bearer {token}"}
    get_uri = f"{HOST}/kafka/topics/{TOPIC}/latest"
    resp = requests.get(get_uri, headers=headers)
    print(resp.json())

if __name__ == "__main__":
    main()