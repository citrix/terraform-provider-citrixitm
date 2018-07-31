FROM python:3.6-slim-stretch
RUN pip install cherrypy
RUN useradd -r -s /bin/false citrix
USER citrix
WORKDIR /webapp
COPY app.py .
CMD ["python", "/webapp/app.py"]
