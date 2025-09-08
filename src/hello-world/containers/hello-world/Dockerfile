# We use Apache HTTP Server as the base image
FROM httpd:latest

# The only thing we change is to override the default web page with our own
# hello world web page
COPY html/ /usr/local/apache2/htdocs/
