import os
import cherrypy

class Root:

    @cherrypy.expose
    def index(self):
        return 'Hello from {}'.format(os.environ.get('EXAMPLE_SERVER_NAME', 'Unknown'))


def run():
    cherrypy.config.update(
        {
            'environment': 'production',
            'server.socket_host': '0.0.0.0',
            'server.socket_port': 8080,
        })
    cherrypy.quickstart(Root(), '/')


if __name__ == '__main__':
    run()
