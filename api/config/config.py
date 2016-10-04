
class Config():

    DEBUG = True
    TESTING = False
    ALLOWED_EXTENSIONS = {'zip'}
    SQLALCHEMY_DATABASE_URI = 'sqlite:///development.db'
