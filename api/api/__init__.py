from flask import Blueprint
from sqlalchemy.ext.declarative import declarative_base

api = Blueprint('api',__name__)

Base = declarative_base()
