import api
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column,Integer,Text,ForeignKey,Boolean,Float
import time

class Connection(api.Base):

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
