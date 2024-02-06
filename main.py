from flask import Flask, redirect, request
app = Flask(__name__)

SECRET = None

def read_ip():
    with open("ip", "r") as f:
        return f.read()
    
def update_ip(ip):
    with open("ip", "w") as f:
        f.write(ip)

def get_secret():
    global SECRET
    if SECRET == None:
        try:
            with open("secret", "r") as f:
                s = f.read().strip()
                if len(s) == 0:
                    raise ValueError("Secret file empty len(s) is 0")
                SECRET = s
        except Exception as e:
            print("Error while reading secret")
            print(e)
            raise Exception("Error while getting secret")
    return SECRET

@app.route("/")
def home():
    ip = read_ip()
    try:
        return redirect(ip)
    except Exception as e:
        print("Something went wrong:")
        print(e)
    return "Failed to redirect, check logs"

@app.route("/update", methods=['POST'])
def update():
    data = request.json
    SECRET = get_secret()
    if data.get('secret') != SECRET or data.get('secret') == None or data.get('ip') == None:
        return "Bad Request"
    else:
        update_ip(data['ip'])
    return "Updated IP"
