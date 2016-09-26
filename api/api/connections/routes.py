from connections import connections
from flask import request, jsonify, g
from model import connection

@connections.route('/',methods=['GET'])
def get_connections():

    connections = g.db.query(Match).filter(Match.id == match_id).all()
    connection_list = []
    for connection in connections:

        connection_list.append(connection.serialize())

    
@connections.route('/',methods=['POST'])
def create_connections():
    
    params = request.json
    source = params['source']
    destination = params['destination']
    
    new_connection = Connection(
        
        source_ip = source['ip'],
        source_destination = source['destination'],
        source_deployment_name = source['deployment_name']
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
    
    g.db.add(new_connection)
    
    return jsonify({"code":200,"message":"Resources created."})

    
    

    
    
