from flask import Blueprint

webapp = Blueprint('webapp',__name__,static_folder='static/',template_folder='templates/')

import routes
