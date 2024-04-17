import subprocess, json
from urllib import request

try:

    db = {
        "ip": subprocess.getoutput("""dig -4 TXT +short o-o.myaddr.l.google.com @ns1.google.com | tr -d '"' """)
    }

    url = "https://hs.timedout.dev/update/"

    db_json = json.dumps(db).encode('utf8')

    req = request.Request(url, data=db_json, headers={'content-type': 'application/json'}, method="POST")
    response = request.urlopen(req)

except Exception as e:
    print("ERROR: ", e)

