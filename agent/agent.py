import os
import sys
import requests
import time
api_url = None
update_period = 5

class Connection():

    def __init__(self,text):

        self.text_representation = ""
        self.source_ip = ""
        self.source_port = ""
        self.source_deployment = ""
        self.source_job = ""
        self.source_index = 0
        self.source_user = ""
        self.source_group = ""
        self.source_pid = ""
        self.source_process_name = ""
        self.source_age = 0
        
        self.destination_ip = ""
        self.destination_port = ""
        
        self.protocol = ""
        
        self.text_representation = text
        words = self.text_representation.split()
        
        self.protocol = words[0]
        
        source_ip_and_port = words[3]
        self.source_ip = source_ip_and_port.split(":")[0]
        self.source_port = source_ip_and_port.split(":")[1]

        destination_ip_and_port = words[4] 
        self.destination_ip = destination_ip_and_port.split(":")[0]
        self.destination_port = destination_ip_and_port.split(":")[1]

        user = words[6]
        if user == "0":
            self.source_user = "root"
        elif user == "1000":
            self.source_user = "vcap"
            
        pid_and_program = words[8]
        self.source_pid = pid_and_program.split("/")[0]
        self.process_name = pid_and_program.split("/")[1]

        
    def print_connection(self):

        print "Source IP: " + self.source_ip
        print "Source Port: " + self.source_port
        print "Source Deployment: " + self.source_deployment
        print "Source Job: " + self.source_job
        print "Source Index: " + str(self.source_index)
        print "Source User: " + self.source_user
        print "Source Group: " + self.source_group
        print "Source PID: " + self.source_pid
        print "Source Process Name: " + self.source_process_name
        print "Source Age: " + str(self.source_age)
        print "Destination IP: " + self.destination_ip
        print "Destination Port: " + self.destination_port
        print "Protocol: " + self.protocol


    def serialize(self):

       return  {
           "protocol": self.protocol,
           
           "source":{
               
               "ip": self.source_ip,
               "port": self.source_port,
               "deployment": self.source_deployment,
               "job": self.source_job,
               "index": self.source_index,
               "user": self.source_user,
               "group": self.source_group,
               "pid": self.source_pid,
               "process_name": self.source_process_name,
               "age": self.source_age
               
           },
           "destination":{
               
               "ip":self.destination_ip,
               "port":self.destination_port
               
           }
           
       }


def publish_connection_info():

    os.system("sudo netstat -e -e -n -p | grep -v unix | grep ESTABLISHED | grep tcp | grep -v 127.0.0.1 > netstat_output.txt")
    connection_list = []

    print " "
    print "##############################################"
    print " "
    
    with open('netstat_output.txt') as net_file:
        
        for line in net_file:
            print "Creating Connection object from following line:"
            print line
            connection = Connection(line)
            connection.print_connection()
            connection_list.append(connection.serialize())

        for conn in connection_list:
            response = requests.post(api_url+'/connections',json=conn,headers={'Content-Type':'application/json'})

    print " "
    print "###############################################"
    print " "
       
if __name__ == "__main__":

    if len(sys.argv) < 3:

        print "Covalence agent requires location of Covalence API as the first argument and reporting period (in seconds) for second argument"
        print "Example: python agent.py http://localhost:9200 5"
        sys.exit(1)

    else:

        api_url = sys.argv[1]
        update_period = float(sys.argv[2])
        
    while True:

        publish_connection_info()
        time.sleep(update_period)


