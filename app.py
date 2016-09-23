from flask import Flask
from webapp import webapp
from api import api

#Create and configure app object. 
app = Flask(__name__,static_folder='webapp/static')
app.config.from_object('config.Config')

#Register webapp and API portions of app with the app object created above.
app.register_blueprint(webapp)
app.register_blueprint(api,url_prefix='/api')

if __name__ == "__main__":
    
    app.run(host='0.0.0.0',port=9200)
