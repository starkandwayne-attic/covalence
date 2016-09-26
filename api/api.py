import sys
sys.dont_write_bytecode = True
from flask import Flask, g
from api import api, Base
from flask.ext.sqlalchemy import SQLAlchemy

#Create and configure app object. 
app = Flask(__name__,static_folder='webapp/static')
app.config.from_object('config.Config')

#Throws a warning if we don't set this.
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

#We'll use this db instance throughout the API.
db = SQLAlchemy(app)
Base.metadata.create_all(bind=db.engine)

#Register webapp and API portions of app with the app object created above.
app.register_blueprint(api)

if __name__ == "__main__":
    
    app.run(host='0.0.0.0',port=9201)
