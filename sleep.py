import time
import requests

url = "http://localhost:8080/"

while True:
    try:
        response = requests.get(url)
    except requests.exceptions.RequestException as e:
        print("Error:", e)

    time.sleep(3)  # wait for 3 seconds before the next request
