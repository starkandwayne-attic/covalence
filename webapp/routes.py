from . import webapp
from flask import render_template

@webapp.route('/',methods=['GET'])
def index():

    return render_template("index.html")


@webapp.route('/login',methods=['GET'])
def login():

    return render_template("login.html")
