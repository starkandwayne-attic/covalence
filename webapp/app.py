from flask import Flask
from webapp import webapp
from api import api

#Create and configure app object. 
api = Flask(__name__,static_folder='webapp/static')
api.config.from_object('config.Config')

#Register webapp and API portions of app with the app object created above.
api.register_blueprint(api,url_prefix='/api')

if __name__ == "__main__":
    
    app.run(host='0.0.0.0',port=9200)
