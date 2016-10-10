from flask import Flask
from flask import render_template

#Create and configure app object. 
webapp = Flask(__name__,static_folder='static')
webapp.config.from_object('config.Config')

@webapp.route('/',methods=['GET'])
def index():

    return render_template("index.html")


@webapp.route('/login',methods=['GET'])
def login():

    return render_template("login.html")

if __name__ == "__main__":
    
    webapp.run(host='0.0.0.0',port=9200)

