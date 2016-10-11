import sys
sys.dont_write_bytecode = True
from flask import Flask, g, Blueprint, url_for, request, jsonify, render_template
#import connection
from sqlalchemy.ext.declarative import declarative_base
from flask.ext.sqlalchemy import SQLAlchemy
import sqlalchemy
from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine
import time
from sqlalchemy import Column,Integer,Text,ForeignKey,Boolean,Float

#Create and configure app object. 
api = Flask(__name__, static_folder='static')
api.config.from_object('config.Config')

#Throws a warning if we don't set this.
api.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

db = SQLAlchemy(api)
#We'll use this db instance throughout the API.

#engine = create_engine('postgresql://localhost:5432')
#Session = sessionmaker(bind=engine)

Base = declarative_base()


class Connection(Base):

    __tablename__ = 'connection'

    id = Column(Integer,primary_key=True)
    uuid = Column(Text)
    source_ip = Column(Text)
    source_port = Column(Text)
    source_deployment = Column(Text)
    source_job = Column(Text)
    source_index = Column(Integer)
    source_user = Column(Text)
    source_group = Column(Text)
    source_pid = Column(Integer)
    source_process_name = Column(Text)
    source_age = Column(Integer)

    destination_ip = Column(Text)
    destination_port = Column(Text)

    def __init__(self,**kwargs):
        self.__dict__.update(**kwargs)
        self.created_at = time.time()

    def serialize(self):

        return {

            'connection': {

                'source': {

                    'ip': self.source_ip,
                    'port': self.source_port,
                    'deployment': self.source_deployment,
                    'job': self.source_job,
                    'index': self.source_index,
                    'user': self.source_user,
                    'group': self.source_group,
                    'pid': self.source_pid,
                    'process_name': self.source_process_name,
                    'age': self.source_age

                },
                'destination': {

                    'ip': self.destination_ip,
                    'port': self.destination_port

                }
            },

            'connection_uuid': self.uuid

        }

Base.metadata.create_all(bind=db.engine)


@api.route('/connections',methods=['GET'])
def get_connections():

    connections = db.session.query(Connection).all()
    connection_list = []
    for con in connections:

        connection_list.append(con.serialize())

    return jsonify({"code":200,"resource":connection_list})


@api.route('/connections',methods=['POST'])
def create_connections():

    print request.data
    params = request.json
    print params
    source = params['source']
    destination = params['destination']

    new_connection = Connection(

        source_ip = source['ip'],
        source_port = source['port'],
        source_deployment_name = source['deployment'],
        source_job = source['job'],
        source_index = source['index'],
        source_user = source['user'],
        source_group = source['group'],
        source_pid = source['pid'],
        source_process_name = source['process_name'],
        source_age = source['age'],
        destination_ip = destination['ip'],
        destination_port = destination['port']

        )

    db.session.add(new_connection)
    db.session.commit()

    return jsonify({"code":200,"message":"Resources created."})

@api.route('/',methods=['GET'])
def index():

    return render_template("index.html")


@api.route('/login',methods=['GET'])
def login():

    return render_template("login.html")


if __name__ == "__main__":
    
    api.run(host='0.0.0.0',port=9201,debug=True)


