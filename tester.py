import requests
import time
from datetime import datetime
import signal
import sys

versions = {}

def interrupt_handler(signum, frame):
    print(f'Handling signal {signum} ({signal.Signals(signum).name}).')
    sum = 0
    for version in versions:
        sum += versions[version]
    for version in versions:
        print("\t" + version.strip() + " " + str(versions[version]) + " " + str(round(versions[version]/sum, 2) * 100))
    sys.exit(0) 


def main():
    i = 0
    while True:
        i+=1
        req = requests.get("http://localhost:8080")
        print(datetime.now().strftime('%H:%M:%S') + " " + str(req.text).strip() + " " + str(req.status_code).strip())
        txt = req.json()["Text"]
        version = req.headers["X-App-Version"]
        if version + "=>" + str(req.status_code) in versions:
            versions[version + "=>" + str(req.status_code)]+=1
        else:
            versions[version + "=>" + str(req.status_code)]=1
        time.sleep(0.005)
        if i % 100 == 0:
            print("Current_status: ")
            for version in versions:
                print("\t" + version.strip() + " " + str(versions[version]))


if __name__ == '__main__':
    signal.signal(signal.SIGINT, interrupt_handler)

    main()