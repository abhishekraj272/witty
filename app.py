from flask import Flask, request
app = Flask(__name__)

@app.route('/upload', methods = ['POST'])
def upload():
    print(request.headers.__dict__)
    return {'success':True}, 200, {'ContentType':'application/json'} 